package model

import (
	"fmt"
	"time"
)

// ExecutorNode 执行器节点
type ExecutorNode struct {
	ID            string    `gorm:"primaryKey;size:128" json:"id"`
	GroupID       uint64    `gorm:"not null;index" json:"group_id"`
	AppName       string    `gorm:"size:64;not null;index" json:"app_name"`
	Host          string    `gorm:"size:128;not null" json:"host"`
	Port          uint      `gorm:"not null" json:"port"`
	Weight        uint      `gorm:"default:100" json:"weight"`
	MaxConcurrent uint      `gorm:"default:100" json:"max_concurrent"`
	CurrentLoad   uint      `gorm:"default:0" json:"current_load"`
	CPUUsage      float64   `gorm:"type:decimal(5,2);default:0" json:"cpu_usage"`
	MemoryUsage   float64   `gorm:"type:decimal(5,2);default:0" json:"memory_usage"`
	Status        int8      `gorm:"default:1;index" json:"status"`
	LastHeartbeat time.Time `gorm:"index" json:"last_heartbeat"`
	RegisteredAt  time.Time `json:"registered_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ExecutorNode) TableName() string {
	return "executor_node"
}

// Address 获取执行器地址
func (e *ExecutorNode) Address() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

// IsOnline 是否在线
func (e *ExecutorNode) IsOnline() bool {
	return e.Status == ExecutorStatusOnline
}

// IsOverload 是否过载
func (e *ExecutorNode) IsOverload() bool {
	return e.CurrentLoad >= e.MaxConcurrent
}

// 执行器状态常量
const (
	ExecutorStatusOffline = 0 // 离线
	ExecutorStatusOnline  = 1 // 在线
)

// ExecutorHeartbeat 执行器心跳数据
type ExecutorHeartbeat struct {
	ExecutorID  string  `json:"executor_id"`
	AppName     string  `json:"app_name"`
	Host        string  `json:"host"`
	Port        uint    `json:"port"`
	CurrentLoad uint    `json:"current_load"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

// ExecutorTask 执行器任务参数
type ExecutorTask struct {
	InstanceID      uint64 `json:"instance_id"`
	TaskID          uint64 `json:"task_id"`
	ExecutorHandler string `json:"executor_handler"`
	ExecutorParam   string `json:"executor_param"`
	ShardIndex      uint   `json:"shard_index"`
	ShardTotal      uint   `json:"shard_total"`
	Timeout         uint   `json:"timeout"`
}

// ExecutorResult 执行器返回结果
type ExecutorResult struct {
	InstanceID uint64 `json:"instance_id"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

