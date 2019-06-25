package repository

import (
	"strconv"

	"github.com/korrawit/school/database"
)

type Todo struct {
	ID     int     `json:"id"`
	Title  *string `json:"title"`
	Status string  `json:"status"`
}

type Repository struct {
	DB database.DB
}

type TodoRepository interface {
	GetTodos() ([]Todo, error)
	GetTodoById(id string) (*Todo, error)
	CreateNewTodo(t *Todo) error
	UpdateTodo(id string, t *Todo) error
	DeleteTodoById(id string) error
}

func (r Repository) GetTodos() ([]Todo, error) {
	db, err := r.DB.Connect()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos ORDER BY id")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	todos := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (r Repository) GetTodoById(id string) (*Todo, error) {
	db, err := r.DB.Connect()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos WHERE id = $1")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	var t Todo
	err = row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r Repository) CreateNewTodo(t *Todo) error {
	db, err := r.DB.Connect()
	defer db.Close()

	if err != nil {
		return err
	}

	query := `
			INSERT INTO todos (title, status) VALUES ($1,$2) RETURNING id
		`

	row := db.QueryRow(query, t.Title, t.Status)
	err = row.Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateTodo(id string, t *Todo) error {
	db, err := r.DB.Connect()
	defer db.Close()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE todos SET title=$2, status=$3 WHERE id=$1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, t.Title, t.Status)
	if err != nil {
		return err
	}

	t.ID, err = strconv.Atoi(id)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) DeleteTodoById(id string) error {
	db, err := r.DB.Connect()
	defer db.Close()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("DELETE FROM todos WHERE id = $1")
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
