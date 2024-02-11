package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tasuke/go-mux/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TaskRepository interface {
	CreateTodo(ctx context.Context, newTodo *model.Todo) (*model.Todo, error)
	GetAllTodos(ctx context.Context) (model.TodoSlice, error)
	GetTodoByID(ctx context.Context, id int) (*model.Todo, error)
	UpdateTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	DeleteTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error)
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (tr *taskRepository) CreateTodo(ctx context.Context, newTodo *model.Todo) (*model.Todo, error) {
	if err := newTodo.Insert(ctx, tr.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("todoの作成に失敗しました: %w", err)
	}
	return newTodo, nil
}

func (tr *taskRepository) GetAllTodos(ctx context.Context) (model.TodoSlice, error) {
	todos, err := model.Todos().All(ctx, tr.db)
	if err != nil {
		return nil, fmt.Errorf("todoの取得に失敗しました: %w", err)
	}
	return todos, nil
}

func (tr *taskRepository) GetTodoByID(ctx context.Context, id int) (*model.Todo, error) {
	todo, err := model.FindTodo(ctx, tr.db, id)
	if err != nil {
		return nil, fmt.Errorf("IDによるtodoの取得に失敗しました: %w", err)
	}
	return todo, nil
}

func (tr *taskRepository) UpdateTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	_, err := todo.Update(ctx, tr.db, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("todoの更新に失敗しました: %w", err)
	}
	return todo, nil
}

func (tr *taskRepository) DeleteTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	_, err := todo.Delete(ctx, tr.db)
	if err != nil {
		return nil, fmt.Errorf("todoの削除に失敗しました: %w", err)
	}
	return todo, nil
}
