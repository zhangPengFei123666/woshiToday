package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/common/response"
	"distributed-scheduler/internal/middleware"
	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/service"
)

// GroupHandler 任务组处理器
type GroupHandler struct {
	groupService service.TaskGroupService
}

// NewGroupHandler 创建任务组处理器
func NewGroupHandler() *GroupHandler {
	return &GroupHandler{
		groupService: service.NewTaskGroupService(),
	}
}

// CreateGroupRequest 创建任务组请求
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description" binding:"max=512"`
	AppName     string `json:"app_name" binding:"required,max=64"`
}

// Create 创建任务组
// @Summary 创建任务组
// @Tags 任务组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateGroupRequest true "创建任务组请求"
// @Success 200 {object} response.Response{data=model.TaskGroup}
// @Router /api/v1/group [post]
func (h *GroupHandler) Create(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	group := &model.TaskGroup{
		Name:        req.Name,
		Description: req.Description,
		AppName:     req.AppName,
		Status:      1,
		CreatedBy:   middleware.GetUserID(c),
	}

	if err := h.groupService.Create(c.Request.Context(), group); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, group)
}

// Update 更新任务组
// @Summary 更新任务组
// @Tags 任务组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务组ID"
// @Param request body CreateGroupRequest true "更新任务组请求"
// @Success 200 {object} response.Response
// @Router /api/v1/group/{id} [put]
func (h *GroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务组ID")
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	group, err := h.groupService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	group.Name = req.Name
	group.Description = req.Description
	group.AppName = req.AppName

	if err := h.groupService.Update(c.Request.Context(), group); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, group)
}

// Delete 删除任务组
// @Summary 删除任务组
// @Tags 任务组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务组ID"
// @Success 200 {object} response.Response
// @Router /api/v1/group/{id} [delete]
func (h *GroupHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务组ID")
		return
	}

	if err := h.groupService.Delete(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 获取任务组详情
// @Summary 获取任务组详情
// @Tags 任务组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务组ID"
// @Success 200 {object} response.Response{data=model.TaskGroup}
// @Router /api/v1/group/{id} [get]
func (h *GroupHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务组ID")
		return
	}

	group, err := h.groupService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, group)
}

// GroupListRequest 任务组列表请求
type GroupListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// List 任务组列表
// @Summary 任务组列表
// @Tags 任务组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "关键字"
// @Success 200 {object} response.Response{data=response.PageResult}
// @Router /api/v1/group [get]
func (h *GroupHandler) List(c *gin.Context) {
	var req GroupListRequest
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

	groups, total, err := h.groupService.List(c.Request.Context(), req.Page, req.PageSize, req.Keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, groups, total, req.Page, req.PageSize)
}

// GetAll 获取所有任务组
// @Summary 获取所有任务组
// @Tags 任务组管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/v1/group/all [get]
func (h *GroupHandler) GetAll(c *gin.Context) {
	groups, err := h.groupService.GetAll(c.Request.Context())
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, groups)
}

