package routes

import (
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/controllers"
	"net/http"
)

func Init() *mux.Router {
	route := mux.NewRouter()

	// Todo routes
	route.HandleFunc("/", controllers.Show) // 支持分页查询参数: ?page=1
	route.HandleFunc("/home", controllers.ShowHomePage)
	route.HandleFunc("/add", controllers.Add).Methods("POST")
	route.HandleFunc("/delete/{id}", controllers.Delete)
	route.HandleFunc("/complete/{id}", controllers.Complete)
	route.HandleFunc("/todos/{id}", controllers.UpdateTodo).Methods("PUT")

	// Blog routes
	route.HandleFunc("/blog", controllers.GetBlog)
	route.HandleFunc("/blog/{id}", controllers.GetBlogDetail)

	// Static file serving
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return route
}
