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
	ID     int     `json:"id"`
	Title  *string `json:"title"`
	Status string  `json:"status"`
}

func main() {
	r := gin.Default()

	r.GET("/api/todos", getTodos)
	r.GET("/api/todos/:id", getTodosByIdHandler)
	r.POST("/api/todos", postTodosHandler)
	r.DELETE("/api/todos/:id", deleteTodosByIDHandler)

	r.Run(":1234")
}

func getTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stmt, err := db.Prepare("SELECT id, title, status FROM todos ORDER BY id")
	defer stmt.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

func getTodosByIdHandler(c *gin.Context) {
	id := c.Param("id")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos WHERE id = $1")
	defer stmt.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	row := stmt.QueryRow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var t todo
	err = row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func postTodosHandler(c *gin.Context) {

	var t todo
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO todos (title, status) VALUES ($1,$2) RETURNING id
	`

	row := db.QueryRow(query, t.Title, t.Status)
	err = row.Scan(&t.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func deleteTodosByIDHandler(c *gin.Context) {
	id := c.Param("id")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("DELETE FROM todos WHERE id = $1")
	defer stmt.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rs, err := stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r, err := rs.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("row affected", r)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
