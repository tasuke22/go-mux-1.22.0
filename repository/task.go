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
	GetTodoByID(ctx context.Context, id int, userId string) (*model.Todo, error)
	UpdateTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	DeleteTodo(ctx context.Context, todoId int, userId string) (*model.Todo, error)
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

func (tr *taskRepository) GetTodoByID(ctx context.Context, id int, userId string) (*model.Todo, error) {
	todo, err := model.Todos(model.TodoWhere.ID.EQ(id), model.TodoWhere.UserID.EQ(userId)).One(ctx, tr.db)
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

func (tr *taskRepository) DeleteTodo(ctx context.Context, todoId int, userId string) (*model.Todo, error) {
	// 最初に、指定されたIDとuserIdに一致するToDoをデータベースから検索します。
	todo, err := model.Todos(model.TodoWhere.ID.EQ(todoId), model.TodoWhere.UserID.EQ(userId)).One(ctx, tr.db)
	if err != nil {
		if err == sql.ErrNoRows {
			// 指定された条件に一致するToDoが存在しない場合
			return nil, fmt.Errorf("指定されたToDoが見つかりませんでした: %w", err)
		}
		// その他のエラーの場合
		return nil, fmt.Errorf("ToDoの検索に失敗しました: %w", err)
	}

	// ToDoをデータベースから削除します。
	if _, err := todo.Delete(ctx, tr.db); err != nil {
		return nil, fmt.Errorf("ToDoの削除に失敗しました: %w", err)
	}

	return todo, nil
}
