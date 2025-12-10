package first

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

// erfen search
func search(nums []int, target int) int {
	l, r := 0, len(nums)-1
	res := -1
	for l <= r {
		mid := l + (r-l)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			r = mid - 1
		} else {
			l = mid + 1
		}
		res = mid
	}
	if res == 0 {
		return 0
	}
	return res + 1
}

func searchRange(nums []int, target int) []int {
	l := firstAndSecondErfen(nums, target)
	if l == len(nums) || nums[l] != target {
		return []int{-1, -1}
	}
	r := firstAndSecondErfen(nums, target+1) - 1
	return []int{l, r}
}

func firstAndSecondErfen(nums []int, target int) int {
	l, r := 0, len(nums)
	for l < r {
		mid := l + (r-l)/2
		if nums[mid] >= target {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

func mySqrt(x int) int {
	l, r := 0, x
	ans := -1

	for l <= r {
		mid := l + (r-l)/2
		if mid*mid <= x {
			ans = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return ans
}

// yichu yuansu
func removeElement(nums []int, val int) int {
	show := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] == val {
			continue
		}
		nums[show] = nums[fast]
		show++
	}
	return show
}

func removeElementxiangxiang(nums []int, val int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		for l <= r && nums[l] != val {
			l++
		}
		for l <= r && nums[r] == val {
			r--
		}

		if l <= r {
			nums[l] = nums[r]
			l++
			r--
		}
	}
	return l
}

func removeNumsNotDizeng(nums []int) int {
	show := 1
	if len(nums) == 0 {
		return 0
	}
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[fast-1] {
			nums[show] = nums[fast]
			show++
		}
	}
	return show
}

// yidong zero
func moveZeroes(nums []int) {
	show := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			nums[show], nums[fast] = nums[fast], nums[show]
			show++
		}
	}
}
