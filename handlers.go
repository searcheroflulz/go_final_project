package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NextDateHandler(c *gin.Context) {
	now, err := time.Parse("20060102", c.Query("now"))
	if err != nil {
		c.String(400, "Ошибка: неправильный формат параметра now")
		return
	}

	date := c.Query("date")
	repeat := c.Query("repeat")

	next, err := NextDate(now, date, repeat)
	if err != nil {
		c.String(400, "Ошибка: "+err.Error())
		return
	}

	c.String(200, next)
}

func AddTaskHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось прочитать JSON"})
			return
		}

		if task.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан заголовок задачи"})
			return
		}

		err := task.VerifyChange()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := db.Exec(AddTask, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Не удалось выполнить запрос ДБ"})
			return
		}
		id, err := res.LastInsertId()
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func GetTasksListHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(GetTasksList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var tasks []Task

		for rows.Next() {
			var task Task
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tasks = append(tasks, task)
		}

		if len(tasks) == 0 {
			tasks = []Task{}
		}

		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}
}

func GetTaskHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан идентификатор задачи"})
			return
		}

		var task Task
		err := db.QueryRow(GetTask, id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Задача не найдена"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		})
	}
}

func PutTaskHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось прочитать JSON"})
			return
		}

		if task.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан заголовок задачи"})
			return
		}

		var existingTask Task
		err := db.QueryRow(GetTask, task.ID).Scan(&existingTask.ID, &existingTask.Date, &existingTask.Title, &existingTask.Comment, &existingTask.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса к базе данных"})
			return
		}
		if existingTask.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача с указанным идентификатором не найдена"})
			return
		}

		if err := task.VerifyChange(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = db.Exec(UpdateTask, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса к базе данных"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func TaskDoneHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан идентификатор задачи"})
			return
		}

		var task Task
		err := db.QueryRow(GetTask, c.Query("id")).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса к базе данных"})
			return
		}

		if task.Repeat == "" {
			_, err := db.Exec(DeleteTask, task.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении задачи из базы данных"})
				return
			}
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при вычислении следующей даты выполнения"})
			return
		}

		_, err = db.Exec(UpdateDate, nextDate, task.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении задачи в базе данных"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func TaskDeleteHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан идентификатор задачи"})
			return
		}

		var task Task
		err := db.QueryRow(GetTask, c.Query("id")).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса к базе данных"})
			return
		}

		_, err = db.Exec(DeleteTask, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении задачи из базы данных"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
