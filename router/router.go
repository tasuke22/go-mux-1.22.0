package router

import (
	"github.com/tasuke/go-mux/controller"
	"net/http"
)

func NewRouter(uc *controller.UserController, tc *controller.TaskController) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /api/v1/auth/signup", uc.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", uc.Login)
	mux.HandleFunc("POST /api/v1/auth/logout", uc.Logout)

	mux.HandleFunc("POST /api/v1/todos", tc.CreateTodo)
	mux.HandleFunc("GET /api/v1/todos", tc.GetAllTodos)
	mux.HandleFunc("GET /api/v1/todos/{id}", tc.GetTodoByID)
	mux.HandleFunc("POST /api/v1/todos/{id}", tc.UpdateTodo)
	mux.HandleFunc("DELETE /api/v1/todos/{id}", tc.DeleteTodo)
	return mux
}
