package router

import (
	"github.com/tasuke/go-mux/controller"
	"net/http"
)

func NewRouter(uc *controller.UserController) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/signup", uc.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", uc.Login)

	//mux.HandleFunc("POST /api/v1/todos", tc.createTodo)
	//mux.HandleFunc("GET /api/v1/todos", tc.getAllTodos)
	//mux.HandleFunc("GET /api/v1/todos/{id}", tc.getTodoByID)
	//mux.HandleFunc("POST /api/v1/todos/{id}", tc.updateTodo)
	//mux.HandleFunc("DELETE /api/v1/todos/{id}", tc.deleteTodo)
	return mux
}
