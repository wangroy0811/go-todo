package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/config"
	"github.com/ichtrojan/go-todo/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var (
	id        int
	item      string
	completed int
	createdAt time.Time
	view      = template.Must(template.ParseFiles("./views/index.html"))
	database  = config.Database()
)

func Show(w http.ResponseWriter, r *http.Request) {
	statement, err := database.Query(`SELECT id, item, completed, created_at FROM todos`)

	if err != nil {
		fmt.Println(err)
	}

	var todos []models.Todo

	for statement.Next() {
		err = statement.Scan(&id, &item, &completed, &createdAt)

		if err != nil {
			fmt.Println(err)
		}

		todo := models.Todo{
			Id:        id,
			Item:      item,
			Completed: completed,
			CreatedAt: createdAt,
		}

		todos = append(todos, todo)
	}

	data := models.View{
		Todos: todos,
	}

	_ = view.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {

	item := r.FormValue("item")
	currentTime := time.Now()

	_, err := database.Exec(`INSERT INTO todos (item, created_at) VALUE (?, ?)`, item, currentTime)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec(`DELETE FROM todos WHERE id = ?`, id)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

func Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec(`UPDATE todos SET completed = 1 WHERE id = ?`, id)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

// UpdateTodo handles the PUT request to update an existing todo item
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Get the id parameter from request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID format"})
		return
	}

	// Parse request body
	var todoData struct {
		Item string `json:"item"`
	}

	// Decode JSON body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todoData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Get the existing todo to preserve the created_at value
	var existingTodo models.Todo
	err = database.QueryRow("SELECT id, item, completed, created_at FROM todos WHERE id = ?", id).Scan(
		&existingTodo.Id, 
		&existingTodo.Item, 
		&existingTodo.Completed,
		&existingTodo.CreatedAt,
	)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch existing todo"})
		fmt.Println(err)
		return
	}

	// Create todo object with updated data
	todo := models.Todo{
		Id:        id,
		Item:      todoData.Item,
		Completed: existingTodo.Completed,
		CreatedAt: existingTodo.CreatedAt, // Preserve the original creation time
	}

	// Update the todo in database
	err = models.UpdateTodo(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update todo"})
		fmt.Println(err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated successfully"})
}