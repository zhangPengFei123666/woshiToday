package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/common/response"
	"distributed-scheduler/internal/service"
)

// InstanceHandler 任务实例处理器
type InstanceHandler struct {
	instanceService service.InstanceService
}

// NewInstanceHandler 创建任务实例处理器
func NewInstanceHandler() *InstanceHandler {
	return &InstanceHandler{
		instanceService: service.NewInstanceService(),
	}
}

// GetByID 获取任务实例详情
// @Summary 获取任务实例详情
// @Tags 执行记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "实例ID"
// @Success 200 {object} response.Response{data=model.TaskInstance}
// @Router /api/v1/instance/{id} [get]
func (h *InstanceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的实例ID")
		return
	}

	instance, err := h.instanceService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrInstanceNotFound {
			response.NotFound(c, "任务实例不存在")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, instance)
}

// InstanceListRequest 实例列表请求
type InstanceListRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	PageSize  int    `form:"page_size" binding:"min=1,max=100"`
	TaskID    uint64 `form:"task_id"`
	Status    int8   `form:"status" binding:"min=-1,max=5"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}

// List 任务实例列表
// @Summary 任务实例列表
// @Tags 执行记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param task_id query int false "任务ID"
// @Param status query int false "状态"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=response.PageResult}
// @Router /api/v1/instance [get]
func (h *InstanceHandler) List(c *gin.Context) {
	var req InstanceListRequest
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

	var startTime, endTime *time.Time
	if req.StartTime != "" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
		if err == nil {
			startTime = &t
		}
	}
	if req.EndTime != "" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local)
		if err == nil {
			endTime = &t
		}
	}

	instances, total, err := h.instanceService.List(c.Request.Context(), req.Page, req.PageSize, req.TaskID, req.Status, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, instances, total, req.Page, req.PageSize)
}

// Cancel 取消任务实例
// @Summary 取消任务实例
// @Tags 执行记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "实例ID"
// @Success 200 {object} response.Response
// @Router /api/v1/instance/{id}/cancel [post]
func (h *InstanceHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的实例ID")
		return
	}

	if err := h.instanceService.Cancel(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Retry 重试任务实例
// @Summary 重试任务实例
// @Tags 执行记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "实例ID"
// @Success 200 {object} response.Response{data=model.TaskInstance}
// @Router /api/v1/instance/{id}/retry [post]
func (h *InstanceHandler) Retry(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的实例ID")
		return
	}

	instance, err := h.instanceService.Retry(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, instance)
}

// GetLogs 获取任务实例日志
// @Summary 获取任务实例日志
// @Tags 执行记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "实例ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult}
// @Router /api/v1/instance/{id}/logs [get]
func (h *InstanceHandler) GetLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的实例ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	logs, total, err := h.instanceService.GetLogs(c.Request.Context(), id, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetStatistics 获取统计信息
// @Summary 获取统计信息
// @Tags 执行记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param task_id query int false "任务ID"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=service.InstanceStatistics}
// @Router /api/v1/instance/statistics [get]
func (h *InstanceHandler) GetStatistics(c *gin.Context) {
	taskID, _ := strconv.ParseUint(c.Query("task_id"), 10, 64)

	startTimeStr := c.DefaultQuery("start_time", time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05"))
	endTimeStr := c.DefaultQuery("end_time", time.Now().Format("2006-01-02 15:04:05"))

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endTimeStr, time.Local)

	stats, err := h.instanceService.GetStatistics(c.Request.Context(), taskID, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetRecentInstances 获取最近的实例
// @Summary 获取最近的实例
// @Tags 执行记录
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "数量"
// @Success 200 {object} response.Response
// @Router /api/v1/instance/recent [get]
func (h *InstanceHandler) GetRecentInstances(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	instances, err := h.instanceService.GetRecentInstances(c.Request.Context(), limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, instances)
}

