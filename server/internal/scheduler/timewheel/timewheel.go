package timewheel

import (
	"container/list"
	"sync"
	"time"

	"distributed-scheduler/pkg/logger"
)

// Task 定时任务
type Task struct {
	delay    time.Duration // 延迟时间
	circle   int           // 圈数
	key      string        // 唯一标识
	callback func()        // 回调函数
}

// TimeWheel 时间轮
type TimeWheel struct {
	interval    time.Duration // 时间间隔
	slots       []*list.List  // 时间槽
	slotNum     int           // 槽数量
	currentPos  int           // 当前位置
	ticker      *time.Ticker  // 定时器
	taskMap     map[string]*TaskElement // 任务映射
	addTaskCh   chan *Task    // 添加任务通道
	removeTaskCh chan string  // 移除任务通道
	stopCh      chan struct{} // 停止信号
	mu          sync.RWMutex
}

// TaskElement 任务元素(用于快速删除)
type TaskElement struct {
	task *Task
	pos  int
	elem *list.Element
}

// NewTimeWheel 创建时间轮
// interval: 时间间隔(如1秒)
// slotNum: 槽数量(如3600，则可支持1小时内的定时任务)
func NewTimeWheel(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}

	tw := &TimeWheel{
		interval:     interval,
		slots:        make([]*list.List, slotNum),
		slotNum:      slotNum,
		currentPos:   0,
		taskMap:      make(map[string]*TaskElement),
		addTaskCh:    make(chan *Task, 1000),
		removeTaskCh: make(chan string, 1000),
		stopCh:       make(chan struct{}),
	}

	// 初始化槽
	for i := 0; i < slotNum; i++ {
		tw.slots[i] = list.New()
	}

	return tw
}

// Start 启动时间轮
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
	logger.Info("时间轮启动成功")
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	close(tw.stopCh)
	if tw.ticker != nil {
		tw.ticker.Stop()
	}
	logger.Info("时间轮已停止")
}

// run 运行时间轮
func (tw *TimeWheel) run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tick()
		case task := <-tw.addTaskCh:
			tw.addTask(task)
		case key := <-tw.removeTaskCh:
			tw.removeTask(key)
		case <-tw.stopCh:
			return
		}
	}
}

// tick 时间轮滴答
func (tw *TimeWheel) tick() {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	// 获取当前槽
	slot := tw.slots[tw.currentPos]

	// 遍历当前槽的所有任务
	for e := slot.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}

		// 执行任务
		go task.callback()

		// 移除任务
		next := e.Next()
		slot.Remove(e)
		delete(tw.taskMap, task.key)
		e = next
	}

	// 移动指针
	tw.currentPos = (tw.currentPos + 1) % tw.slotNum
}

// AddTask 添加任务
func (tw *TimeWheel) AddTask(delay time.Duration, key string, callback func()) {
	if delay <= 0 {
		return
	}

	task := &Task{
		delay:    delay,
		key:      key,
		callback: callback,
	}

	select {
	case tw.addTaskCh <- task:
	default:
		// 通道满了，直接添加
		tw.mu.Lock()
		tw.addTask(task)
		tw.mu.Unlock()
	}
}

// addTask 内部添加任务
func (tw *TimeWheel) addTask(task *Task) {
	// 如果已存在相同key的任务，先移除
	if _, ok := tw.taskMap[task.key]; ok {
		tw.removeTask(task.key)
	}

	// 计算位置和圈数
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	// 添加到对应槽
	elem := tw.slots[pos].PushBack(task)
	tw.taskMap[task.key] = &TaskElement{
		task: task,
		pos:  pos,
		elem: elem,
	}
}

// RemoveTask 移除任务
func (tw *TimeWheel) RemoveTask(key string) {
	select {
	case tw.removeTaskCh <- key:
	default:
		tw.mu.Lock()
		tw.removeTask(key)
		tw.mu.Unlock()
	}
}

// removeTask 内部移除任务
func (tw *TimeWheel) removeTask(key string) {
	taskElem, ok := tw.taskMap[key]
	if !ok {
		return
	}

	tw.slots[taskElem.pos].Remove(taskElem.elem)
	delete(tw.taskMap, key)
}

// getPositionAndCircle 计算位置和圈数
func (tw *TimeWheel) getPositionAndCircle(delay time.Duration) (int, int) {
	delaySeconds := int(delay / tw.interval)
	circle := delaySeconds / tw.slotNum
	pos := (tw.currentPos + delaySeconds) % tw.slotNum
	return pos, circle
}

// HasTask 检查任务是否存在
func (tw *TimeWheel) HasTask(key string) bool {
	tw.mu.RLock()
	defer tw.mu.RUnlock()
	_, ok := tw.taskMap[key]
	return ok
}

// TaskCount 获取任务数量
func (tw *TimeWheel) TaskCount() int {
	tw.mu.RLock()
	defer tw.mu.RUnlock()
	return len(tw.taskMap)
}

