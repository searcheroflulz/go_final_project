package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	port := os.Getenv("port")
	router := gin.Default()

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		panic(err)
	}

	router.NoRoute(gin.WrapH(http.FileServer(gin.Dir(os.Getenv("WebDir"), false))))
	router.GET("/api/nextdate", NextDateHandler)
	router.POST("/api/task", AddTaskHandler(db))
	router.GET("/api/tasks", GetTasksListHandler(db))
	router.GET("/api/task", GetTaskHandler(db))
	router.PUT("/api/task", PutTaskHandler(db))
	router.POST("/api/task/done", TaskDoneHandler(db))
	router.DELETE("/api/task", TaskDeleteHandler(db))

	router.Run(":" + port)
}
