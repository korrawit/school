package main

import (
	"github.com/gin-gonic/gin"
	"github.com/korrawit/school/database"
	"github.com/korrawit/school/repository"
	"github.com/korrawit/school/todo"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()

	db := database.PostgresDB{}

	t := todo.TodoContext{
		Repo: repository.Repository{
			DB: db,
		},
	}

	r.GET("/api/todos", t.GetTodos)
	r.GET("/api/todos/:id", t.GetTodosByIdHandler)
	r.POST("/api/todos", t.PostTodosHandler)
	r.PUT("/api/todos/:id", t.PutTodosHandler)
	r.DELETE("/api/todos/:id", t.DeleteTodosByIDHandler)

	r.Run(":1234")
}
