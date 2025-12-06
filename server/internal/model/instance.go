package model

import (
	"time"
)

// TaskInstance 任务实例(执行记录)
type TaskInstance struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID          uint64     `gorm:"not null;index" json:"task_id"`
	GroupID         uint64     `gorm:"not null;index" json:"group_id"`
	ExecutorID      string     `gorm:"size:128;index" json:"executor_id"`
	ExecutorAddress string     `gorm:"size:256" json:"executor_address"`
	ExecutorHandler string     `gorm:"size:256" json:"executor_handler"`
	ExecutorParam   string     `gorm:"type:text" json:"executor_param"`
	ShardIndex      uint       `gorm:"default:0" json:"shard_index"`
	ShardTotal      uint       `gorm:"default:1" json:"shard_total"`
	TriggerType     string     `gorm:"size:32;default:CRON" json:"trigger_type"`
	TriggerTime     time.Time  `gorm:"not null;index" json:"trigger_time"`
	ScheduleTime    *time.Time `json:"schedule_time"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	Status          int8       `gorm:"default:0;index" json:"status"`
	ResultCode      int        `gorm:"default:0" json:"result_code"`
	ResultMsg       string     `gorm:"type:text" json:"result_msg"`
	RetryCount      uint       `gorm:"default:0" json:"retry_count"`
	AlarmStatus     int8       `gorm:"default:0" json:"alarm_status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Task            *Task      `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

// TableName 指定表名
func (TaskInstance) TableName() string {
	return "task_instance"
}

// TaskLog 任务执行日志
type TaskLog struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	InstanceID uint64    `gorm:"not null;index" json:"instance_id"`
	TaskID     uint64    `gorm:"not null;index" json:"task_id"`
	LogTime    time.Time `gorm:"not null;index;type:datetime(3)" json:"log_time"`
	LogLevel   string    `gorm:"size:16;default:INFO" json:"log_level"`
	LogContent string    `gorm:"type:text" json:"log_content"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 指定表名
func (TaskLog) TableName() string {
	return "task_log"
}

// 任务实例状态常量
const (
	InstanceStatusPending    = 0 // 待调度
	InstanceStatusScheduling = 1 // 调度中
	InstanceStatusRunning    = 2 // 执行中
	InstanceStatusSuccess    = 3 // 执行成功
	InstanceStatusFailed     = 4 // 执行失败
	InstanceStatusCancelled  = 5 // 已取消
)

// 触发类型常量
const (
	TriggerTypeCron   = "CRON"   // Cron触发
	TriggerTypeManual = "MANUAL" // 手动触发
	TriggerTypeParent = "PARENT" // 父任务触发
	TriggerTypeAPI    = "API"    // API触发
	TriggerTypeRetry  = "RETRY"  // 重试触发
)

// 日志级别常量
const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

// 告警状态常量
const (
	AlarmStatusDefault = 0 // 默认
	AlarmStatusAlarmed = 1 // 已告警
)

// GetStatusText 获取状态文本
func (i *TaskInstance) GetStatusText() string {
	switch i.Status {
	case InstanceStatusPending:
		return "待调度"
	case InstanceStatusScheduling:
		return "调度中"
	case InstanceStatusRunning:
		return "执行中"
	case InstanceStatusSuccess:
		return "执行成功"
	case InstanceStatusFailed:
		return "执行失败"
	case InstanceStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// Duration 计算执行时长(毫秒)
func (i *TaskInstance) Duration() int64 {
	if i.StartTime == nil || i.EndTime == nil {
		return 0
	}
	return i.EndTime.Sub(*i.StartTime).Milliseconds()
}

