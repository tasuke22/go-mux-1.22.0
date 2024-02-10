package controller

import (
	"context"
	"encoding/json"
	"github.com/tasuke/go-mux/service"
	"net/http"
	"strconv"
)

type TaskController struct {
	ts *service.TaskService
}

func NewTaskController(ts *service.TaskService) *TaskController {
	return &TaskController{ts: ts}
}

func (tc *TaskController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todoRequest service.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	newTodo, err := tc.ts.CreateTodo(ctx, todoRequest)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTodo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (tc *TaskController) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	allTodos, err := tc.ts.GetAllTodos(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allTodos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tc *TaskController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	taskID, _ := strconv.Atoi(id)
	ctx := context.Background()
	todo, err := tc.ts.GetTodoByID(ctx, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (tc *TaskController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	taskID, _ := strconv.Atoi(id)
	ctx := context.Background()
	deleteTodo, err := tc.ts.DeleteTodo(ctx, taskID)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(deleteTodo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (tc *TaskController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var updateTodoRequest service.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&updateTodoRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// URLパラメータからToDoのIDを取得
	id := r.PathValue("id")
	taskId, _ := strconv.Atoi(id)

	ctx := context.Background()
	updateTodo, err := tc.ts.UpdateTodo(ctx, taskId, updateTodoRequest)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return

	}

	// 更新されたToDoをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updateTodo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
