package main

import (
	"fmt"
	"github.com/tasuke/go-mux/config"
)

func main() {
	fmt.Println("Hello, World!")
	db, err := config.NewDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return
	}
	fmt.Println(db)
}
