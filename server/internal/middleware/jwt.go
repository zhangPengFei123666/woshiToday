package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/common/response"
	"distributed-scheduler/internal/common/utils"
)

// 上下文键
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
	ContextKeyRoleCode = "role_code"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		// 解析Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			if err == utils.ErrTokenExpired {
				response.Error(c, response.CodeTokenExpired, "Token已过期")
			} else {
				response.Error(c, response.CodeTokenInvalid, "Token无效")
			}
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyRoleCode, claims.RoleCode)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint64 {
	if userID, exists := c.Get(ContextKeyUserID); exists {
		return userID.(uint64)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(ContextKeyUsername); exists {
		return username.(string)
	}
	return ""
}

// GetRoleCode 从上下文获取角色编码
func GetRoleCode(c *gin.Context) string {
	if roleCode, exists := c.Get(ContextKeyRoleCode); exists {
		return roleCode.(string)
	}
	return ""
}

