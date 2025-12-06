package repository

import (
	"context"

	"gorm.io/gorm"

	"distributed-scheduler/internal/model"
	"distributed-scheduler/pkg/mysql"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(ctx context.Context, user *model.SysUser) error
	Update(ctx context.Context, user *model.SysUser) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.SysUser, error)
	GetByUsername(ctx context.Context, username string) (*model.SysUser, error)
	List(ctx context.Context, page, pageSize int, keyword string) ([]*model.SysUser, int64, error)
	UpdateLoginInfo(ctx context.Context, id uint64, ip string) error
	GetUserRoles(ctx context.Context, userID uint64) ([]model.SysRole, error)
	AssignRoles(ctx context.Context, userID uint64, roleIDs []uint64) error
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository() UserRepository {
	return &userRepository{db: mysql.GetDB()}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.SysUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.SysUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户(软删除)
func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysUser{}, id).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]*model.SysUser, int64, error) {
	var users []*model.SysUser
	var total int64

	db := r.db.WithContext(ctx).Model(&model.SysUser{})

	if keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(mysql.Paginate(page, pageSize)).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLoginInfo 更新登录信息
func (r *userRepository) UpdateLoginInfo(ctx context.Context, id uint64, ip string) error {
	return r.db.WithContext(ctx).Model(&model.SysUser{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_time": gorm.Expr("NOW()"),
			"last_login_ip":   ip,
		}).Error
}

// GetUserRoles 获取用户角色
func (r *userRepository) GetUserRoles(ctx context.Context, userID uint64) ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.db.WithContext(ctx).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// AssignRoles 分配角色
func (r *userRepository) AssignRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除原有角色
		if err := tx.Where("user_id = ?", userID).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}
		// 添加新角色
		for _, roleID := range roleIDs {
			userRole := model.SysUserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

