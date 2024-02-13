package router

import (
	"github.com/tasuke/go-mux/controller"
	"github.com/tasuke/go-mux/middleware"
	"net/http"
)

func NewRouter(uc *controller.UserController, tc *controller.TaskController) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /api/v1/auth/signup", uc.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", uc.Login)
	mux.Handle("POST /api/v1/auth/logout", middleware.JWTMiddleware(uc.Logout))

	mux.Handle("POST /api/v1/todos", middleware.JWTMiddleware(tc.CreateTodo))
	mux.HandleFunc("GET /api/v1/todos", tc.GetAllTodos)
	mux.HandleFunc("GET /api/v1/todos/{id}", tc.GetTodoByID)
	mux.Handle("POST /api/v1/todos/{id}", middleware.JWTMiddleware(tc.UpdateTodo))
	mux.Handle("DELETE /api/v1/todos/{id}", middleware.JWTMiddleware(tc.DeleteTodo))
	return mux
}
