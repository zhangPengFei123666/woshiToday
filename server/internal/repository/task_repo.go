package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"distributed-scheduler/internal/model"
	"distributed-scheduler/pkg/mysql"
)

// TaskRepository 任务仓库接口
type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Task, error)
	List(ctx context.Context, page, pageSize int, groupID uint64, keyword string, status int8) ([]*model.Task, int64, error)
	GetEnabledTasks(ctx context.Context) ([]*model.Task, error)
	GetTasksToTrigger(ctx context.Context, beforeTime time.Time, limit int) ([]*model.Task, error)
	UpdateNextTriggerTime(ctx context.Context, id uint64, nextTime time.Time, lastTime time.Time) error
	UpdateStatus(ctx context.Context, id uint64, status int8) error
	GetDependencies(ctx context.Context, taskID uint64) ([]model.Task, error)
	SetDependencies(ctx context.Context, taskID uint64, dependTaskIDs []uint64) error
}

// taskRepository 任务仓库实现
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓库
func NewTaskRepository() TaskRepository {
	return &taskRepository{db: mysql.GetDB()}
}

// Create 创建任务
func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// Update 更新任务
func (r *taskRepository) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete 删除任务(软删除)
func (r *taskRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}

// GetByID 根据ID获取任务
func (r *taskRepository) GetByID(ctx context.Context, id uint64) (*model.Task, error) {
	var task model.Task
	err := r.db.WithContext(ctx).Preload("Group").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List 获取任务列表
func (r *taskRepository) List(ctx context.Context, page, pageSize int, groupID uint64, keyword string, status int8) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Task{})

	if groupID > 0 {
		db = db.Where("group_id = ?", groupID)
	}
	if keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status >= 0 {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(mysql.Paginate(page, pageSize)).
		Preload("Group").
		Order("id DESC").
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetEnabledTasks 获取所有启用的任务
func (r *taskRepository) GetEnabledTasks(ctx context.Context) ([]*model.Task, error) {
	var tasks []*model.Task
	err := r.db.WithContext(ctx).
		Where("status = ?", model.TaskStatusEnabled).
		Find(&tasks).Error
	return tasks, err
}

// GetTasksToTrigger 获取需要触发的任务
func (r *taskRepository) GetTasksToTrigger(ctx context.Context, beforeTime time.Time, limit int) ([]*model.Task, error) {
	var tasks []*model.Task
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_trigger_time <= ?", model.TaskStatusEnabled, beforeTime).
		Order("next_trigger_time ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// UpdateNextTriggerTime 更新下次触发时间
func (r *taskRepository) UpdateNextTriggerTime(ctx context.Context, id uint64, nextTime time.Time, lastTime time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"next_trigger_time": nextTime,
			"last_trigger_time": lastTime,
		}).Error
}

// UpdateStatus 更新任务状态
func (r *taskRepository) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).
		Update("status", status).Error
}

// GetDependencies 获取任务依赖
func (r *taskRepository) GetDependencies(ctx context.Context, taskID uint64) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).
		Joins("JOIN task_dependency ON task_dependency.depend_task_id = task.id").
		Where("task_dependency.task_id = ?", taskID).
		Find(&tasks).Error
	return tasks, err
}

// SetDependencies 设置任务依赖
func (r *taskRepository) SetDependencies(ctx context.Context, taskID uint64, dependTaskIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除原有依赖
		if err := tx.Where("task_id = ?", taskID).Delete(&model.TaskDependency{}).Error; err != nil {
			return err
		}
		// 添加新依赖
		for _, depID := range dependTaskIDs {
			dep := model.TaskDependency{
				TaskID:       taskID,
				DependTaskID: depID,
			}
			if err := tx.Create(&dep).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// TaskGroupRepository 任务组仓库接口
type TaskGroupRepository interface {
	Create(ctx context.Context, group *model.TaskGroup) error
	Update(ctx context.Context, group *model.TaskGroup) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.TaskGroup, error)
	GetByAppName(ctx context.Context, appName string) (*model.TaskGroup, error)
	List(ctx context.Context, page, pageSize int, keyword string) ([]*model.TaskGroup, int64, error)
	GetAll(ctx context.Context) ([]*model.TaskGroup, error)
}

// taskGroupRepository 任务组仓库实现
type taskGroupRepository struct {
	db *gorm.DB
}

// NewTaskGroupRepository 创建任务组仓库
func NewTaskGroupRepository() TaskGroupRepository {
	return &taskGroupRepository{db: mysql.GetDB()}
}

// Create 创建任务组
func (r *taskGroupRepository) Create(ctx context.Context, group *model.TaskGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// Update 更新任务组
func (r *taskGroupRepository) Update(ctx context.Context, group *model.TaskGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

// Delete 删除任务组(软删除)
func (r *taskGroupRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.TaskGroup{}, id).Error
}

// GetByID 根据ID获取任务组
func (r *taskGroupRepository) GetByID(ctx context.Context, id uint64) (*model.TaskGroup, error) {
	var group model.TaskGroup
	err := r.db.WithContext(ctx).First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetByAppName 根据应用名获取任务组
func (r *taskGroupRepository) GetByAppName(ctx context.Context, appName string) (*model.TaskGroup, error) {
	var group model.TaskGroup
	err := r.db.WithContext(ctx).Where("app_name = ?", appName).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// List 获取任务组列表
func (r *taskGroupRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]*model.TaskGroup, int64, error) {
	var groups []*model.TaskGroup
	var total int64

	db := r.db.WithContext(ctx).Model(&model.TaskGroup{})

	if keyword != "" {
		db = db.Where("name LIKE ? OR app_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(mysql.Paginate(page, pageSize)).Order("id DESC").Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// GetAll 获取所有任务组
func (r *taskGroupRepository) GetAll(ctx context.Context) ([]*model.TaskGroup, error) {
	var groups []*model.TaskGroup
	err := r.db.WithContext(ctx).Where("status = ?", 1).Find(&groups).Error
	return groups, err
}

