package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"distributed-scheduler/internal/common/response"
)

// RateLimiter 令牌桶限流器
type RateLimiter struct {
	rate       int           // 每秒产生的令牌数
	capacity   int           // 桶容量
	tokens     int           // 当前令牌数
	lastUpdate time.Time     // 上次更新时间
	mu         sync.Mutex    // 互斥锁
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate, capacity int) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastUpdate: time.Now(),
	}
}

// Allow 是否允许请求
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// 计算需要添加的令牌数
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	tokensToAdd := int(elapsed * float64(rl.rate))

	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.lastUpdate = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	limiters map[string]*RateLimiter
	rate     int
	capacity int
	mu       sync.Mutex
}

// NewIPRateLimiter 创建IP限流器
func NewIPRateLimiter(rate, capacity int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
		rate:     rate,
		capacity: capacity,
	}
}

// GetLimiter 获取指定IP的限流器
func (ipl *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	ipl.mu.Lock()
	defer ipl.mu.Unlock()

	limiter, exists := ipl.limiters[ip]
	if !exists {
		limiter = NewRateLimiter(ipl.rate, ipl.capacity)
		ipl.limiters[ip] = limiter
	}
	return limiter
}

// 全局IP限流器实例
var globalIPLimiter *IPRateLimiter

// InitRateLimiter 初始化限流器
func InitRateLimiter(rate, capacity int) {
	globalIPLimiter = NewIPRateLimiter(rate, capacity)
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if globalIPLimiter == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		limiter := globalIPLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			response.Error(c, response.CodeError, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

