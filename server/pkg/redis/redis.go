package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"distributed-scheduler/internal/config"
	"distributed-scheduler/pkg/logger"
)

var Client *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	logger.Info("Redis连接成功")
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// Close 关闭Redis连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Exists 判断键是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
	return Client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return Client.TTL(ctx, key).Result()
}

// Incr 自增
func Incr(ctx context.Context, key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// Decr 自减
func Decr(ctx context.Context, key string) (int64, error) {
	return Client.Decr(ctx, key).Result()
}

// HSet 设置哈希字段
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return Client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段
func HGet(ctx context.Context, key, field string) (string, error) {
	return Client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return Client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return Client.HDel(ctx, key, fields...).Err()
}

// LPush 左侧入队
func LPush(ctx context.Context, key string, values ...interface{}) error {
	return Client.LPush(ctx, key, values...).Err()
}

// RPush 右侧入队
func RPush(ctx context.Context, key string, values ...interface{}) error {
	return Client.RPush(ctx, key, values...).Err()
}

// LPop 左侧出队
func LPop(ctx context.Context, key string) (string, error) {
	return Client.LPop(ctx, key).Result()
}

// RPop 右侧出队
func RPop(ctx context.Context, key string) (string, error) {
	return Client.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
	return Client.LLen(ctx, key).Result()
}

// SAdd 添加集合成员
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return Client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(ctx context.Context, key string) ([]string, error) {
	return Client.SMembers(ctx, key).Result()
}

// SRem 移除集合成员
func SRem(ctx context.Context, key string, members ...interface{}) error {
	return Client.SRem(ctx, key, members...).Err()
}

// ZAdd 添加有序集合成员
func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return Client.ZAdd(ctx, key, members...).Err()
}

// ZRange 按排名范围获取有序集合成员
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return Client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeByScore 按分数范围获取有序集合成员
func ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return Client.ZRangeByScore(ctx, key, opt).Result()
}

// ZRem 移除有序集合成员
func ZRem(ctx context.Context, key string, members ...interface{}) error {
	return Client.ZRem(ctx, key, members...).Err()
}

