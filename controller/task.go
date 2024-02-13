package controller

import (
	"encoding/json"
	"github.com/tasuke/go-mux/auth"
	"github.com/tasuke/go-mux/usecase"
	"net/http"
	"strconv"
)

type TaskController struct {
	tu *usecase.TaskUsecase
	at *auth.AuthToken
}

func NewTaskController(tu *usecase.TaskUsecase, at *auth.AuthToken) *TaskController {
	return &TaskController{tu: tu, at: at}
}

func (tc *TaskController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todoRequest usecase.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		sendErrorResponse(w, "リクエストボディが無効です。", http.StatusBadRequest)
		return
	}

	// ユーザーIDを取得
	userId, err := tc.at.ExtractUserIDFromToken(r)
	if err != nil {
		sendErrorResponse(w, "認証が必要です。", http.StatusUnauthorized)
		return
	}

	newTodo, err := tc.tu.CreateTodo(r.Context(), todoRequest, userId)
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

	// ユーザーIDを取得
	userId, err := tc.at.ExtractUserIDFromToken(r)
	if err != nil {
		sendErrorResponse(w, "認証が必要です。", http.StatusUnauthorized)
		return
	}

	todo, err := tc.tu.GetTodoByID(r.Context(), taskID, userId)
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

	// ユーザーIDを取得
	userId, err := tc.at.ExtractUserIDFromToken(r)
	if err != nil {
		sendErrorResponse(w, "認証が必要です。", http.StatusUnauthorized)
		return
	}

	deleteTodo, err := tc.tu.DeleteTodo(r.Context(), taskID, userId)
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

	// ユーザーIDを取得
	userId, err := tc.at.ExtractUserIDFromToken(r)
	if err != nil {
		sendErrorResponse(w, "認証が必要です。", http.StatusUnauthorized)
		return
	}

	updatedTodo, err := tc.tu.UpdateTodo(r.Context(), taskID, userId, updateTodoRequest)
	if err != nil {
		sendErrorResponse(w, "ToDoの更新に失敗しました。", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, updatedTodo, http.StatusOK)
}
