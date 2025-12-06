package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"distributed-scheduler/internal/common/utils"
	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/repository"
)

var (
	ErrTaskNotFound  = errors.New("任务不存在")
	ErrGroupNotFound = errors.New("任务组不存在")
	ErrInvalidCron   = errors.New("无效的Cron表达式")
)

// TaskService 任务服务接口
type TaskService interface {
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.Task, error)
	List(ctx context.Context, page, pageSize int, groupID uint64, keyword string, status int8) ([]*model.Task, int64, error)
	Start(ctx context.Context, id uint64) error
	Stop(ctx context.Context, id uint64) error
	Trigger(ctx context.Context, id uint64, param string) (*model.TaskInstance, error)
	GetNextTriggerTimes(ctx context.Context, cron string, count int) ([]time.Time, error)
}

// taskService 任务服务实现
type taskService struct {
	taskRepo     repository.TaskRepository
	groupRepo    repository.TaskGroupRepository
	instanceRepo repository.InstanceRepository
}

// NewTaskService 创建任务服务
func NewTaskService() TaskService {
	return &taskService{
		taskRepo:     repository.NewTaskRepository(),
		groupRepo:    repository.NewTaskGroupRepository(),
		instanceRepo: repository.NewInstanceRepository(),
	}
}

// Create 创建任务
func (s *taskService) Create(ctx context.Context, task *model.Task) error {
	// 验证Cron表达式
	if err := utils.ValidateCron(task.Cron); err != nil {
		return ErrInvalidCron
	}

	// 验证任务组是否存在
	_, err := s.groupRepo.GetByID(ctx, task.GroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrGroupNotFound
		}
		return err
	}

	// 计算下次触发时间
	nextTime, err := utils.GetNextTriggerTime(task.Cron, time.Now())
	if err != nil {
		return err
	}
	task.NextTriggerTime = &nextTime

	return s.taskRepo.Create(ctx, task)
}

// Update 更新任务
func (s *taskService) Update(ctx context.Context, task *model.Task) error {
	// 验证Cron表达式
	if err := utils.ValidateCron(task.Cron); err != nil {
		return ErrInvalidCron
	}

	// 重新计算下次触发时间
	nextTime, err := utils.GetNextTriggerTime(task.Cron, time.Now())
	if err != nil {
		return err
	}
	task.NextTriggerTime = &nextTime

	// 乐观锁更新
	task.Version++

	return s.taskRepo.Update(ctx, task)
}

// Delete 删除任务
func (s *taskService) Delete(ctx context.Context, id uint64) error {
	return s.taskRepo.Delete(ctx, id)
}

// GetByID 根据ID获取任务
func (s *taskService) GetByID(ctx context.Context, id uint64) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	// 获取依赖任务
	deps, err := s.taskRepo.GetDependencies(ctx, id)
	if err != nil {
		return nil, err
	}
	task.Dependencies = deps

	return task, nil
}

// List 获取任务列表
func (s *taskService) List(ctx context.Context, page, pageSize int, groupID uint64, keyword string, status int8) ([]*model.Task, int64, error) {
	return s.taskRepo.List(ctx, page, pageSize, groupID, keyword, status)
}

// Start 启动任务
func (s *taskService) Start(ctx context.Context, id uint64) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 计算下次触发时间
	nextTime, err := utils.GetNextTriggerTime(task.Cron, time.Now())
	if err != nil {
		return err
	}

	task.Status = model.TaskStatusEnabled
	task.NextTriggerTime = &nextTime

	return s.taskRepo.Update(ctx, task)
}

// Stop 停止任务
func (s *taskService) Stop(ctx context.Context, id uint64) error {
	return s.taskRepo.UpdateStatus(ctx, id, model.TaskStatusDisabled)
}

// Trigger 手动触发任务
func (s *taskService) Trigger(ctx context.Context, id uint64, param string) (*model.TaskInstance, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	// 创建任务实例
	instance := &model.TaskInstance{
		TaskID:          task.ID,
		GroupID:         task.GroupID,
		ExecutorHandler: task.ExecutorHandler,
		ExecutorParam:   param,
		ShardIndex:      0,
		ShardTotal:      1,
		TriggerType:     model.TriggerTypeManual,
		TriggerTime:     time.Now(),
		Status:          model.InstanceStatusPending,
	}

	if param == "" {
		instance.ExecutorParam = task.ExecutorParam
	}

	if err := s.instanceRepo.Create(ctx, instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// GetNextTriggerTimes 获取下N次触发时间
func (s *taskService) GetNextTriggerTimes(ctx context.Context, cron string, count int) ([]time.Time, error) {
	return utils.GetNextNTriggerTimes(cron, time.Now(), count)
}

// TaskGroupService 任务组服务接口
type TaskGroupService interface {
	Create(ctx context.Context, group *model.TaskGroup) error
	Update(ctx context.Context, group *model.TaskGroup) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.TaskGroup, error)
	List(ctx context.Context, page, pageSize int, keyword string) ([]*model.TaskGroup, int64, error)
	GetAll(ctx context.Context) ([]*model.TaskGroup, error)
}

// taskGroupService 任务组服务实现
type taskGroupService struct {
	groupRepo repository.TaskGroupRepository
}

// NewTaskGroupService 创建任务组服务
func NewTaskGroupService() TaskGroupService {
	return &taskGroupService{
		groupRepo: repository.NewTaskGroupRepository(),
	}
}

// Create 创建任务组
func (s *taskGroupService) Create(ctx context.Context, group *model.TaskGroup) error {
	return s.groupRepo.Create(ctx, group)
}

// Update 更新任务组
func (s *taskGroupService) Update(ctx context.Context, group *model.TaskGroup) error {
	return s.groupRepo.Update(ctx, group)
}

// Delete 删除任务组
func (s *taskGroupService) Delete(ctx context.Context, id uint64) error {
	return s.groupRepo.Delete(ctx, id)
}

// GetByID 根据ID获取任务组
func (s *taskGroupService) GetByID(ctx context.Context, id uint64) (*model.TaskGroup, error) {
	return s.groupRepo.GetByID(ctx, id)
}

// List 获取任务组列表
func (s *taskGroupService) List(ctx context.Context, page, pageSize int, keyword string) ([]*model.TaskGroup, int64, error) {
	return s.groupRepo.List(ctx, page, pageSize, keyword)
}

// GetAll 获取所有任务组
func (s *taskGroupService) GetAll(ctx context.Context) ([]*model.TaskGroup, error) {
	return s.groupRepo.GetAll(ctx)
}

