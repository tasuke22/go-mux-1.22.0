package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tasuke/go-mux/config"
	"github.com/tasuke/go-mux/controller"
	"github.com/tasuke/go-mux/model"
	"github.com/tasuke/go-mux/router"
	"github.com/tasuke/go-mux/service"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
)

func main() {
	db, err := config.NewDB()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return
	}

	us := service.NewUserService(db)
	//ts := service.NewTaskService(db)
	uc := controller.NewUserController(us)
	//tc := controller.NewTaskController(ts)
	mux := router.NewRouter(uc)

	//mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusOK)
	//})
	//
	//mux.HandleFunc("POST /api/v1/todos", createTodo)
	//mux.HandleFunc("GET /api/v1/todos", getAllTodos)
	//mux.HandleFunc("GET /api/v1/todos/{id}", getTodoByID)
	//mux.HandleFunc("POST /api/v1/todos/{id}", updateTodo)
	//mux.HandleFunc("DELETE /api/v1/todos/{id}", deleteTodo)

	http.ListenAndServe(":8080", mux)
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createTodo")
	// リクエストボディをCreateTodoRequest構造体にデコード
	var todoRequest CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// データベース接続を取得（dbは*sql.DB型のグローバル変数とします）
	db, err := config.NewDB() // 既にデータベース接続を確立している場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// context.Background()は、特別なコンテキストが必要ない場合に便利です。
	ctx := context.Background()
	// jwtからログインしているユーザーIDを取得 TODO
	userId := "8947f187-6f5c-4a99-a237-e5bd032cf543"

	// SQLBoilerを使用してToDoを作成
	newTodo := &model.Todo{
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		Completed:   todoRequest.Completed,
		UserID:      userId,
	}

	if err := newTodo.Insert(ctx, db, boil.Infer()); err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	// 作成されたToDoをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTodo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	// データベース接続を取得（dbは*sql.DB型のグローバル変数とします）
	db, err := config.NewDB() // 既にデータベース接続を確立している場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// URLパラメータからToDoのIDを取得
	id := r.PathValue("id")
	numId, _ := strconv.Atoi(id)
	fmt.Println(id)

	// context.Background()は、特別なコンテキストが必要ない場合に便利です。
	ctx := context.Background()

	// SQLBoilerを使用してToDoを削除
	todo, err := model.FindTodo(ctx, db, numId)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if _, err := todo.Delete(ctx, db); err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	// 削除されたToDoをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	// データベース接続を取得（dbは*sql.DB型のグローバル変数とします）
	db, err := config.NewDB() // 既にデータベース接続を確立している場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// リクエストボディをUpdateTodoRequest構造体にデコード
	var todoRequest UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// context.Background()は、特別なコンテキストが必要ない場合に便利です。
	ctx := context.Background()

	// URLパラメータからToDoのIDを取得
	id := r.PathValue("id")
	numId, _ := strconv.Atoi(id)
	fmt.Println(id)

	// SQLBoilerを使用してToDoを更新
	todo, err := model.FindTodo(ctx, db, numId)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todo.Title = todoRequest.Title
	todo.Description = todoRequest.Description
	todo.Completed = todoRequest.Completed

	if _, err := todo.Update(ctx, db, boil.Infer()); err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	// 更新されたToDoをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func getTodoByID(w http.ResponseWriter, r *http.Request) {
	// データベース接続を取得（dbは*sql.DB型のグローバル変数とします）
	fmt.Println("getTodoByID")
	db, err := config.NewDB() // 既にデータベース接続を確立している場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// URLパラメータからToDoのIDを取得
	id := r.PathValue("id")
	fmt.Println(id)

	// context.Background()は、特別なコンテキストが必要ない場合に便利です。
	// 実際のアプリケーションでは、リクエストスコープのコンテキストを使用することを検討してください。
	ctx := context.Background()

	// SQLBoilerを使用してToDoを取得
	todo, err := model.Todos(qm.Where("id=?", id)).One(ctx, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(todo)

	// 結果をJSONとしてエンコードしてレスポンスに書き込む
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getAllTodos(w http.ResponseWriter, r *http.Request) {
	// データベース接続を取得（dbは*sql.DB型のグローバル変数とします）
	db, err := config.NewDB() // 既にデータベース接続を確立している場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// context.Background()は、特別なコンテキストが必要ない場合に便利です。
	// 実際のアプリケーションでは、リクエストスコープのコンテキストを使用することを検討してください。
	ctx := context.Background()

	// SQLBoilerを使用してすべてのToDoを取得
	todos, err := model.Todos().All(ctx, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 結果をJSONとしてエンコードしてレスポンスに書き込む
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
