package controller

import (
	"encoding/json"
	"github.com/tasuke/go-mux/usecase"
	"net/http"
	"strconv"
)

type TaskController struct {
	ts *usecase.TaskUsecase
}

func NewTaskController(ts *usecase.TaskUsecase) *TaskController {
	return &TaskController{ts: ts}
}

func (tc *TaskController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todoRequest usecase.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newTodo, err := tc.ts.CreateTodo(r.Context(), todoRequest)
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
	allTodos, err := tc.ts.GetAllTodos(r.Context())
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
	todo, err := tc.ts.GetTodoByID(r.Context(), taskID)
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
	deleteTodo, err := tc.ts.DeleteTodo(r.Context(), taskID)
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
	var updateTodoRequest usecase.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&updateTodoRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// URLパラメータからToDoのIDを取得
	id := r.PathValue("id")
	taskId, _ := strconv.Atoi(id)

	updateTodo, err := tc.ts.UpdateTodo(r.Context(), taskId, updateTodoRequest)
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
