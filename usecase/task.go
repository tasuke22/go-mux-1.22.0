package usecase

import (
	"context"
	"fmt"
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

func (ts *TaskUsecase) CreateTodo(ctx context.Context, todoRequest CreateTodoRequest, userId string) (model.Todo, error) {
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
		return model.TodoSlice{}, fmt.Errorf("toDoの取得に失敗しました: %w", err)
	}
	return todos, nil
}

func (ts *TaskUsecase) GetTodoByID(ctx context.Context, id int) (model.Todo, error) {
	todo, err := ts.tr.GetTodoByID(ctx, id)
	if err != nil {
		return model.Todo{}, fmt.Errorf("IDによるtodoの取得に失敗しました: %w", err)
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
		return model.Todo{}, fmt.Errorf("更新するtodoの取得に失敗しました: %w", err)
	}

	todo.Title = updateTodoRequest.Title
	todo.Description = updateTodoRequest.Description
	todo.Completed = updateTodoRequest.Completed

	updatedTodo, err := ts.tr.UpdateTodo(ctx, todo)
	if err != nil {
		return model.Todo{}, fmt.Errorf("todoの更新に失敗しました: %w", err)
	}
	return *updatedTodo, nil
}

func (ts *TaskUsecase) DeleteTodo(ctx context.Context, id int) (*model.Todo, error) {
	todo, err := ts.tr.GetTodoByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("削除するtodoの取得に失敗しました: %w", err)
	}

	deletedTodo, err := ts.tr.DeleteTodo(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("todoの削除に失敗しました: %w", err)
	}
	return deletedTodo, nil
}
