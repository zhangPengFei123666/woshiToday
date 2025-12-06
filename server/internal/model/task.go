package model

import (
	"time"

	"gorm.io/gorm"
)

// TaskGroup 任务组
type TaskGroup struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"size:128;not null" json:"name"`
	Description string         `gorm:"size:512" json:"description"`
	AppName     string         `gorm:"size:64;not null;uniqueIndex" json:"app_name"`
	Status      int8           `gorm:"default:1" json:"status"`
	CreatedBy   uint64         `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (TaskGroup) TableName() string {
	return "task_group"
}

// Task 任务定义
type Task struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID         uint64         `gorm:"not null;index" json:"group_id"`
	Name            string         `gorm:"size:128;not null" json:"name"`
	Description     string         `gorm:"size:512" json:"description"`
	Cron            string         `gorm:"size:64;not null" json:"cron"`
	ExecutorType    string         `gorm:"size:32;not null;default:HTTP" json:"executor_type"`
	ExecutorHandler string         `gorm:"size:256;not null" json:"executor_handler"`
	ExecutorParam   string         `gorm:"type:text" json:"executor_param"`
	RouteStrategy   string         `gorm:"size:32;default:ROUND_ROBIN" json:"route_strategy"`
	BlockStrategy   string         `gorm:"size:32;default:SERIAL_EXECUTION" json:"block_strategy"`
	ShardNum        uint           `gorm:"default:1" json:"shard_num"`
	RetryCount      uint           `gorm:"default:0" json:"retry_count"`
	RetryInterval   uint           `gorm:"default:0" json:"retry_interval"`
	Timeout         uint           `gorm:"default:0" json:"timeout"`
	AlarmEmail      string         `gorm:"size:512" json:"alarm_email"`
	Priority        int            `gorm:"default:0" json:"priority"`
	Status          int8           `gorm:"default:1;index" json:"status"`
	Version         uint           `gorm:"default:0" json:"version"`
	NextTriggerTime *time.Time     `gorm:"index" json:"next_trigger_time"`
	LastTriggerTime *time.Time     `json:"last_trigger_time"`
	CreatedBy       uint64         `json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Group           *TaskGroup     `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Dependencies    []Task         `gorm:"many2many:task_dependency;joinForeignKey:TaskID;joinReferences:DependTaskID" json:"dependencies,omitempty"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "task"
}

// TaskDependency 任务依赖关系
type TaskDependency struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID       uint64    `gorm:"not null;uniqueIndex:uk_task_depend" json:"task_id"`
	DependTaskID uint64    `gorm:"not null;uniqueIndex:uk_task_depend;index" json:"depend_task_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (TaskDependency) TableName() string {
	return "task_dependency"
}

// 任务状态常量
const (
	TaskStatusDisabled = 0 // 禁用
	TaskStatusEnabled  = 1 // 启用
)

// 执行器类型常量
const (
	ExecutorTypeHTTP   = "HTTP"
	ExecutorTypeGRPC   = "GRPC"
	ExecutorTypeSCRIPT = "SCRIPT"
)

// 路由策略常量
const (
	RouteStrategyRoundRobin          = "ROUND_ROBIN"           // 轮询
	RouteStrategyRandom              = "RANDOM"                // 随机
	RouteStrategyConsistentHash      = "CONSISTENT_HASH"       // 一致性哈希
	RouteStrategyLeastFrequentlyUsed = "LEAST_FREQUENTLY_USED" // 最少使用
	RouteStrategyLeastRecentlyUsed   = "LEAST_RECENTLY_USED"   // 最近最少使用
	RouteStrategyFailover            = "FAILOVER"              // 故障转移
	RouteStrategyShardingBroadcast   = "SHARDING_BROADCAST"    // 分片广播
)

// 阻塞策略常量
const (
	BlockStrategySerialExecution = "SERIAL_EXECUTION" // 串行执行
	BlockStrategyDiscardLater    = "DISCARD_LATER"    // 丢弃后续调度
	BlockStrategyCoverEarly      = "COVER_EARLY"      // 覆盖之前调度
)

