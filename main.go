package main

import (
	"github.com/tasuke/go-mux/config"
	"github.com/tasuke/go-mux/controller"
	"github.com/tasuke/go-mux/repository"
	"github.com/tasuke/go-mux/router"
	"github.com/tasuke/go-mux/usecase"
	"log"
	"net/http"
)

func main() {
	// データベース接続
	db, err := config.NewDB()
	if err != nil {
		log.Fatalf("データベースへの接続に失敗しました: %v", err)
	}

	// リポジトリの初期化
	ur := repository.NewUserRepository(db)
	tr := repository.NewTaskRepository(db)

	// ユースケースの初期化
	uu := usecase.NewUserUsecase(ur)
	tu := usecase.NewTaskUsecase(tr)

	// コントローラの初期化
	uc := controller.NewUserController(uu)
	tc := controller.NewTaskController(tu)

	// ルータの初期化
	mux := router.NewRouter(uc, tc)

	// HTTPサーバの起動
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("HTTPサーバの起動に失敗しました: %v", err)
	}
}
