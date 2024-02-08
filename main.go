package main

import (
	"fmt"
	"github.com/tasuke/go-mux/config"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")
	db, err := config.NewDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return
	}
	fmt.Println(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8080", mux)

}
