package dag

import (
	"errors"
	"sync"
)

var (
	ErrCycleDetected = errors.New("检测到循环依赖")
	ErrNodeNotFound  = errors.New("节点不存在")
)

// Node DAG节点
type Node struct {
	ID       uint64
	Name     string
	Data     interface{}
	InDegree int           // 入度
	Children []*Node       // 子节点
	Parents  []*Node       // 父节点
}

// DAG 有向无环图
type DAG struct {
	nodes map[uint64]*Node
	mu    sync.RWMutex
}

// NewDAG 创建DAG
func NewDAG() *DAG {
	return &DAG{
		nodes: make(map[uint64]*Node),
	}
}

// AddNode 添加节点
func (d *DAG) AddNode(id uint64, name string, data interface{}) *Node {
	d.mu.Lock()
	defer d.mu.Unlock()

	if node, exists := d.nodes[id]; exists {
		return node
	}

	node := &Node{
		ID:       id,
		Name:     name,
		Data:     data,
		Children: make([]*Node, 0),
		Parents:  make([]*Node, 0),
	}
	d.nodes[id] = node
	return node
}

// AddEdge 添加边(from -> to)
func (d *DAG) AddEdge(fromID, toID uint64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	from, exists := d.nodes[fromID]
	if !exists {
		return ErrNodeNotFound
	}

	to, exists := d.nodes[toID]
	if !exists {
		return ErrNodeNotFound
	}

	// 检查是否会形成环
	if d.wouldCreateCycle(from, to) {
		return ErrCycleDetected
	}

	from.Children = append(from.Children, to)
	to.Parents = append(to.Parents, from)
	to.InDegree++

	return nil
}

// wouldCreateCycle 检查添加边是否会形成环
func (d *DAG) wouldCreateCycle(from, to *Node) bool {
	// 如果to能到达from，则会形成环
	visited := make(map[uint64]bool)
	return d.canReach(to, from, visited)
}

// canReach 检查from是否能到达to
func (d *DAG) canReach(from, to *Node, visited map[uint64]bool) bool {
	if from.ID == to.ID {
		return true
	}

	if visited[from.ID] {
		return false
	}
	visited[from.ID] = true

	for _, child := range from.Children {
		if d.canReach(child, to, visited) {
			return true
		}
	}

	return false
}

// GetNode 获取节点
func (d *DAG) GetNode(id uint64) (*Node, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	node, exists := d.nodes[id]
	return node, exists
}

// RemoveNode 移除节点
func (d *DAG) RemoveNode(id uint64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	node, exists := d.nodes[id]
	if !exists {
		return
	}

	// 移除所有指向该节点的边
	for _, parent := range node.Parents {
		for i, child := range parent.Children {
			if child.ID == id {
				parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
				break
			}
		}
	}

	// 移除所有从该节点出发的边
	for _, child := range node.Children {
		child.InDegree--
		for i, parent := range child.Parents {
			if parent.ID == id {
				child.Parents = append(child.Parents[:i], child.Parents[i+1:]...)
				break
			}
		}
	}

	delete(d.nodes, id)
}

// TopologicalSort 拓扑排序
func (d *DAG) TopologicalSort() ([]*Node, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// 复制入度信息
	inDegree := make(map[uint64]int)
	for id, node := range d.nodes {
		inDegree[id] = node.InDegree
	}

	// 找出所有入度为0的节点
	queue := make([]*Node, 0)
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, d.nodes[id])
		}
	}

	result := make([]*Node, 0, len(d.nodes))

	for len(queue) > 0 {
		// 取出队首节点
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// 减少所有子节点的入度
		for _, child := range node.Children {
			inDegree[child.ID]--
			if inDegree[child.ID] == 0 {
				queue = append(queue, child)
			}
		}
	}

	// 如果结果数量不等于节点数量，说明存在环
	if len(result) != len(d.nodes) {
		return nil, ErrCycleDetected
	}

	return result, nil
}

// GetRoots 获取所有根节点(入度为0)
func (d *DAG) GetRoots() []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	roots := make([]*Node, 0)
	for _, node := range d.nodes {
		if node.InDegree == 0 {
			roots = append(roots, node)
		}
	}
	return roots
}

// GetLeaves 获取所有叶子节点(出度为0)
func (d *DAG) GetLeaves() []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	leaves := make([]*Node, 0)
	for _, node := range d.nodes {
		if len(node.Children) == 0 {
			leaves = append(leaves, node)
		}
	}
	return leaves
}

// NodeCount 获取节点数量
func (d *DAG) NodeCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.nodes)
}

// GetExecutableNodes 获取可执行的节点(所有父节点都已完成)
func (d *DAG) GetExecutableNodes(completedIDs map[uint64]bool) []*Node {
	d.mu.RLock()
	defer d.mu.RUnlock()

	executable := make([]*Node, 0)
	for _, node := range d.nodes {
		// 跳过已完成的节点
		if completedIDs[node.ID] {
			continue
		}

		// 检查所有父节点是否都已完成
		allParentsDone := true
		for _, parent := range node.Parents {
			if !completedIDs[parent.ID] {
				allParentsDone = false
				break
			}
		}

		if allParentsDone {
			executable = append(executable, node)
		}
	}

	return executable
}

