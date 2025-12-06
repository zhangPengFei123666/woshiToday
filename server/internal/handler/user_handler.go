package handler

import (
	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/common/response"
	"distributed-scheduler/internal/middleware"
	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string          `json:"token"`
	User  *model.SysUser `json:"user"`
}

// Login 用户登录
// @Summary 用户登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	token, user, err := h.userService.Login(c.Request.Context(), req.Username, req.Password, c.ClientIP())
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.Error(c, response.CodeUserNotFound, "")
		case service.ErrPasswordInvalid:
			response.Error(c, response.CodePasswordError, "")
		case service.ErrUserDisabled:
			response.Error(c, response.CodeUserDisabled, "")
		default:
			response.ServerError(c, err.Error())
		}
		return
	}

	response.Success(c, LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetCurrentUser 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=model.SysUser}
// @Router /api/v1/user/current [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, user)
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
}

// Create 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateUserRequest true "创建用户请求"
// @Success 200 {object} response.Response
// @Router /api/v1/user [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	user := &model.SysUser{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   model.UserStatusEnabled,
	}

	if err := h.userService.Create(c.Request.Context(), user, req.Password); err != nil {
		if err == service.ErrUsernameExists {
			response.Error(c, response.CodeDuplicateEntry, "用户名已存在")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// ListRequest 列表请求
type ListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// List 用户列表
// @Summary 用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "关键字"
// @Success 200 {object} response.Response{data=response.PageResult}
// @Router /api/v1/user [get]
func (h *UserHandler) List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, total, err := h.userService.List(c.Request.Context(), req.Page, req.PageSize, req.Keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, users, total, req.Page, req.PageSize)
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} response.Response
// @Router /api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err == service.ErrPasswordInvalid {
			response.Error(c, response.CodePasswordError, "原密码错误")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Logout 退出登录
// @Summary 退出登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	// TODO: 可以在Redis中记录Token黑名单
	response.Success(c, nil)
}

