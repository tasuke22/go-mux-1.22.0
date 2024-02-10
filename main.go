package main

import (
	"fmt"
	"github.com/tasuke/go-mux/config"
	"github.com/tasuke/go-mux/controller"
	"github.com/tasuke/go-mux/repository"
	"github.com/tasuke/go-mux/router"
	"github.com/tasuke/go-mux/usecase"
	"net/http"
)

func main() {
	db, err := config.NewDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return
	}

	ur := repository.NewUserRepository(db)
	tr := repository.NewTaskRepository(db)
	uu := usecase.NewUserUsecase(ur)
	tu := usecase.NewTaskUsecase(tr)
	uc := controller.NewUserController(uu)
	tc := controller.NewTaskController(tu)
	mux := router.NewRouter(uc, tc)

	http.ListenAndServe(":8080", mux)
}
