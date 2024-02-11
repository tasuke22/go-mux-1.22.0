package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func NewDB() (*sql.DB, error) {
	// .env ファイルから環境変数をロード（エラーがあってもプログラムは終了しない）
	err := godotenv.Load()
	if err != nil {
		fmt.Println("警告: .env ファイルのロードに失敗しました。環境変数が設定されていることを確認してください。")
	}

	// データベース接続用のDSNを構築
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"))

	// データベースに接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("データベース接続のオープンに失敗しました: %w", err)
	}

	// データベース接続を確認
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベースへのpingに失敗しました: %w", err)
	}

	fmt.Println("MySQLへの接続に成功しました")
	return db, nil
}
