package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/tasuke/go-mux/config"
	"github.com/tasuke/go-mux/model"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
)

var db *sql.DB

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

	mux.HandleFunc("POST /api/v1/auth/signup", signUp)
	mux.HandleFunc("POST /api/v1/auth/login", login)

	mux.HandleFunc("POST /api/v1/todos", createTodo)
	mux.HandleFunc("GET /api/v1/todos", getAllTodos)
	mux.HandleFunc("GET /api/v1/todos/{id}", getTodoByID)
	mux.HandleFunc("POST /api/v1/todos/{id}", updateTodo)
	mux.HandleFunc("DELETE /api/v1/todos/{id}", deleteTodo)

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
	newTodo := &model.Task{
		Title:       todoRequest.Title,
		Description: null.NewString(todoRequest.Description, true), // Descriptionがstring型の場合
		Completed:   null.NewBool(todoRequest.Completed, todoRequest.Completed),
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
	todo, err := model.FindTask(ctx, db, numId)
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
	todo, err := model.FindTask(ctx, db, numId)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todo.Title = todoRequest.Title
	todo.Description = null.NewString(todoRequest.Description, true) // Descriptionがstring型の場合
	todo.Completed = null.NewBool(todoRequest.Completed, todoRequest.Completed)

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
	todo, err := model.Tasks(qm.Where("id=?", id)).One(ctx, db)
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
	todos, err := model.Tasks().All(ctx, db)
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	db, _ := config.NewDB()
	currentUser, err := fetchUserByEmail(db, loginRequest.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// JWTトークンの生成準備
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": currentUser.ID,
		"sub":     currentUser.ID,
	})
	// JWTトークンの生成(署名)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// トークンの返却
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: tokenString}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// fetchUserByEmail は与えられたメールアドレスに一致するユーザーをデータベースから検索します。
func fetchUserByEmail(db *sql.DB, email string) (*model.User, error) {
	// context.Background() は、実際のアプリケーションでは適切なコンテキストに置き換えるべきです。
	user, err := model.Users(qm.Where("email=?", email)).One(context.Background(), db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var userRequest SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// データベース操作で使用され、その操作によってオブジェクトが変更される可能性があるため
	newUser := &model.User{
		ID:       uuid.New().String(),
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashedPassword),
	}

	db, _ := config.NewDB()
	ctx := r.Context()
	err = newUser.Insert(ctx, db, boil.Infer())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 値でOK。小さい構造体、レスポンスするためだけの構造体, 構造体を変更する必要がない,=> コードの意図がより明確
	signUpResponse := SignUpResponse{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	// レスポンスを設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created ステータスコードを返す
	if err := json.NewEncoder(w).Encode(signUpResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
