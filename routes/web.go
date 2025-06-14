package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/controllers"
)

func Init() *mux.Router {
	route := mux.NewRouter()

	// CORS middleware
	route.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	// Static file server
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	route.HandleFunc("/", controllers.Show)
	route.HandleFunc("/home", controllers.ShowHomePage)
	route.HandleFunc("/add", controllers.Add).Methods("POST")
	route.HandleFunc("/delete/{id}", controllers.Delete)
	route.HandleFunc("/complete/{id}", controllers.Complete)
	route.HandleFunc("/todos/{id}", controllers.UpdateTodo).Methods("PUT")
	route.HandleFunc("/blog", controllers.GetBlog)
	route.HandleFunc("/blog/{id}", controllers.GetBlogDetail)

	return route
}