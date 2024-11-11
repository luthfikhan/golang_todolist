package repository

import "todolist/model"

type TasksRepository interface {
	InsertTask(task *model.Task) (*int, error)
	GetTaskById(id int) (*model.Task, error)
	GetAllTasks(status *string, page, limit int, search *string) (*[]model.Task, *model.Pagination, error)
	UpdateTask(task *model.Task) error
	DeleteTask(id int) error
}
