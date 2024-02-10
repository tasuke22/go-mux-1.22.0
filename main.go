package main

import (
	"fmt"
	"github.com/tasuke/go-mux/config"
	"github.com/tasuke/go-mux/controller"
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

	us := service.NewUserService(db)
	ts := service.NewTaskService(db)
	uc := controller.NewUserController(us)
	tc := controller.NewTaskController(ts)
	mux := router.NewRouter(uc, tc)

	http.ListenAndServe(":8080", mux)
}
