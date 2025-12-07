package main

func lengthStr(s string) int {
	ans := 0
	left := 0
	str := make([]bool, 128)
	for i, c := range s {
		for str[c] {
			str[s[left]] = false
			left++
		}
		str[c] = true
		ans = max(ans, i-left+1)
	}
	return ans
}

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */

type ListLRU struct {
	key, value int
	pre, next  *ListLRU
}
type LRUCache struct {
	size, capacity int
	head, tail     *ListLRU
	mapCache       map[int]*ListLRU
}

func initListLRU(key, value int) *ListLRU {
	return &ListLRU{
		key:   key,
		value: value,
	}
}
func Constructor(capacity int) LRUCache {
	l := LRUCache{
		size:     0,
		capacity: capacity,
		head:     initListLRU(0, 0),
		tail:     initListLRU(0, 0),
		mapCache: map[int]*ListLRU{},
	}
	l.head.next = l.tail
	l.tail.pre = l.head
	return l
}

func (this *LRUCache) AddToHead(node *ListLRU) {
	node.pre = this.head
	node.next = this.head.next
	this.head.next.pre = node
	this.head.next = node
}

func (this *LRUCache) RemoveNode(node *ListLRU) {
	node.pre.next = node.next
	node.next.pre = node.pre
}

func (this *LRUCache) MoveToHead(node *ListLRU) {
	this.RemoveNode(node)
	this.AddToHead(node)
}

func (this *LRUCache) TailNode() *ListLRU {
	tailNode := this.tail.pre
	this.RemoveNode(tailNode)
	return tailNode
}
func (this *LRUCache) Get(key int) int {
	if _, ok := this.mapCache[key]; !ok {
		return -1
	}
	node := this.mapCache[key]
	this.MoveToHead(node)
	return node.value
}

func (this *LRUCache) Put(key int, value int) {
	if _, ok := this.mapCache[key]; !ok {
		node := initListLRU(key, value)
		this.mapCache[key] = node
		this.size++
		this.AddToHead(node)
		if this.size > this.capacity {
			tailNode := this.TailNode()
			delete(this.mapCache, tailNode.key)
			this.size--
		}
	} else {
		node := this.mapCache[key]
		node.value = value
		this.MoveToHead(node)
	}
}
