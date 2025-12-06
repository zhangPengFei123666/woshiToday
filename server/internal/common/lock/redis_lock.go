package lock

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	pkgRedis "distributed-scheduler/pkg/redis"
)

var (
	ErrLockFailed   = errors.New("获取锁失败")
	ErrLockNotHeld  = errors.New("锁不属于当前持有者")
	ErrUnlockFailed = errors.New("释放锁失败")
)

// RedisLock Redis分布式锁
type RedisLock struct {
	key        string
	value      string
	expiration time.Duration
	client     *redis.Client
}

// NewRedisLock 创建Redis分布式锁
func NewRedisLock(key string, expiration time.Duration) *RedisLock {
	return &RedisLock{
		key:        "lock:" + key,
		value:      uuid.New().String(),
		expiration: expiration,
		client:     pkgRedis.GetClient(),
	}
}

// Lock 获取锁
func (l *RedisLock) Lock(ctx context.Context) error {
	success, err := l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
	if err != nil {
		return err
	}
	if !success {
		return ErrLockFailed
	}
	return nil
}

// TryLock 尝试获取锁(带超时)
func (l *RedisLock) TryLock(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := l.Lock(ctx); err == nil {
			return nil
		}
		// 等待一段时间后重试
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
			continue
		}
	}
	return ErrLockFailed
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	// Lua脚本保证原子性：只有锁的持有者才能释放锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return ErrLockNotHeld
	}
	return nil
}

// Refresh 刷新锁的过期时间
func (l *RedisLock) Refresh(ctx context.Context) error {
	// Lua脚本保证原子性：只有锁的持有者才能刷新过期时间
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value, l.expiration.Milliseconds()).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return ErrLockNotHeld
	}
	return nil
}

// WithLock 在锁保护下执行函数
func WithLock(ctx context.Context, key string, expiration time.Duration, fn func() error) error {
	lock := NewRedisLock(key, expiration)
	if err := lock.Lock(ctx); err != nil {
		return err
	}
	defer lock.Unlock(ctx)
	return fn()
}

// WithTryLock 尝试获取锁并执行函数
func WithTryLock(ctx context.Context, key string, expiration, timeout time.Duration, fn func() error) error {
	lock := NewRedisLock(key, expiration)
	if err := lock.TryLock(ctx, timeout); err != nil {
		return err
	}
	defer lock.Unlock(ctx)
	return fn()
}

// SchedulerLock 调度器锁(用于防止任务重复调度)
type SchedulerLock struct {
	*RedisLock
}

// NewSchedulerLock 创建调度器锁
func NewSchedulerLock(taskID uint64, triggerTime time.Time) *SchedulerLock {
	key := "scheduler:" + string(rune(taskID)) + ":" + triggerTime.Format("20060102150405")
	return &SchedulerLock{
		RedisLock: NewRedisLock(key, 5*time.Minute),
	}
}

// ExecutorLock 执行器锁(用于任务执行互斥)
type ExecutorLock struct {
	*RedisLock
}

// NewExecutorLock 创建执行器锁
func NewExecutorLock(instanceID uint64) *ExecutorLock {
	key := "executor:" + string(rune(instanceID))
	return &ExecutorLock{
		RedisLock: NewRedisLock(key, 10*time.Minute),
	}
}

