package models

import (
	"fmt"
	"github.com/ichtrojan/go-todo/config"
	"time"
)

type Todo struct {
	Id        int
	Item      string
	Completed int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FormatCreatedAt formats and returns the creation time as a string
func (todo Todo) FormatCreatedAt() string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// Fallback to UTC if location loading fails
		return todo.CreatedAt.Format("2006-01-02 15:04:05")
	}
	// Convert the time to China timezone (UTC+8)
	return todo.CreatedAt.In(loc).Format("2006-01-02 15:04:05")
}

// FormatUpdatedAt formats and returns the update time as a string
func (todo Todo) FormatUpdatedAt() string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// Fallback to UTC if location loading fails
		return todo.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	// Convert the time to China timezone (UTC+8)
	return todo.UpdatedAt.In(loc).Format("2006-01-02 15:04:05")
}

// UpdateTodo updates an existing todo item in the database
func UpdateTodo(todo Todo) error {
	fmt.Println("UpdateTodo received todo:", todo.Id, todo.Item, todo.Completed, "UpdatedAt:", todo.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	database := config.Database()

	fmt.Println("Executing SQL with parameters:", 
		"item=", todo.Item, 
		"completed=", todo.Completed, 
		"updated_at=", todo.UpdatedAt.Format("2006-01-02 15:04:05"), 
		"id=", todo.Id)
	
	_, err := database.Exec(`UPDATE todos SET item = ?, completed = ?, updated_at = ? WHERE id = ?`, 
		todo.Item, todo.Completed, todo.UpdatedAt, todo.Id)

	if err != nil {
		fmt.Println("SQL update error:", err)
		return err
	}
	
	fmt.Println("SQL update successful for todo ID:", todo.Id)
	return err
}