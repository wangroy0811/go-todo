package models

import (
	"fmt"
	"github.com/ichtrojan/go-todo/config"
	"time"
)

// Blog represents a blog post in the database
type Blog struct {
	ID        int
	Title     string
	Content   string
	Author    string
	CreatedAt time.Time
	Image     string
}

// FormatCreatedAt formats and returns the blog creation time as a string
func (blog Blog) FormatCreatedAt() string {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// Fallback to UTC if location loading fails
		return blog.CreatedAt.Format("2006-01-02 15:04:05")
	}
	// Convert the time to China timezone (UTC+8)
	return blog.CreatedAt.In(loc).Format("2006-01-02 15:04:05")
}

// GetAllBlogs retrieves all blog posts from the database
func GetAllBlogs() ([]Blog, error) {
	var blogs []Blog
	database := config.Database()

	statement, err := database.Query(`SELECT id, title, content, author, created_at, image FROM blog`)
	if err != nil {
		fmt.Println("Error retrieving blogs:", err)
		return nil, err
	}
	defer statement.Close()

	for statement.Next() {
		var blog Blog
		var id int
		var title, content, author, image string
		var createdAt time.Time

		err = statement.Scan(&id, &title, &content, &author, &createdAt, &image)
		if err != nil {
			fmt.Println("Error scanning blog row:", err)
			continue
		}

		blog = Blog{
			ID:        id,
			Title:     title,
			Content:   content,
			Author:    author,
			CreatedAt: createdAt,
			Image:     image,
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// GetBlog retrieves a single blog post by ID from the database
func GetBlog(id int) (Blog, error) {
	database := config.Database()
	var blog Blog

	err := database.QueryRow(`SELECT id, title, content, author, created_at, image FROM blog WHERE id = ?`, id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.Author,
		&blog.CreatedAt,
		&blog.Image,
	)

	if err != nil {
		fmt.Println("Error retrieving blog with ID", id, ":", err)
		return blog, err
	}

	return blog, nil
}

// CreateBlog adds a new blog post to the database
func CreateBlog(blog Blog) error {
	database := config.Database()
	currentTime := time.Now()

	_, err := database.Exec(
		`INSERT INTO blog (title, content, author, created_at, image) VALUES (?, ?, ?, ?, ?)`,
		blog.Title,
		blog.Content,
		blog.Author,
		currentTime,
		blog.Image,
	)

	if err != nil {
		fmt.Println("Error creating blog:", err)
		return err
	}

	return nil
}

// UpdateBlog updates an existing blog post in the database
func UpdateBlog(blog Blog) error {
	database := config.Database()

	_, err := database.Exec(
		`UPDATE blog SET title = ?, content = ?, author = ?, image = ? WHERE id = ?`,
		blog.Title,
		blog.Content,
		blog.Author,
		blog.Image,
		blog.ID,
	)

	if err != nil {
		fmt.Println("Error updating blog:", err)
		return err
	}

	return nil
}

// DeleteBlog removes a blog post from the database by ID
func DeleteBlog(id int) error {
	database := config.Database()

	_, err := database.Exec(`DELETE FROM blog WHERE id = ?`, id)
	if err != nil {
		fmt.Println("Error deleting blog:", err)
		return err
	}

	return nil
}

// InitBlogTable creates the blog table if it doesn't exist
func InitBlogTable() error {
	database := config.Database()

	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS blog (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author VARCHAR(100) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			image VARCHAR(255)
		);
	`)

	if err != nil {
		fmt.Println("Error creating blog table:", err)
		return err
	}

	return nil
}