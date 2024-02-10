package usecase

import (
	"context"
	"github.com/tasuke/go-mux/model"
	"github.com/tasuke/go-mux/repository"
)

type TaskUsecase struct {
	tr repository.TaskRepository
}

func NewTaskUsecase(tr repository.TaskRepository) *TaskUsecase {
	return &TaskUsecase{tr: tr}
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (ts *TaskUsecase) CreateTodo(ctx context.Context, todoRequest CreateTodoRequest) (model.Todo, error) {
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

	newTodo, err := ts.tr.CreateTodo(ctx, newTodo)
	if err != nil {
		return model.Todo{}, err
	}
	return *newTodo, nil
}

func (ts *TaskUsecase) GetAllTodos(ctx context.Context) (model.TodoSlice, error) {
	todos, err := ts.tr.GetAllTodos(ctx)
	if err != nil {
		return model.TodoSlice{}, err
	}
	return todos, nil
}

func (ts *TaskUsecase) GetTodoByID(ctx context.Context, id int) (model.Todo, error) {
	todo, err := ts.tr.GetTodoByID(ctx, id)
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

func (ts *TaskUsecase) UpdateTodo(ctx context.Context, id int, updateTodoRequest UpdateTodoRequest) (model.Todo, error) {

	todo, err := ts.tr.GetTodoByID(ctx, id)
	if err != nil {
		return model.Todo{}, err
	}

	todo.Title = updateTodoRequest.Title
	todo.Description = updateTodoRequest.Description
	todo.Completed = updateTodoRequest.Completed

	todo, err = ts.tr.UpdateTodo(ctx, todo)
	if err != nil {
		return model.Todo{}, err
	}
	return *todo, nil
}

func (ts *TaskUsecase) DeleteTodo(ctx context.Context, id int) (model.Todo, error) {
	todo, err := ts.tr.GetTodoByID(ctx, id)
	if err != nil {
		return model.Todo{}, err
	}

	todo, err = ts.tr.DeleteTodo(ctx, todo)
	if err != nil {
		return model.Todo{}, err
	}
	return *todo, nil
}