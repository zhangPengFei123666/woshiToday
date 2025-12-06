package router

import (
	"errors"
	"hash/crc32"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"distributed-scheduler/internal/model"
)

var (
	ErrNoAvailableExecutor = errors.New("没有可用的执行器")
)

// Strategy 路由策略接口
type Strategy interface {
	Select(executors []*model.ExecutorNode, param string) (*model.ExecutorNode, error)
}

// NewStrategy 根据策略名称创建路由策略
func NewStrategy(strategy string) Strategy {
	switch strategy {
	case model.RouteStrategyRandom:
		return &RandomStrategy{}
	case model.RouteStrategyConsistentHash:
		return &ConsistentHashStrategy{}
	case model.RouteStrategyLeastFrequentlyUsed:
		return &LFUStrategy{}
	case model.RouteStrategyLeastRecentlyUsed:
		return &LRUStrategy{}
	case model.RouteStrategyFailover:
		return &FailoverStrategy{}
	default:
		return &RoundRobinStrategy{}
	}
}

// RoundRobinStrategy 轮询策略
type RoundRobinStrategy struct {
	counter uint64
}

func (s *RoundRobinStrategy) Select(executors []*model.ExecutorNode, param string) (*model.ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	// 过滤在线且未过载的执行器
	available := filterAvailable(executors)
	if len(available) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	index := atomic.AddUint64(&s.counter, 1) % uint64(len(available))
	return available[index], nil
}

// RandomStrategy 随机策略
type RandomStrategy struct {
	r *rand.Rand
}

func (s *RandomStrategy) Select(executors []*model.ExecutorNode, param string) (*model.ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	available := filterAvailable(executors)
	if len(available) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	if s.r == nil {
		s.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return available[s.r.Intn(len(available))], nil
}

// ConsistentHashStrategy 一致性哈希策略
type ConsistentHashStrategy struct {
	mu       sync.RWMutex
	ring     *ConsistentHashRing
	replicas int
}

type ConsistentHashRing struct {
	hashSortedNodes []uint32
	circle          map[uint32]*model.ExecutorNode
}

func (s *ConsistentHashStrategy) Select(executors []*model.ExecutorNode, param string) (*model.ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	available := filterAvailable(executors)
	if len(available) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 重建哈希环
	s.replicas = 100
	s.ring = &ConsistentHashRing{
		circle: make(map[uint32]*model.ExecutorNode),
	}

	for _, node := range available {
		for i := 0; i < s.replicas; i++ {
			hash := hashKey(node.ID + string(rune(i)))
			s.ring.circle[hash] = node
			s.ring.hashSortedNodes = append(s.ring.hashSortedNodes, hash)
		}
	}
	sort.Slice(s.ring.hashSortedNodes, func(i, j int) bool {
		return s.ring.hashSortedNodes[i] < s.ring.hashSortedNodes[j]
	})

	// 查找节点
	hash := hashKey(param)
	idx := sort.Search(len(s.ring.hashSortedNodes), func(i int) bool {
		return s.ring.hashSortedNodes[i] >= hash
	})
	if idx >= len(s.ring.hashSortedNodes) {
		idx = 0
	}

	return s.ring.circle[s.ring.hashSortedNodes[idx]], nil
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// LFUStrategy 最不经常使用策略
type LFUStrategy struct{}

func (s *LFUStrategy) Select(executors []*model.ExecutorNode, param string) (*model.ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	available := filterAvailable(executors)
	if len(available) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	// 选择负载最小的执行器
	var minNode *model.ExecutorNode
	minLoad := uint(^uint(0))

	for _, node := range available {
		if node.CurrentLoad < minLoad {
			minLoad = node.CurrentLoad
			minNode = node
		}
	}

	return minNode, nil
}

// LRUStrategy 最近最少使用策略
type LRUStrategy struct {
	mu       sync.RWMutex
	lastUsed map[string]time.Time
}

func (s *LRUStrategy) Select(executors []*model.ExecutorNode, param string) (*model.ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	available := filterAvailable(executors)
	if len(available) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.lastUsed == nil {
		s.lastUsed = make(map[string]time.Time)
	}

	// 选择最久未使用的执行器
	var selectedNode *model.ExecutorNode
	oldestTime := time.Now()

	for _, node := range available {
		lastTime, ok := s.lastUsed[node.ID]
		if !ok {
			selectedNode = node
			break
		}
		if lastTime.Before(oldestTime) {
			oldestTime = lastTime
			selectedNode = node
		}
	}

	if selectedNode != nil {
		s.lastUsed[selectedNode.ID] = time.Now()
	}

	return selectedNode, nil
}

// FailoverStrategy 故障转移策略
type FailoverStrategy struct{}

func (s *FailoverStrategy) Select(executors []*model.ExecutorNode, param string) (*model.ExecutorNode, error) {
	if len(executors) == 0 {
		return nil, ErrNoAvailableExecutor
	}

	// 按权重排序，选择第一个可用的
	sorted := make([]*model.ExecutorNode, len(executors))
	copy(sorted, executors)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Weight > sorted[j].Weight
	})

	for _, node := range sorted {
		if node.IsOnline() && !node.IsOverload() {
			return node, nil
		}
	}

	return nil, ErrNoAvailableExecutor
}

// filterAvailable 过滤可用的执行器
func filterAvailable(executors []*model.ExecutorNode) []*model.ExecutorNode {
	available := make([]*model.ExecutorNode, 0, len(executors))
	for _, node := range executors {
		if node.IsOnline() && !node.IsOverload() {
			available = append(available, node)
		}
	}
	return available
}

