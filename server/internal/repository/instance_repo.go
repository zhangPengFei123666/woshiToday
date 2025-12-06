package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"distributed-scheduler/internal/model"
	"distributed-scheduler/pkg/mysql"
)

// InstanceRepository 任务实例仓库接口
type InstanceRepository interface {
	Create(ctx context.Context, instance *model.TaskInstance) error
	Update(ctx context.Context, instance *model.TaskInstance) error
	GetByID(ctx context.Context, id uint64) (*model.TaskInstance, error)
	List(ctx context.Context, page, pageSize int, taskID uint64, status int8, startTime, endTime *time.Time) ([]*model.TaskInstance, int64, error)
	UpdateStatus(ctx context.Context, id uint64, status int8, resultCode int, resultMsg string) error
	UpdateStartTime(ctx context.Context, id uint64) error
	UpdateEndTime(ctx context.Context, id uint64, status int8, resultCode int, resultMsg string) error
	GetRunningInstances(ctx context.Context, taskID uint64) ([]*model.TaskInstance, error)
	GetInstancesByTriggerTime(ctx context.Context, taskID uint64, triggerTime time.Time) ([]*model.TaskInstance, error)
	CountByStatus(ctx context.Context, taskID uint64, startTime, endTime time.Time) (map[int8]int64, error)
	GetRecentInstances(ctx context.Context, limit int) ([]*model.TaskInstance, error)
}

// instanceRepository 任务实例仓库实现
type instanceRepository struct {
	db *gorm.DB
}

// NewInstanceRepository 创建任务实例仓库
func NewInstanceRepository() InstanceRepository {
	return &instanceRepository{db: mysql.GetDB()}
}

// Create 创建任务实例
func (r *instanceRepository) Create(ctx context.Context, instance *model.TaskInstance) error {
	return r.db.WithContext(ctx).Create(instance).Error
}

// Update 更新任务实例
func (r *instanceRepository) Update(ctx context.Context, instance *model.TaskInstance) error {
	return r.db.WithContext(ctx).Save(instance).Error
}

// GetByID 根据ID获取任务实例
func (r *instanceRepository) GetByID(ctx context.Context, id uint64) (*model.TaskInstance, error) {
	var instance model.TaskInstance
	err := r.db.WithContext(ctx).Preload("Task").First(&instance, id).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// List 获取任务实例列表
func (r *instanceRepository) List(ctx context.Context, page, pageSize int, taskID uint64, status int8, startTime, endTime *time.Time) ([]*model.TaskInstance, int64, error) {
	var instances []*model.TaskInstance
	var total int64

	db := r.db.WithContext(ctx).Model(&model.TaskInstance{})

	if taskID > 0 {
		db = db.Where("task_id = ?", taskID)
	}
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	if startTime != nil {
		db = db.Where("trigger_time >= ?", startTime)
	}
	if endTime != nil {
		db = db.Where("trigger_time <= ?", endTime)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(mysql.Paginate(page, pageSize)).
		Preload("Task").
		Order("id DESC").
		Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	return instances, total, nil
}

// UpdateStatus 更新任务实例状态
func (r *instanceRepository) UpdateStatus(ctx context.Context, id uint64, status int8, resultCode int, resultMsg string) error {
	updates := map[string]interface{}{
		"status":      status,
		"result_code": resultCode,
		"result_msg":  resultMsg,
	}
	return r.db.WithContext(ctx).Model(&model.TaskInstance{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateStartTime 更新开始时间
func (r *instanceRepository) UpdateStartTime(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.TaskInstance{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     model.InstanceStatusRunning,
			"start_time": gorm.Expr("NOW()"),
		}).Error
}

// UpdateEndTime 更新结束时间
func (r *instanceRepository) UpdateEndTime(ctx context.Context, id uint64, status int8, resultCode int, resultMsg string) error {
	return r.db.WithContext(ctx).Model(&model.TaskInstance{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"result_code": resultCode,
			"result_msg":  resultMsg,
			"end_time":    gorm.Expr("NOW()"),
		}).Error
}

// GetRunningInstances 获取运行中的实例
func (r *instanceRepository) GetRunningInstances(ctx context.Context, taskID uint64) ([]*model.TaskInstance, error) {
	var instances []*model.TaskInstance
	db := r.db.WithContext(ctx).Where("status IN ?", []int8{model.InstanceStatusPending, model.InstanceStatusScheduling, model.InstanceStatusRunning})
	if taskID > 0 {
		db = db.Where("task_id = ?", taskID)
	}
	err := db.Find(&instances).Error
	return instances, err
}

// GetInstancesByTriggerTime 根据触发时间获取实例
func (r *instanceRepository) GetInstancesByTriggerTime(ctx context.Context, taskID uint64, triggerTime time.Time) ([]*model.TaskInstance, error) {
	var instances []*model.TaskInstance
	err := r.db.WithContext(ctx).
		Where("task_id = ? AND trigger_time = ?", taskID, triggerTime).
		Find(&instances).Error
	return instances, err
}

// CountByStatus 按状态统计实例数量
func (r *instanceRepository) CountByStatus(ctx context.Context, taskID uint64, startTime, endTime time.Time) (map[int8]int64, error) {
	type Result struct {
		Status int8
		Count  int64
	}
	var results []Result

	db := r.db.WithContext(ctx).Model(&model.TaskInstance{}).
		Select("status, COUNT(*) as count").
		Where("trigger_time BETWEEN ? AND ?", startTime, endTime)

	if taskID > 0 {
		db = db.Where("task_id = ?", taskID)
	}

	if err := db.Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int8]int64)
	for _, r := range results {
		countMap[r.Status] = r.Count
	}
	return countMap, nil
}

// GetRecentInstances 获取最近的实例
func (r *instanceRepository) GetRecentInstances(ctx context.Context, limit int) ([]*model.TaskInstance, error) {
	var instances []*model.TaskInstance
	err := r.db.WithContext(ctx).
		Preload("Task").
		Order("id DESC").
		Limit(limit).
		Find(&instances).Error
	return instances, err
}

// TaskLogRepository 任务日志仓库接口
type TaskLogRepository interface {
	Create(ctx context.Context, log *model.TaskLog) error
	BatchCreate(ctx context.Context, logs []*model.TaskLog) error
	GetByInstanceID(ctx context.Context, instanceID uint64, page, pageSize int) ([]*model.TaskLog, int64, error)
}

// taskLogRepository 任务日志仓库实现
type taskLogRepository struct {
	db *gorm.DB
}

// NewTaskLogRepository 创建任务日志仓库
func NewTaskLogRepository() TaskLogRepository {
	return &taskLogRepository{db: mysql.GetDB()}
}

// Create 创建任务日志
func (r *taskLogRepository) Create(ctx context.Context, log *model.TaskLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// BatchCreate 批量创建任务日志
func (r *taskLogRepository) BatchCreate(ctx context.Context, logs []*model.TaskLog) error {
	return r.db.WithContext(ctx).CreateInBatches(logs, 100).Error
}

// GetByInstanceID 根据实例ID获取日志
func (r *taskLogRepository) GetByInstanceID(ctx context.Context, instanceID uint64, page, pageSize int) ([]*model.TaskLog, int64, error) {
	var logs []*model.TaskLog
	var total int64

	db := r.db.WithContext(ctx).Model(&model.TaskLog{}).Where("instance_id = ?", instanceID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(mysql.Paginate(page, pageSize)).Order("log_time ASC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

