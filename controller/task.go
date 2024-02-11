package controller

import (
	"encoding/json"
	"github.com/tasuke/go-mux/usecase"
	"net/http"
	"strconv"
)

type TaskController struct {
	tu *usecase.TaskUsecase
}

func NewTaskController(tu *usecase.TaskUsecase) *TaskController {
	return &TaskController{tu: tu}
}

func (tc *TaskController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todoRequest usecase.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		sendErrorResponse(w, "リクエストボディが無効です。", http.StatusBadRequest)
		return
	}

	newTodo, err := tc.tu.CreateTodo(r.Context(), todoRequest)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, newTodo, http.StatusCreated)
}

func (tc *TaskController) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	allTodos, err := tc.tu.GetAllTodos(r.Context())
	if err != nil {
		sendErrorResponse(w, "ToDoの取得に失敗しました。", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, allTodos, http.StatusOK)
}

func (tc *TaskController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendErrorResponse(w, "ToDoのIDが無効です。", http.StatusBadRequest)
		return
	}

	todo, err := tc.tu.GetTodoByID(r.Context(), taskID)
	if err != nil {
		sendErrorResponse(w, "指定されたToDoの取得に失敗しました。", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, todo, http.StatusOK)
}

func (tc *TaskController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendErrorResponse(w, "ToDoのIDが無効です。", http.StatusBadRequest)
		return
	}

	deleteTodo, err := tc.tu.DeleteTodo(r.Context(), taskID)
	if err != nil {
		sendErrorResponse(w, "ToDoの削除に失敗しました。", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, deleteTodo, http.StatusOK)
}

func (tc *TaskController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendErrorResponse(w, "ToDoのIDが無効です。", http.StatusBadRequest)
		return
	}

	var updateTodoRequest usecase.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&updateTodoRequest); err != nil {
		sendErrorResponse(w, "リクエストボディが無効です。", http.StatusBadRequest)
		return
	}

	updatedTodo, err := tc.tu.UpdateTodo(r.Context(), taskID, updateTodoRequest)
	if err != nil {
		sendErrorResponse(w, "ToDoの更新に失敗しました。", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, updatedTodo, http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sendErrorResponse(w, "レスポンスのエンコードに失敗しました。", http.StatusInternalServerError)
	}
}
