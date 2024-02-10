package service

import (
	"context"
	"database/sql"
	"github.com/tasuke/go-mux/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TaskService struct {
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db}
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (ts *TaskService) CreateTodo(ctx context.Context, todoRequest CreateTodoRequest) (model.Todo, error) {
	// データベース接続を取得（dbは*sql.DB型のグローバル変数とします）
	// jwtからログインしているユーザーIDを取得 TODO
	userId := "7b31ecae-ab14-4f4a-a8b3-8b63cc10ddb2"

	// SQLBoilerを使用してToDoを作成
	newTodo := &model.Todo{
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		Completed:   todoRequest.Completed,
		UserID:      userId,
	}

	if err := newTodo.Insert(ctx, ts.db, boil.Infer()); err != nil {
		return model.Todo{}, err
	}
	return *newTodo, nil
}

func (ts *TaskService) GetAllTodos(ctx context.Context) (model.TodoSlice, error) {
	todos, err := model.Todos().All(ctx, ts.db)
	if err != nil {
		return model.TodoSlice{}, err
	}
	return todos, nil
}

func (ts *TaskService) GetTodoByID(ctx context.Context, id int) (model.Todo, error) {
	todo, err := model.Todos(qm.Where("id=?", id)).One(ctx, ts.db)
	if err != nil {
		return model.Todo{}, err
	}
	return *todo, nil
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (ts *TaskService) UpdateTodo(ctx context.Context, id int, updateTodoRequest UpdateTodoRequest) (model.Todo, error) {

	// SQLBoilerを使用してToDoを更新
	todo, err := model.FindTodo(ctx, ts.db, id)
	if err != nil {
		return model.Todo{}, err
	}

	todo.Title = updateTodoRequest.Title
	todo.Description = updateTodoRequest.Description
	todo.Completed = updateTodoRequest.Completed

	if _, err := todo.Update(ctx, ts.db, boil.Infer()); err != nil {
		return model.Todo{}, err
	}
	return *todo, nil
}

func (ts *TaskService) DeleteTodo(ctx context.Context, id int) (model.Todo, error) {
	todo, err := model.FindTodo(ctx, ts.db, id)
	if err != nil {
		return model.Todo{}, err
	}

	if _, err := todo.Delete(ctx, ts.db); err != nil {
		return model.Todo{}, err
	}
	return *todo, nil
}
