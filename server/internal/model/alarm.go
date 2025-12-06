package model

import (
	"time"

	"gorm.io/gorm"
)

// AlarmRule 告警规则
type AlarmRule struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	RuleType     string         `gorm:"size:32;not null;index" json:"rule_type"`
	GroupID      uint64         `gorm:"default:0;index" json:"group_id"`
	TaskID       uint64         `gorm:"default:0;index" json:"task_id"`
	Threshold    uint           `gorm:"default:1" json:"threshold"`
	AlarmLevel   string         `gorm:"size:16;default:WARNING" json:"alarm_level"`
	NotifyType   string         `gorm:"size:64;default:EMAIL" json:"notify_type"`
	NotifyTarget string         `gorm:"type:text" json:"notify_target"`
	Status       int8           `gorm:"default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AlarmRule) TableName() string {
	return "alarm_rule"
}

// AlarmRecord 告警记录
type AlarmRecord struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RuleID       uint64     `gorm:"not null;index" json:"rule_id"`
	TaskID       uint64     `gorm:"default:0;index" json:"task_id"`
	InstanceID   uint64     `gorm:"default:0" json:"instance_id"`
	AlarmType    string     `gorm:"size:32;not null;index" json:"alarm_type"`
	AlarmLevel   string     `gorm:"size:16;not null" json:"alarm_level"`
	AlarmTitle   string     `gorm:"size:256;not null" json:"alarm_title"`
	AlarmContent string     `gorm:"type:text" json:"alarm_content"`
	NotifyStatus int8       `gorm:"default:0" json:"notify_status"`
	NotifyTime   *time.Time `json:"notify_time"`
	NotifyResult string     `gorm:"type:text" json:"notify_result"`
	CreatedAt    time.Time  `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (AlarmRecord) TableName() string {
	return "alarm_record"
}

// 告警规则类型常量
const (
	AlarmRuleTypeTaskFail        = "TASK_FAIL"        // 任务失败
	AlarmRuleTypeTaskTimeout     = "TASK_TIMEOUT"     // 任务超时
	AlarmRuleTypeExecutorOffline = "EXECUTOR_OFFLINE" // 执行器离线
)

// 告警级别常量
const (
	AlarmLevelInfo     = "INFO"
	AlarmLevelWarning  = "WARNING"
	AlarmLevelError    = "ERROR"
	AlarmLevelCritical = "CRITICAL"
)

// 通知类型常量
const (
	NotifyTypeEmail   = "EMAIL"
	NotifyTypeSMS     = "SMS"
	NotifyTypeWebhook = "WEBHOOK"
)

// 通知状态常量
const (
	NotifyStatusPending = 0 // 待发送
	NotifyStatusSent    = 1 // 已发送
	NotifyStatusFailed  = 2 // 发送失败
)

