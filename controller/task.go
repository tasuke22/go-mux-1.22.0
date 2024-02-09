package controller

import "github.com/tasuke/go-mux/service"

type TaskController struct {
	ts *service.TaskService
}

func NewTaskController(ts *service.TaskService) *TaskController {
	return &TaskController{ts: ts}
}
