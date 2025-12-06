package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"distributed-scheduler/internal/common/utils"
	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/repository"
)

var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrPasswordInvalid = errors.New("密码错误")
	ErrUserDisabled    = errors.New("用户已禁用")
	ErrUsernameExists  = errors.New("用户名已存在")
)

// UserService 用户服务接口
type UserService interface {
	Login(ctx context.Context, username, password, ip string) (string, *model.SysUser, error)
	Create(ctx context.Context, user *model.SysUser, password string) error
	Update(ctx context.Context, user *model.SysUser) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*model.SysUser, error)
	List(ctx context.Context, page, pageSize int, keyword string) ([]*model.SysUser, int64, error)
	ChangePassword(ctx context.Context, id uint64, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, id uint64, newPassword string) error
	UpdateStatus(ctx context.Context, id uint64, status int8) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService() UserService {
	return &userService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, username, password, ip string) (string, *model.SysUser, error) {
	// 查询用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrUserNotFound
		}
		return "", nil, err
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return "", nil, ErrPasswordInvalid
	}

	// 检查状态
	if user.Status != model.UserStatusEnabled {
		return "", nil, ErrUserDisabled
	}

	// 获取角色
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return "", nil, err
	}

	roleCode := ""
	if len(roles) > 0 {
		roleCode = roles[0].Code
		user.Roles = roles
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, roleCode)
	if err != nil {
		return "", nil, err
	}

	// 更新登录信息
	_ = s.userRepo.UpdateLoginInfo(ctx, user.ID, ip)

	return token, user, nil
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, user *model.SysUser, password string) error {
	// 检查用户名是否存在
	existing, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return ErrUsernameExists
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Create(ctx, user)
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, user *model.SysUser) error {
	return s.userRepo.Update(ctx, user)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint64) error {
	return s.userRepo.Delete(ctx, id)
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint64) (*model.SysUser, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 获取角色
	roles, err := s.userRepo.GetUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return user, nil
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, page, pageSize int, keyword string) ([]*model.SysUser, int64, error) {
	return s.userRepo.List(ctx, page, pageSize, keyword)
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, id uint64, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(oldPassword, user.Password) {
		return ErrPasswordInvalid
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Update(ctx, user)
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, id uint64, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Update(ctx, user)
}

// UpdateStatus 更新用户状态
func (s *userService) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.Status = status
	return s.userRepo.Update(ctx, user)
}

