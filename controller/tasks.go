package controller

import "github.com/gin-gonic/gin"

type TaskController interface {
	InsertTask(ctx *gin.Context)
	GetAllTasks(ctx *gin.Context)
	UpdateTask(ctx *gin.Context)
	GetTask(ctx *gin.Context)
	DeleteTask(ctx *gin.Context)
}
