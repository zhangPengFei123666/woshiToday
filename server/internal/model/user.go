package model

import (
	"time"

	"gorm.io/gorm"
)

// SysUser 系统用户
type SysUser struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Username      string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Password      string         `gorm:"size:128;not null" json:"-"`
	Nickname      string         `gorm:"size:64" json:"nickname"`
	Email         string         `gorm:"size:128" json:"email"`
	Phone         string         `gorm:"size:20" json:"phone"`
	Avatar        string         `gorm:"size:256" json:"avatar"`
	Status        int8           `gorm:"default:1" json:"status"`
	LastLoginTime *time.Time     `json:"last_login_time"`
	LastLoginIP   string         `gorm:"size:64" json:"last_login_ip"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Roles         []SysRole      `gorm:"many2many:sys_user_role" json:"roles,omitempty"`
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}

// SysRole 系统角色
type SysRole struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"size:64;not null" json:"name"`
	Code        string         `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Description string         `gorm:"size:256" json:"description"`
	Status      int8           `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SysRole) TableName() string {
	return "sys_role"
}

// SysUserRole 用户角色关联
type SysUserRole struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;uniqueIndex:uk_user_role" json:"user_id"`
	RoleID    uint64    `gorm:"not null;uniqueIndex:uk_user_role;index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (SysUserRole) TableName() string {
	return "sys_user_role"
}

// UserStatus 用户状态
const (
	UserStatusDisabled = 0 // 禁用
	UserStatusEnabled  = 1 // 启用
)

