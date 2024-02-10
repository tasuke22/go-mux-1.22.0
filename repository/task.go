package repository

import (
	"context"
	"database/sql"
	"github.com/tasuke/go-mux/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (tr *TaskRepository) CreateTodo(ctx context.Context, newTodo *model.Todo) (*model.Todo, error) {
	if err := newTodo.Insert(ctx, tr.db, boil.Infer()); err != nil {
		return nil, err
	}
	return newTodo, nil
}

func (tr *TaskRepository) GetAllTodos(ctx context.Context) (model.TodoSlice, error) {
	todos, err := model.Todos().All(ctx, tr.db)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (tr *TaskRepository) GetTodoByID(ctx context.Context, id int) (*model.Todo, error) {
	todo, err := model.FindTodo(ctx, tr.db, id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (tr *TaskRepository) UpdateTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	_, err := todo.Update(ctx, tr.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (tr *TaskRepository) DeleteTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	_, err := todo.Delete(ctx, tr.db)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (tr *TaskRepository) FindTodoById(ctx context.Context, id int) (*model.Todo, error) {
	todo, err := model.FindTodo(ctx, tr.db, id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}
