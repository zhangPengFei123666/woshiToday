package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/repository"
)

var (
	ErrInstanceNotFound = errors.New("任务实例不存在")
)

// InstanceService 任务实例服务接口
type InstanceService interface {
	GetByID(ctx context.Context, id uint64) (*model.TaskInstance, error)
	List(ctx context.Context, page, pageSize int, taskID uint64, status int8, startTime, endTime *time.Time) ([]*model.TaskInstance, int64, error)
	Cancel(ctx context.Context, id uint64) error
	Retry(ctx context.Context, id uint64) (*model.TaskInstance, error)
	GetLogs(ctx context.Context, instanceID uint64, page, pageSize int) ([]*model.TaskLog, int64, error)
	GetStatistics(ctx context.Context, taskID uint64, startTime, endTime time.Time) (*InstanceStatistics, error)
	GetRecentInstances(ctx context.Context, limit int) ([]*model.TaskInstance, error)
}

// InstanceStatistics 实例统计
type InstanceStatistics struct {
	Total     int64   `json:"total"`
	Success   int64   `json:"success"`
	Failed    int64   `json:"failed"`
	Running   int64   `json:"running"`
	Pending   int64   `json:"pending"`
	Cancelled int64   `json:"cancelled"`
	Rate      float64 `json:"rate"` // 成功率
}

// instanceService 任务实例服务实现
type instanceService struct {
	instanceRepo repository.InstanceRepository
	logRepo      repository.TaskLogRepository
	taskRepo     repository.TaskRepository
}

// NewInstanceService 创建任务实例服务
func NewInstanceService() InstanceService {
	return &instanceService{
		instanceRepo: repository.NewInstanceRepository(),
		logRepo:      repository.NewTaskLogRepository(),
		taskRepo:     repository.NewTaskRepository(),
	}
}

// GetByID 根据ID获取任务实例
func (s *instanceService) GetByID(ctx context.Context, id uint64) (*model.TaskInstance, error) {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInstanceNotFound
		}
		return nil, err
	}
	return instance, nil
}

// List 获取任务实例列表
func (s *instanceService) List(ctx context.Context, page, pageSize int, taskID uint64, status int8, startTime, endTime *time.Time) ([]*model.TaskInstance, int64, error) {
	return s.instanceRepo.List(ctx, page, pageSize, taskID, status, startTime, endTime)
}

// Cancel 取消任务实例
func (s *instanceService) Cancel(ctx context.Context, id uint64) error {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 只能取消待调度或调度中的任务
	if instance.Status != model.InstanceStatusPending && instance.Status != model.InstanceStatusScheduling {
		return errors.New("只能取消待调度或调度中的任务")
	}

	return s.instanceRepo.UpdateStatus(ctx, id, model.InstanceStatusCancelled, 0, "用户取消")
}

// Retry 重试任务实例
func (s *instanceService) Retry(ctx context.Context, id uint64) (*model.TaskInstance, error) {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 只能重试失败的任务
	if instance.Status != model.InstanceStatusFailed {
		return nil, errors.New("只能重试失败的任务")
	}

	// 创建新的任务实例
	newInstance := &model.TaskInstance{
		TaskID:          instance.TaskID,
		GroupID:         instance.GroupID,
		ExecutorHandler: instance.ExecutorHandler,
		ExecutorParam:   instance.ExecutorParam,
		ShardIndex:      instance.ShardIndex,
		ShardTotal:      instance.ShardTotal,
		TriggerType:     model.TriggerTypeRetry,
		TriggerTime:     time.Now(),
		Status:          model.InstanceStatusPending,
	}

	if err := s.instanceRepo.Create(ctx, newInstance); err != nil {
		return nil, err
	}

	return newInstance, nil
}

// GetLogs 获取任务实例日志
func (s *instanceService) GetLogs(ctx context.Context, instanceID uint64, page, pageSize int) ([]*model.TaskLog, int64, error) {
	return s.logRepo.GetByInstanceID(ctx, instanceID, page, pageSize)
}

// GetStatistics 获取实例统计
func (s *instanceService) GetStatistics(ctx context.Context, taskID uint64, startTime, endTime time.Time) (*InstanceStatistics, error) {
	countMap, err := s.instanceRepo.CountByStatus(ctx, taskID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	stats := &InstanceStatistics{
		Success:   countMap[model.InstanceStatusSuccess],
		Failed:    countMap[model.InstanceStatusFailed],
		Running:   countMap[model.InstanceStatusRunning],
		Pending:   countMap[model.InstanceStatusPending] + countMap[model.InstanceStatusScheduling],
		Cancelled: countMap[model.InstanceStatusCancelled],
	}
	stats.Total = stats.Success + stats.Failed + stats.Running + stats.Pending + stats.Cancelled

	if stats.Success+stats.Failed > 0 {
		stats.Rate = float64(stats.Success) / float64(stats.Success+stats.Failed) * 100
	}

	return stats, nil
}

// GetRecentInstances 获取最近的实例
func (s *instanceService) GetRecentInstances(ctx context.Context, limit int) ([]*model.TaskInstance, error) {
	return s.instanceRepo.GetRecentInstances(ctx, limit)
}

