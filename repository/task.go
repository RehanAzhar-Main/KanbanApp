package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	var listTask []entity.Task

	// get task by user id
	rows, err := r.db.WithContext(ctx).Table("tasks").
		Where("user_id = ?", id).Rows()

	if err != nil {
		return []entity.Task{}, err
	}

	defer rows.Close()
	for rows.Next() {
		r.db.ScanRows(rows, &listTask)
	}

	// task not found
	if len(listTask) == 0 {
		return []entity.Task{}, nil
	}

	return listTask, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	if err := r.db.WithContext(ctx).Table("tasks").
		Create(&task).Error; err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	// var task = &entity.Task{}
	// var storeTask entity.Task

	// // get task by id
	// if err := r.db.WithContext(ctx).Table("tasks").
	// 	Where("id = ?", id).
	// 	Find(task).Scan(&storeTask).Error; err != nil {
	// 	return entity.Task{}, err
	// }

	// return storeTask, nil

	task := entity.Task{}
	// err := r.db.First(&entity.Task{}, id).Scan(&task).Error
	err := r.db.Raw("SELECT * FROM tasks WHERE id = ?", id).Scan(&task).Error

	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	var listTask []entity.Task

	// get categories by user id
	rows, err := r.db.WithContext(ctx).Table("tasks").
		Joins("left join categories on categories.id = tasks.category_id").
		Where("deleted_at is null").Rows()

	if err != nil {
		return []entity.Task{}, err
	}

	defer rows.Close()
	for rows.Next() {
		r.db.ScanRows(rows, &listTask)
	}

	// categories not found
	if len(listTask) == 0 {
		return []entity.Task{}, nil
	}

	return listTask, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	if err := r.db.WithContext(ctx).Table("tasks").Where("id = ?", task.ID).
		Updates(&task).Error; err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	var task entity.Task
	if err := r.db.WithContext(ctx).
		Delete(&task, id).Error; err != nil {
		return err
	}

	return nil
}
