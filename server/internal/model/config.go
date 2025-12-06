package model

import (
	"time"
)

// SysConfig 系统配置
type SysConfig struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"size:128;not null;uniqueIndex" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	ConfigType  string    `gorm:"size:32;default:STRING" json:"config_type"`
	Description string    `gorm:"size:256" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (SysConfig) TableName() string {
	return "sys_config"
}

// OperationLog 操作日志
type OperationLog struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"default:0;index" json:"user_id"`
	Username      string    `gorm:"size:64" json:"username"`
	Module        string    `gorm:"size:64;index" json:"module"`
	Action        string    `gorm:"size:64" json:"action"`
	TargetType    string    `gorm:"size:64" json:"target_type"`
	TargetID      uint64    `gorm:"default:0" json:"target_id"`
	RequestMethod string    `gorm:"size:16" json:"request_method"`
	RequestURL    string    `gorm:"size:512" json:"request_url"`
	RequestParam  string    `gorm:"type:text" json:"request_param"`
	RequestIP     string    `gorm:"size:64" json:"request_ip"`
	UserAgent     string    `gorm:"size:512" json:"user_agent"`
	ResultCode    int       `gorm:"default:0" json:"result_code"`
	ResultMsg     string    `gorm:"size:512" json:"result_msg"`
	Duration      uint      `gorm:"default:0" json:"duration"`
	CreatedAt     time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_log"
}

// 配置类型常量
const (
	ConfigTypeString  = "STRING"
	ConfigTypeNumber  = "NUMBER"
	ConfigTypeBoolean = "BOOLEAN"
	ConfigTypeJSON    = "JSON"
)

