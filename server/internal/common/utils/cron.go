package utils

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// CronParser Cron解析器
var CronParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

// ValidateCron 验证Cron表达式
func ValidateCron(expr string) error {
	_, err := CronParser.Parse(expr)
	if err != nil {
		return fmt.Errorf("无效的Cron表达式: %w", err)
	}
	return nil
}

// GetNextTriggerTime 获取下次触发时间
func GetNextTriggerTime(expr string, from time.Time) (time.Time, error) {
	schedule, err := CronParser.Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return schedule.Next(from), nil
}

// GetNextNTriggerTimes 获取下N次触发时间
func GetNextNTriggerTimes(expr string, from time.Time, n int) ([]time.Time, error) {
	schedule, err := CronParser.Parse(expr)
	if err != nil {
		return nil, err
	}

	times := make([]time.Time, 0, n)
	current := from
	for i := 0; i < n; i++ {
		next := schedule.Next(current)
		times = append(times, next)
		current = next
	}
	return times, nil
}

// CronDescription Cron表达式描述
func CronDescription(expr string) string {
	// 简单的Cron表达式描述
	// 格式: 秒 分 时 日 月 周
	// 示例: "0 0 * * * *" -> 每小时执行一次
	// 这里只做简单描述，完整实现需要更复杂的逻辑
	return fmt.Sprintf("Cron表达式: %s", expr)
}

