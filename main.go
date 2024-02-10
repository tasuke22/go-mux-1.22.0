package main

import (
	"fmt"
	"github.com/tasuke/go-mux/config"
	"github.com/tasuke/go-mux/controller"
	"github.com/tasuke/go-mux/repository"
	"github.com/tasuke/go-mux/router"
	"github.com/tasuke/go-mux/service"
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
	us := service.NewUserService(ur)
	ts := service.NewTaskService(tr)
	uc := controller.NewUserController(us)
	tc := controller.NewTaskController(ts)
	mux := router.NewRouter(uc, tc)

	http.ListenAndServe(":8080", mux)
}
