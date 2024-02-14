package router

import (
	"database/sql"
	"github.com/tasuke/go-mux/auth"
	"github.com/tasuke/go-mux/controller"
	"github.com/tasuke/go-mux/middleware"
	"github.com/tasuke/go-mux/repository"
	"github.com/tasuke/go-mux/usecase"
	"net/http"
)

func InitRoute(api *http.ServeMux, db *sql.DB) {
	api.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	{
		todoRoute(api, db)
		userRoute(api, db)
	}
}

func todoRoute(api *http.ServeMux, db *sql.DB) {
	tr := repository.NewTaskRepository(db)
	tu := usecase.NewTaskUsecase(tr)
	at := auth.NewAuthToken()
	tc := controller.NewTaskController(tu, at)

	api.Handle("POST /api/v1/todos", middleware.JWTMiddleware(tc.CreateTodo))
	api.HandleFunc("GET /api/v1/todos", tc.GetAllTodos)
	api.HandleFunc("GET /api/v1/todos/{id}", tc.GetTodoByID)
	api.Handle("POST /api/v1/todos/{id}", middleware.JWTMiddleware(tc.UpdateTodo))
	api.Handle("DELETE /api/v1/todos/{id}", middleware.JWTMiddleware(tc.DeleteTodo))
}

func userRoute(api *http.ServeMux, db *sql.DB) {
	ur := repository.NewUserRepository(db)
	uu := usecase.NewUserUsecase(ur)
	uc := controller.NewUserController(uu)

	api.HandleFunc("POST /api/v1/auth/signup", uc.SignUp)
	api.HandleFunc("POST /api/v1/auth/login", uc.Login)
	api.Handle("POST /api/v1/auth/logout", middleware.JWTMiddleware(uc.Logout))
}
