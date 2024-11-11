package repository

import (
	"database/sql"
	"math"
	"strconv"
	"strings"
	"todolist/helper"
	"todolist/model"

	"github.com/gin-gonic/gin"
)

type tasksRepositoryImpl struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) TasksRepository {
	return &tasksRepositoryImpl{
		db: db,
	}
}

// GetTaskById implements TasksRepository.
func (t *tasksRepositoryImpl) GetTaskById(id int) (*model.Task, error) {
	var task model.Task

	err := t.db.QueryRow(
		"SELECT id, title, description, status, due_date from tasks WHERE id = :1", id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDate)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// InsertTask implements TasksRepository.
func (t *tasksRepositoryImpl) InsertTask(task *model.Task) (*int, error) {
	var id int

	_, err := t.db.Exec(
		"INSERT INTO tasks (title, description, status, due_date) VALUES (:1, :2, :3, TO_DATE(:4, 'YYYY-MM-DD')) RETURNING id INTO :5",
		task.Title, task.Description, task.Status, task.DueDate, &id,
	)

	if err != nil {
		return nil, err
	}

	helper.Log.Error(gin.H{"title": task.Title, "id": id}, "Success insert task")

	return &id, nil
}
func (t *tasksRepositoryImpl) GetAllTasks(status *string, page, limit int, search *string) (*[]model.Task, *model.Pagination, error) {
	tasks := []model.Task{}
	var pagination *model.Pagination = nil

	countQuery := "SELECT COUNT(*) FROM tasks"
	query := "SELECT id, title, description, status, due_date FROM tasks"
	var conditions []string
	conditionCount := 0

	args := []interface{}{}

	if *status != "" {
		conditionCount++
		conditions = append(conditions, "status = :"+strconv.Itoa(conditionCount))
		args = append(args, status)
	}

	if *search != "" {
		conditionCount++
		conditions = append(conditions, "(LOWER(title) LIKE '%' || LOWER("+":"+strconv.Itoa(conditionCount)+") || '%' OR LOWER(description) LIKE '%' || LOWER("+":"+strconv.Itoa(conditionCount)+") || '%')")
		args = append(args, search)
	}

	if conditionCount > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	if page > 0 && limit > 0 {
		conditionCount++
		offset := 0
		totalTasks := 0

		err := t.db.QueryRow(countQuery, args...).Scan(&totalTasks)
		if err != nil {
			return nil, nil, err
		}

		offset = (page - 1) * limit
		query += " OFFSET :" + strconv.Itoa(conditionCount) + " ROWS FETCH NEXT :" + strconv.Itoa(conditionCount+1) + " ROWS ONLY"
		args = append(args, offset, limit)

		pagination = &model.Pagination{
			CurrentPage: page,
			TotalPages:  int(math.Ceil(float64(totalTasks) / float64(limit))),
			TotalTasks:  totalTasks,
		}
	}

	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDate); err != nil {
			return nil, nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return &tasks, pagination, nil
}

func (t *tasksRepositoryImpl) UpdateTask(task *model.Task) error {
	_, err := t.db.Exec(
		"UPDATE tasks SET title = :1, description = :2, status = :3, due_date = TO_DATE(:4, 'YYYY-MM-DD') WHERE id = :5",
		task.Title, task.Description, task.Status, task.DueDate, task.ID,
	)

	return err
}

func (t *tasksRepositoryImpl) DeleteTask(id int) error {
	_, err := t.db.Exec(
		"DELETE FROM tasks WHERE id = :1",
		id,
	)

	return err
}
