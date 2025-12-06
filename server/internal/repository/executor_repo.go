package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"distributed-scheduler/internal/model"
	"distributed-scheduler/pkg/mysql"
)

// ExecutorRepository 执行器仓库接口
type ExecutorRepository interface {
	Register(ctx context.Context, node *model.ExecutorNode) error
	Unregister(ctx context.Context, id string) error
	UpdateHeartbeat(ctx context.Context, heartbeat *model.ExecutorHeartbeat) error
	GetByID(ctx context.Context, id string) (*model.ExecutorNode, error)
	GetOnlineByGroupID(ctx context.Context, groupID uint64) ([]*model.ExecutorNode, error)
	GetOnlineByAppName(ctx context.Context, appName string) ([]*model.ExecutorNode, error)
	GetAllOnline(ctx context.Context) ([]*model.ExecutorNode, error)
	SetOffline(ctx context.Context, id string) error
	SetOfflineByTimeout(ctx context.Context, timeout time.Duration) (int64, error)
	UpdateLoad(ctx context.Context, id string, load uint) error
	List(ctx context.Context, page, pageSize int, groupID uint64, status int8) ([]*model.ExecutorNode, int64, error)
}

// executorRepository 执行器仓库实现
type executorRepository struct {
	db *gorm.DB
}

// NewExecutorRepository 创建执行器仓库
func NewExecutorRepository() ExecutorRepository {
	return &executorRepository{db: mysql.GetDB()}
}

// Register 注册执行器
func (r *executorRepository) Register(ctx context.Context, node *model.ExecutorNode) error {
	// 使用ON DUPLICATE KEY UPDATE实现注册/更新
	return r.db.WithContext(ctx).Save(node).Error
}

// Unregister 注销执行器
func (r *executorRepository) Unregister(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.ExecutorNode{}, "id = ?", id).Error
}

// UpdateHeartbeat 更新心跳
func (r *executorRepository) UpdateHeartbeat(ctx context.Context, heartbeat *model.ExecutorHeartbeat) error {
	return r.db.WithContext(ctx).Model(&model.ExecutorNode{}).Where("id = ?", heartbeat.ExecutorID).
		Updates(map[string]interface{}{
			"current_load":   heartbeat.CurrentLoad,
			"cpu_usage":      heartbeat.CPUUsage,
			"memory_usage":   heartbeat.MemoryUsage,
			"status":         model.ExecutorStatusOnline,
			"last_heartbeat": gorm.Expr("NOW()"),
		}).Error
}

// GetByID 根据ID获取执行器
func (r *executorRepository) GetByID(ctx context.Context, id string) (*model.ExecutorNode, error) {
	var node model.ExecutorNode
	err := r.db.WithContext(ctx).First(&node, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// GetOnlineByGroupID 获取任务组的在线执行器
func (r *executorRepository) GetOnlineByGroupID(ctx context.Context, groupID uint64) ([]*model.ExecutorNode, error) {
	var nodes []*model.ExecutorNode
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND status = ?", groupID, model.ExecutorStatusOnline).
		Find(&nodes).Error
	return nodes, err
}

// GetOnlineByAppName 根据应用名获取在线执行器
func (r *executorRepository) GetOnlineByAppName(ctx context.Context, appName string) ([]*model.ExecutorNode, error) {
	var nodes []*model.ExecutorNode
	err := r.db.WithContext(ctx).
		Where("app_name = ? AND status = ?", appName, model.ExecutorStatusOnline).
		Find(&nodes).Error
	return nodes, err
}

// GetAllOnline 获取所有在线执行器
func (r *executorRepository) GetAllOnline(ctx context.Context) ([]*model.ExecutorNode, error) {
	var nodes []*model.ExecutorNode
	err := r.db.WithContext(ctx).
		Where("status = ?", model.ExecutorStatusOnline).
		Find(&nodes).Error
	return nodes, err
}

// SetOffline 设置执行器离线
func (r *executorRepository) SetOffline(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.ExecutorNode{}).Where("id = ?", id).
		Update("status", model.ExecutorStatusOffline).Error
}

// SetOfflineByTimeout 设置超时的执行器为离线
func (r *executorRepository) SetOfflineByTimeout(ctx context.Context, timeout time.Duration) (int64, error) {
	deadLine := time.Now().Add(-timeout)
	result := r.db.WithContext(ctx).Model(&model.ExecutorNode{}).
		Where("status = ? AND last_heartbeat < ?", model.ExecutorStatusOnline, deadLine).
		Update("status", model.ExecutorStatusOffline)
	return result.RowsAffected, result.Error
}

// UpdateLoad 更新执行器负载
func (r *executorRepository) UpdateLoad(ctx context.Context, id string, load uint) error {
	return r.db.WithContext(ctx).Model(&model.ExecutorNode{}).Where("id = ?", id).
		Update("current_load", load).Error
}

// List 获取执行器列表
func (r *executorRepository) List(ctx context.Context, page, pageSize int, groupID uint64, status int8) ([]*model.ExecutorNode, int64, error) {
	var nodes []*model.ExecutorNode
	var total int64

	db := r.db.WithContext(ctx).Model(&model.ExecutorNode{})

	if groupID > 0 {
		db = db.Where("group_id = ?", groupID)
	}
	if status >= 0 {
		db = db.Where("status = ?", status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(mysql.Paginate(page, pageSize)).Order("registered_at DESC").Find(&nodes).Error; err != nil {
		return nil, 0, err
	}

	return nodes, total, nil
}

