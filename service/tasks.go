package service

import (
	"todolist/model"

	"github.com/gin-gonic/gin"
)

type TasksService interface {
	InsertTask(task *model.Task) (code int, response *gin.H)
	GetAllTasks(status *string, page, limit int, search *string) (code int, response *gin.H)
	UpdateTask(task *model.Task) (code int, response *gin.H)
	GetTask(id int) (code int, response *gin.H)
	DeleteTask(id int) (code int, response *gin.H)
}
