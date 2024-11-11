package controller

import (
	"net/http"
	"strconv"
	"time"
	"todolist/helper"
	"todolist/model"
	"todolist/service"

	"github.com/gin-gonic/gin"
)

type tasksControllerImpl struct {
	service service.TasksService
}

func NewTasksController(service service.TasksService) TaskController {
	return &tasksControllerImpl{
		service: service,
	}
}

// InsertTask implements TaskController.
func (t *tasksControllerImpl) InsertTask(ctx *gin.Context) {
	var payload model.Task

	if err := ctx.ShouldBindJSON(&payload); err == nil {
		_, err = time.Parse("2006-01-02", payload.DueDate)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

			return
		}

		code, response := t.service.InsertTask(&payload)

		ctx.JSON(code, response)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

// GetAllTasks implements TaskController.
func (t *tasksControllerImpl) GetAllTasks(ctx *gin.Context) {
	status, page, limit, search := ctx.Query("status"), ctx.Query("page"), ctx.Query("limit"), ctx.Query("search")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	code, response := t.service.GetAllTasks(&status, pageInt, limitInt, &search)

	ctx.JSON(code, response)
}

// InsertTask implements TaskController.
func (t *tasksControllerImpl) UpdateTask(ctx *gin.Context) {
	var payload model.Task

	if err := ctx.ShouldBindJSON(&payload); err == nil {
		_, err = time.Parse("2006-01-02", payload.DueDate)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

			return
		}

		id := ctx.Param("id")
		payload.ID, err = strconv.Atoi(id)
		helper.PanicIfError(err)

		code, response := t.service.UpdateTask(&payload)

		ctx.JSON(code, response)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func (t *tasksControllerImpl) GetTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.PanicIfError(err)

	code, response := t.service.GetTask(id)

	ctx.JSON(code, response)
}

func (t *tasksControllerImpl) DeleteTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helper.PanicIfError(err)

	code, response := t.service.DeleteTask(id)

	ctx.JSON(code, response)
}
