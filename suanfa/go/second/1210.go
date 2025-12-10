package second

type LRUList struct {
	key, value int
	pre, next  *LRUList
}

func initLRUList(k, v int) *LRUList {
	return &LRUList{key: k, value: v}
}

type LRUCache struct {
	size, capacity int
	head, tail     *LRUList
	mapCache       map[int]*LRUList
}

func Constructor(capacity int) LRUCache {
	l := LRUCache{
		size:     0,
		capacity: capacity,
		head:     initLRUList(0, 0),
		tail:     initLRUList(0, 0),
		mapCache: map[int]*LRUList{},
	}
	l.head.next = l.tail
	l.tail.pre = l.head
	return l
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
		node := initLRUList(key, value)
		this.mapCache[key] = node
		this.AddtoHead(node)
		this.size++
		if this.size > this.capacity {
			tail := this.TailNode()
			delete(this.mapCache, tail.key)
			this.size--
		}
		return
	}
	node := this.mapCache[key]
	node.value = value
	this.AddtoHead(node)
	return
}

func (this *LRUCache) AddtoHead(node *LRUList) {
	node.pre = this.head
	node.next = this.head.next
	this.head.next.pre = node
	this.head.next = node
}

func (this *LRUCache) RemoveNode(node *LRUList) {
	node.pre.next = node.next
	node.next.pre = node.pre
}

func (this *LRUCache) MoveToHead(node *LRUList) {
	this.RemoveNode(node)
	this.AddtoHead(node)
}

func (this *LRUCache) TailNode() *LRUList {
	tailNode := this.tail.pre
	this.RemoveNode(tailNode)
	return tailNode
}
