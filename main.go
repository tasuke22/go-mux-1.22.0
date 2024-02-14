package main

import (
	"context"
	"github.com/tasuke/go-mux/config"
	"github.com/tasuke/go-mux/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db, err := config.NewDB()
	if err != nil {
		log.Fatalf("データベースへの接続に失敗しました: %v", err)
	}
	api := http.NewServeMux()
	router.InitRoute(api, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: api,
	}

	go func() {
		log.Printf("HTTPサーバをポート %s で起動しています...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPサーバの起動に失敗しました: %v", err)
		}
	}()

	// シグナルを待機
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// グレースフルシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("サーバのシャットダウン中にエラーが発生しました: %v", err)
	}
}
