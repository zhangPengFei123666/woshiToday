package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/common/response"
	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/service"
)

// ExecutorHandler 执行器处理器
type ExecutorHandler struct {
	executorService service.ExecutorService
}

// NewExecutorHandler 创建执行器处理器
func NewExecutorHandler() *ExecutorHandler {
	return &ExecutorHandler{
		executorService: service.NewExecutorService(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	AppName       string `json:"app_name" binding:"required,max=64"`
	Host          string `json:"host" binding:"required"`
	Port          uint   `json:"port" binding:"required,min=1,max=65535"`
	MaxConcurrent uint   `json:"max_concurrent" binding:"min=1"`
}

// Register 注册执行器
// @Summary 注册执行器
// @Tags 执行器管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求"
// @Success 200 {object} response.Response{data=model.ExecutorNode}
// @Router /api/v1/executor/register [post]
func (h *ExecutorHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if req.MaxConcurrent == 0 {
		req.MaxConcurrent = 100
	}

	node, err := h.executorService.Register(c.Request.Context(), req.AppName, req.Host, req.Port, req.MaxConcurrent)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, node)
}

// UnregisterRequest 注销请求
type UnregisterRequest struct {
	ExecutorID string `json:"executor_id" binding:"required"`
}

// Unregister 注销执行器
// @Summary 注销执行器
// @Tags 执行器管理
// @Accept json
// @Produce json
// @Param request body UnregisterRequest true "注销请求"
// @Success 200 {object} response.Response
// @Router /api/v1/executor/unregister [post]
func (h *ExecutorHandler) Unregister(c *gin.Context) {
	var req UnregisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if err := h.executorService.Unregister(c.Request.Context(), req.ExecutorID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Heartbeat 心跳
// @Summary 执行器心跳
// @Tags 执行器管理
// @Accept json
// @Produce json
// @Param request body model.ExecutorHeartbeat true "心跳请求"
// @Success 200 {object} response.Response
// @Router /api/v1/executor/heartbeat [post]
func (h *ExecutorHandler) Heartbeat(c *gin.Context) {
	var req model.ExecutorHeartbeat
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if err := h.executorService.Heartbeat(c.Request.Context(), &req); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 获取执行器详情
// @Summary 获取执行器详情
// @Tags 执行器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "执行器ID"
// @Success 200 {object} response.Response{data=model.ExecutorNode}
// @Router /api/v1/executor/{id} [get]
func (h *ExecutorHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ParamError(c, "无效的执行器ID")
		return
	}

	node, err := h.executorService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrExecutorNotFound {
			response.NotFound(c, "执行器不存在")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, node)
}

// ExecutorListRequest 执行器列表请求
type ExecutorListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	GroupID  uint64 `form:"group_id"`
	Status   int8   `form:"status" binding:"min=-1,max=1"`
}

// List 执行器列表
// @Summary 执行器列表
// @Tags 执行器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param group_id query int false "任务组ID"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=response.PageResult}
// @Router /api/v1/executor [get]
func (h *ExecutorHandler) List(c *gin.Context) {
	var req ExecutorListRequest
	req.Status = -1
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

	nodes, total, err := h.executorService.List(c.Request.Context(), req.Page, req.PageSize, req.GroupID, req.Status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, nodes, total, req.Page, req.PageSize)
}

// GetOnlineByGroupID 获取任务组的在线执行器
// @Summary 获取任务组的在线执行器
// @Tags 执行器管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param group_id query int true "任务组ID"
// @Success 200 {object} response.Response
// @Router /api/v1/executor/online [get]
func (h *ExecutorHandler) GetOnlineByGroupID(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Query("group_id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务组ID")
		return
	}

	nodes, err := h.executorService.GetOnlineByGroupID(c.Request.Context(), groupID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nodes)
}

