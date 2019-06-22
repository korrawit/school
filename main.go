package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type todo struct {
	ID     int     `json:"id'`
	Title  *string `json:"title"`
	Status string  `json:"status"`
}

func main() {
	r := gin.Default()

	r.GET("/api/todos", getTodos)
	// r.GET("/api/todos/:id", getTodosByIdHandler)
	// r.POST("/api/todos", postTodosHandler)
	// r.DELETE("/api/todos/:id", deleteTodosByIdHandler)

	r.Run(":1234")
}

func getTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err.Error())
	}
	stmt, err := db.Prepare("SELECT id, title, status FROM todos ORDER BY id")
	if err != nil {
		fmt.Println(err.Error())
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err.Error())
	}
	todos := []todo{}
	for rows.Next() {
		t := todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error()},
			)
			return
		}
		todos = append(todos, t)
	}
	c.JSON(http.StatusOK, todos)
}
