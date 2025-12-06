package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 响应码常量
const (
	CodeSuccess      = 0
	CodeError        = -1
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500

	// 业务错误码
	CodeParamError     = 10001 // 参数错误
	CodeUserNotFound   = 10002 // 用户不存在
	CodePasswordError  = 10003 // 密码错误
	CodeUserDisabled   = 10004 // 用户已禁用
	CodeTokenExpired   = 10005 // Token过期
	CodeTokenInvalid   = 10006 // Token无效
	CodePermDenied     = 10007 // 权限不足
	CodeTaskNotFound   = 10008 // 任务不存在
	CodeGroupNotFound  = 10009 // 任务组不存在
	CodeExecutorError  = 10010 // 执行器错误
	CodeScheduleError  = 10011 // 调度错误
	CodeDuplicateEntry = 10012 // 重复记录
)

// 响应消息
var codeMsg = map[int]string{
	CodeSuccess:        "success",
	CodeError:          "error",
	CodeUnauthorized:   "未授权",
	CodeForbidden:      "禁止访问",
	CodeNotFound:       "资源不存在",
	CodeServerError:    "服务器内部错误",
	CodeParamError:     "参数错误",
	CodeUserNotFound:   "用户不存在",
	CodePasswordError:  "密码错误",
	CodeUserDisabled:   "用户已禁用",
	CodeTokenExpired:   "Token已过期",
	CodeTokenInvalid:   "Token无效",
	CodePermDenied:     "权限不足",
	CodeTaskNotFound:   "任务不存在",
	CodeGroupNotFound:  "任务组不存在",
	CodeExecutorError:  "执行器错误",
	CodeScheduleError:  "调度错误",
	CodeDuplicateEntry: "重复记录",
}

// GetCodeMsg 获取响应码对应的消息
func GetCodeMsg(code int) string {
	if msg, ok := codeMsg[code]; ok {
		return msg
	}
	return "未知错误"
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应(自定义消息)
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PageResult{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	if message == "" {
		message = GetCodeMsg(code)
	}
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 错误响应(带数据)
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	if message == "" {
		message = GetCodeMsg(code)
	}
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	if message == "" {
		message = GetCodeMsg(CodeParamError)
	}
	Error(c, CodeParamError, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = GetCodeMsg(CodeUnauthorized)
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = GetCodeMsg(CodeForbidden)
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
	})
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = GetCodeMsg(CodeNotFound)
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = GetCodeMsg(CodeServerError)
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeServerError,
		Message: message,
	})
}

