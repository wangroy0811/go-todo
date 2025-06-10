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
	id int
	item string
	completed int
	createdAt time.Time
	var updatedAt time.Time
	view = template.Must(template.ParseFiles("./views/index.html"))
	blogView = template.Must(template.ParseFiles("./views/blog.html"))
	blogDetailView = template.Must(template.ParseFiles("./views/blog_detail.html"))
	database = config.Database()
)

func Show(w http.ResponseWriter, r *http.Request) {
	statement, err := database.Query(`SELECT id, item, completed, created_at, updated_at FROM todos`)

	if err != nil {
		fmt.Println(err)
	}

	var todos []models.Todo

	for statement.Next() {
		err = statement.Scan(&id, &item, &completed, &createdAt, &updatedAt)

		if err != nil {
			fmt.Println(err)
		}

		todo := models.Todo{
			Id: id,
			Item: item,
			Completed: completed,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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

	_, err := database.Exec(`INSERT INTO todos (item, created_at, updated_at) VALUE (?, ?, ?)`, item, currentTime, currentTime)

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
	currentTime := time.Now()

	_, err := database.Exec(`UPDATE todos SET completed = 1, updated_at = ? WHERE id = ?`, currentTime, id)

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
	err = database.QueryRow("SELECT id, item, completed, created_at, updated_at FROM todos WHERE id = ?", id).Scan(
		&existingTodo.Id,
		&existingTodo.Item,
		&existingTodo.Completed,
		&existingTodo.CreatedAt,
		&existingTodo.UpdatedAt,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch existing todo"})
		fmt.Println(err)
		return
	}

	// Create todo object with updated data
	todo := models.Todo{
		Id: id,
		Item: todoData.Item,
		Completed: existingTodo.Completed,
		CreatedAt: existingTodo.CreatedAt, // Preserve the original creation time
		UpdatedAt: time.Now(), // Set updated_at to current time
	}
	fmt.Println("Setting updated time to:", todo.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Update the todo in database
	fmt.Println("Sending to model layer with updated_at:", todo.UpdatedAt.Format("2006-01-02 15:04:05"))
	err = models.UpdateTodo(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update todo"})
		fmt.Println(err)
		return
	}
	fmt.Println("Todo with ID", id, "updated successfully with new updated_at time")

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated successfully"})
}

// ShowHomePage handles displaying the application homepage
func ShowHomePage(w http.ResponseWriter, r *http.Request) {
	// Create template for homepage
	homepageView := template.Must(template.ParseFiles("./views/homepage.html"))

	// Execute the template
	err := homepageView.Execute(w, nil)

	if err != nil {
		fmt.Println(err)
	}
}

// GetBlog handles displaying the blog page with all blog posts
func GetBlog(w http.ResponseWriter, r *http.Request) {
	// Get all blogs from the database
	blogs, err := models.GetAllBlogs()

	if err != nil {
		fmt.Println("Error retrieving blogs:", err)
		http.Error(w, "Failed to retrieve blog posts", http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := struct {
		Blogs []models.Blog
	}{
		Blogs: blogs,
	}

	// Execute the blog template
	err = blogView.Execute(w, data)

	if err != nil {
		fmt.Println("Error rendering blog template:", err)
		http.Error(w, "Failed to render blog page", http.StatusInternalServerError)
	}
}

// GetBlogDetail handles displaying a single blog post's details
func GetBlogDetail(w http.ResponseWriter, r *http.Request) {
	// Get the blog ID from URL parameters
	vars := mux.Vars(r)
	idParam := vars["id"]

	// Convert string ID to integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println("Invalid blog ID:", err)
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	// Get blog post by ID from the database
	blog, err := models.GetBlog(id)
	if err != nil {
		fmt.Println("Error retrieving blog with ID", id, ":", err)
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}

	// Prepare data for template
	data := struct {
		Blog models.Blog
	}{
		Blog: blog,
	}

	// Execute the blog detail template
	err = blogDetailView.Execute(w, data)

	if err != nil {
		fmt.Println("Error rendering blog detail template:", err)
		http.Error(w, "Failed to render blog detail page", http.StatusInternalServerError)
	}
}
