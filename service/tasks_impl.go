package service

import (
	"net/http"
	"todolist/helper"
	"todolist/model"
	"todolist/repository"

	"github.com/gin-gonic/gin"
)

type tasksServiceImpl struct {
	repository repository.TasksRepository
	cache      *helper.Cache
}

func NewTasksService(repo repository.TasksRepository, cache *helper.Cache) TasksService {
	return &tasksServiceImpl{
		repository: repo,
		cache:      cache,
	}
}

// InsertTask implements TasksService.
func (t *tasksServiceImpl) InsertTask(task *model.Task) (code int, response *gin.H) {
	id, err := t.repository.InsertTask(task)
	helper.PanicIfError(err)

	helper.Log.Info(gin.H{"id": id}, "Task created")
	return http.StatusCreated, &gin.H{
		"message": "Task created successfully",
		"task": model.Task{
			ID:          *id,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			DueDate:     task.DueDate,
		},
	}
}

// GetAllTasks implements TasksService.
func (t *tasksServiceImpl) GetAllTasks(status *string, page int, limit int, search *string) (code int, response *gin.H) {
	tasks, pagination, err := t.repository.GetAllTasks(status, page, limit, search)
	helper.PanicIfError(err)

	res := gin.H{
		"tasks": tasks,
	}

	if pagination != nil {
		res["pagination"] = pagination
	}

	return http.StatusOK, &res
}

// InsertTask implements TasksService.
func (t *tasksServiceImpl) UpdateTask(task *model.Task) (code int, response *gin.H) {
	currentTask, err := t.repository.GetTaskById(task.ID)
	helper.PanicIfError(err)

	if currentTask == nil {
		return http.StatusNotFound, &gin.H{
			"message": "Task not found",
		}
	}

	err = t.repository.UpdateTask(task)
	helper.PanicIfError(err)

	t.cache.DeleteTaskCache(task.ID)
	helper.Log.Info(gin.H{"id": task.ID}, "Task updated")
	return http.StatusOK, &gin.H{
		"message": "Task updated successfully",
		"task": model.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			DueDate:     task.DueDate,
		},
	}
}

// GetTask implements TasksService.
func (t *tasksServiceImpl) GetTask(id int) (code int, response *gin.H) {
	var currentTask *model.Task
	var err error

	currentTask = t.cache.GetTaskCache(id)

	if currentTask == nil {
		currentTask, err = t.repository.GetTaskById(id)
		helper.PanicIfError(err)
	} else {
		helper.Log.Info(gin.H{"id": id}, "Get task from cache")
	}

	if currentTask == nil {
		return http.StatusNotFound, &gin.H{
			"message": "Task not found",
		}
	}

	response, err = helper.TypeConverter[gin.H](currentTask)
	helper.PanicIfError(err)

	t.cache.SetTaskCache(id, currentTask)

	return http.StatusOK, response
}

func (t *tasksServiceImpl) DeleteTask(id int) (code int, response *gin.H) {
	currentTask, err := t.repository.GetTaskById(id)
	helper.PanicIfError(err)

	if currentTask == nil {
		return http.StatusNotFound, &gin.H{
			"message": "Task not found",
		}
	}

	t.cache.DeleteTaskCache(id)
	err = t.repository.DeleteTask(id)
	helper.PanicIfError(err)

	helper.Log.Info(gin.H{"id": id}, "Task deleted")
	return http.StatusOK, &gin.H{
		"message": "Task deleted successfully",
	}
}
