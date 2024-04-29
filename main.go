package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/searcheroflulz/go_final_project/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("variables.env")
	if err != nil {
		panic("Error loading variables.env file")
	}

	port := os.Getenv("PORT")
	router := gin.Default()

	db, err := sql.Open("sqlite3", os.Getenv("DBFILE"))
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

	router.NoRoute(gin.WrapH(http.FileServer(gin.Dir(os.Getenv("WEBDIR"), false))))
	router.GET("/api/nextdate", handlers.NextDateHandler)
	router.POST("/api/task", handlers.AddTaskHandler(db))
	router.GET("/api/tasks", handlers.GetTasksListHandler(db))
	router.GET("/api/task", handlers.GetTaskHandler(db))
	router.PUT("/api/task", handlers.PutTaskHandler(db))
	router.POST("/api/task/done", handlers.TaskDoneHandler(db))
	router.DELETE("/api/task", handlers.TaskDeleteHandler(db))

	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("error running server: %s", err)
	}
}
