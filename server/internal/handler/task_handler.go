package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/common/response"
	"distributed-scheduler/internal/middleware"
	"distributed-scheduler/internal/model"
	"distributed-scheduler/internal/service"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	taskService service.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService: service.NewTaskService(),
	}
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	GroupID         uint64   `json:"group_id" binding:"required"`
	Name            string   `json:"name" binding:"required,max=128"`
	Description     string   `json:"description" binding:"max=512"`
	Cron            string   `json:"cron" binding:"required"`
	ExecutorType    string   `json:"executor_type" binding:"required,oneof=HTTP GRPC SCRIPT"`
	ExecutorHandler string   `json:"executor_handler" binding:"required,max=256"`
	ExecutorParam   string   `json:"executor_param"`
	RouteStrategy   string   `json:"route_strategy" binding:"omitempty,oneof=ROUND_ROBIN RANDOM CONSISTENT_HASH LEAST_FREQUENTLY_USED LEAST_RECENTLY_USED FAILOVER SHARDING_BROADCAST"`
	BlockStrategy   string   `json:"block_strategy" binding:"omitempty,oneof=SERIAL_EXECUTION DISCARD_LATER COVER_EARLY"`
	ShardNum        uint     `json:"shard_num"`
	RetryCount      uint     `json:"retry_count"`
	RetryInterval   uint     `json:"retry_interval"`
	Timeout         uint     `json:"timeout"`
	AlarmEmail      string   `json:"alarm_email"`
	Priority        int      `json:"priority"`
	DependencyIDs   []uint64 `json:"dependency_ids"`
}

// Create 创建任务
// @Summary 创建任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateTaskRequest true "创建任务请求"
// @Success 200 {object} response.Response{data=model.Task}
// @Router /api/v1/task [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	task := &model.Task{
		GroupID:         req.GroupID,
		Name:            req.Name,
		Description:     req.Description,
		Cron:            req.Cron,
		ExecutorType:    req.ExecutorType,
		ExecutorHandler: req.ExecutorHandler,
		ExecutorParam:   req.ExecutorParam,
		RouteStrategy:   req.RouteStrategy,
		BlockStrategy:   req.BlockStrategy,
		ShardNum:        req.ShardNum,
		RetryCount:      req.RetryCount,
		RetryInterval:   req.RetryInterval,
		Timeout:         req.Timeout,
		AlarmEmail:      req.AlarmEmail,
		Priority:        req.Priority,
		Status:          model.TaskStatusDisabled,
		CreatedBy:       middleware.GetUserID(c),
	}

	// 设置默认值
	if task.RouteStrategy == "" {
		task.RouteStrategy = model.RouteStrategyRoundRobin
	}
	if task.BlockStrategy == "" {
		task.BlockStrategy = model.BlockStrategySerialExecution
	}
	if task.ShardNum == 0 {
		task.ShardNum = 1
	}

	if err := h.taskService.Create(c.Request.Context(), task); err != nil {
		switch err {
		case service.ErrInvalidCron:
			response.ParamError(c, "无效的Cron表达式")
		case service.ErrGroupNotFound:
			response.Error(c, response.CodeGroupNotFound, "")
		default:
			response.ServerError(c, err.Error())
		}
		return
	}

	response.Success(c, task)
}

// Update 更新任务
// @Summary 更新任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Param request body CreateTaskRequest true "更新任务请求"
// @Success 200 {object} response.Response
// @Router /api/v1/task/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务ID")
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 获取原任务
	task, err := h.taskService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			response.Error(c, response.CodeTaskNotFound, "")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	// 更新字段
	task.GroupID = req.GroupID
	task.Name = req.Name
	task.Description = req.Description
	task.Cron = req.Cron
	task.ExecutorType = req.ExecutorType
	task.ExecutorHandler = req.ExecutorHandler
	task.ExecutorParam = req.ExecutorParam
	task.RouteStrategy = req.RouteStrategy
	task.BlockStrategy = req.BlockStrategy
	task.ShardNum = req.ShardNum
	task.RetryCount = req.RetryCount
	task.RetryInterval = req.RetryInterval
	task.Timeout = req.Timeout
	task.AlarmEmail = req.AlarmEmail
	task.Priority = req.Priority

	if err := h.taskService.Update(c.Request.Context(), task); err != nil {
		if err == service.ErrInvalidCron {
			response.ParamError(c, "无效的Cron表达式")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, task)
}

// Delete 删除任务
// @Summary 删除任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/v1/task/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务ID")
		return
	}

	if err := h.taskService.Delete(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 获取任务详情
// @Summary 获取任务详情
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response{data=model.Task}
// @Router /api/v1/task/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务ID")
		return
	}

	task, err := h.taskService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			response.Error(c, response.CodeTaskNotFound, "")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, task)
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	GroupID  uint64 `form:"group_id"`
	Keyword  string `form:"keyword"`
	Status   int8   `form:"status" binding:"min=-1,max=1"`
}

// List 任务列表
// @Summary 任务列表
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param group_id query int false "任务组ID"
// @Param keyword query string false "关键字"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=response.PageResult}
// @Router /api/v1/task [get]
func (h *TaskHandler) List(c *gin.Context) {
	var req TaskListRequest
	req.Status = -1 // 默认不过滤状态
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

	tasks, total, err := h.taskService.List(c.Request.Context(), req.Page, req.PageSize, req.GroupID, req.Keyword, req.Status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, tasks, total, req.Page, req.PageSize)
}

// Start 启动任务
// @Summary 启动任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/v1/task/{id}/start [post]
func (h *TaskHandler) Start(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务ID")
		return
	}

	if err := h.taskService.Start(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Stop 停止任务
// @Summary 停止任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/v1/task/{id}/stop [post]
func (h *TaskHandler) Stop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务ID")
		return
	}

	if err := h.taskService.Stop(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// TriggerRequest 手动触发请求
type TriggerRequest struct {
	Param string `json:"param"`
}

// Trigger 手动触发任务
// @Summary 手动触发任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Param request body TriggerRequest true "触发请求"
// @Success 200 {object} response.Response{data=model.TaskInstance}
// @Router /api/v1/task/{id}/trigger [post]
func (h *TaskHandler) Trigger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的任务ID")
		return
	}

	var req TriggerRequest
	_ = c.ShouldBindJSON(&req)

	instance, err := h.taskService.Trigger(c.Request.Context(), id, req.Param)
	if err != nil {
		if err == service.ErrTaskNotFound {
			response.Error(c, response.CodeTaskNotFound, "")
			return
		}
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, instance)
}

// NextTriggerTimesRequest 下次触发时间请求
type NextTriggerTimesRequest struct {
	Cron  string `form:"cron" binding:"required"`
	Count int    `form:"count" binding:"min=1,max=10"`
}

// GetNextTriggerTimes 获取下次触发时间
// @Summary 获取下次触发时间
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param cron query string true "Cron表达式"
// @Param count query int false "获取数量"
// @Success 200 {object} response.Response
// @Router /api/v1/task/next-trigger-times [get]
func (h *TaskHandler) GetNextTriggerTimes(c *gin.Context) {
	var req NextTriggerTimesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if req.Count <= 0 {
		req.Count = 5
	}

	times, err := h.taskService.GetNextTriggerTimes(c.Request.Context(), req.Cron, req.Count)
	if err != nil {
		response.ParamError(c, "无效的Cron表达式")
		return
	}

	response.Success(c, times)
}

