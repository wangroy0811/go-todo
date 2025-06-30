package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/controllers"
)

func Init() *mux.Router {
	route := mux.NewRouter()

	route.HandleFunc("/", controllers.Show)
	route.HandleFunc("/home", controllers.ShowHomePage)
	route.HandleFunc("/add", controllers.Add).Methods("POST")
	route.HandleFunc("/delete/{id}", controllers.Delete)
	route.HandleFunc("/complete/{id}", controllers.Complete)
	route.HandleFunc("/todos/{id}", controllers.UpdateTodo).Methods("PUT")
	route.HandleFunc("/blog", controllers.GetBlog)
	route.HandleFunc("/blog/{id}", controllers.GetBlogDetail)

	// Static file serving
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return route
}