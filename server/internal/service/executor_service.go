package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/repository"
)

var (
	ErrExecutorNotFound = errors.New("执行器不存在")
	ErrNoAvailableNode  = errors.New("没有可用的执行节点")
)

// ExecutorService 执行器服务接口
type ExecutorService interface {
	Register(ctx context.Context, appName, host string, port uint, maxConcurrent uint) (*model.ExecutorNode, error)
	Unregister(ctx context.Context, id string) error
	Heartbeat(ctx context.Context, heartbeat *model.ExecutorHeartbeat) error
	GetByID(ctx context.Context, id string) (*model.ExecutorNode, error)
	GetOnlineByGroupID(ctx context.Context, groupID uint64) ([]*model.ExecutorNode, error)
	List(ctx context.Context, page, pageSize int, groupID uint64, status int8) ([]*model.ExecutorNode, int64, error)
	CheckOfflineExecutors(ctx context.Context, timeout time.Duration) (int64, error)
}

// executorService 执行器服务实现
type executorService struct {
	executorRepo repository.ExecutorRepository
	groupRepo    repository.TaskGroupRepository
}

// NewExecutorService 创建执行器服务
func NewExecutorService() ExecutorService {
	return &executorService{
		executorRepo: repository.NewExecutorRepository(),
		groupRepo:    repository.NewTaskGroupRepository(),
	}
}

// Register 注册执行器
func (s *executorService) Register(ctx context.Context, appName, host string, port uint, maxConcurrent uint) (*model.ExecutorNode, error) {
	// 获取任务组
	group, err := s.groupRepo.GetByAppName(ctx, appName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("任务组不存在: %s", appName)
		}
		return nil, err
	}

	// 生成执行器ID
	nodeID := uuid.New().String()

	node := &model.ExecutorNode{
		ID:            nodeID,
		GroupID:       group.ID,
		AppName:       appName,
		Host:          host,
		Port:          port,
		Weight:        100,
		MaxConcurrent: maxConcurrent,
		CurrentLoad:   0,
		Status:        model.ExecutorStatusOnline,
		LastHeartbeat: time.Now(),
		RegisteredAt:  time.Now(),
	}

	if err := s.executorRepo.Register(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

// Unregister 注销执行器
func (s *executorService) Unregister(ctx context.Context, id string) error {
	return s.executorRepo.Unregister(ctx, id)
}

// Heartbeat 心跳
func (s *executorService) Heartbeat(ctx context.Context, heartbeat *model.ExecutorHeartbeat) error {
	return s.executorRepo.UpdateHeartbeat(ctx, heartbeat)
}

// GetByID 根据ID获取执行器
func (s *executorService) GetByID(ctx context.Context, id string) (*model.ExecutorNode, error) {
	node, err := s.executorRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExecutorNotFound
		}
		return nil, err
	}
	return node, nil
}

// GetOnlineByGroupID 获取任务组的在线执行器
func (s *executorService) GetOnlineByGroupID(ctx context.Context, groupID uint64) ([]*model.ExecutorNode, error) {
	return s.executorRepo.GetOnlineByGroupID(ctx, groupID)
}

// List 获取执行器列表
func (s *executorService) List(ctx context.Context, page, pageSize int, groupID uint64, status int8) ([]*model.ExecutorNode, int64, error) {
	return s.executorRepo.List(ctx, page, pageSize, groupID, status)
}

// CheckOfflineExecutors 检查离线执行器
func (s *executorService) CheckOfflineExecutors(ctx context.Context, timeout time.Duration) (int64, error) {
	return s.executorRepo.SetOfflineByTimeout(ctx, timeout)
}

