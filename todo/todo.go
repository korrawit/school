package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korrawit/school/repository"
)

type TodoContext struct {
	Repo repository.TodoRepository
}

func (tc TodoContext) GetTodos(c *gin.Context) {
	todos, err := tc.Repo.GetTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (tc TodoContext) GetTodosByIdHandler(c *gin.Context) {
	id := c.Param("id")

	todo, err := tc.Repo.GetTodoById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)

}

func (tc TodoContext) PostTodosHandler(c *gin.Context) {
	var t repository.Todo
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = tc.Repo.CreateNewTodo(&t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (tc TodoContext) PutTodosHandler(c *gin.Context) {
	var t repository.Todo
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	err = tc.Repo.UpdateTodo(id, &t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (tc TodoContext) DeleteTodosByIDHandler(c *gin.Context) {
	id := c.Param("id")

	err := tc.Repo.DeleteTodoById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
