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

func sortedSquares(nums []int) []int {
	// 有序数组的平方
	res := make([]int, len(nums))
	k := len(nums) - 1
	l, r := 0, len(nums)-1

	for l <= r {
		if nums[l]*nums[l] > nums[r]*nums[r] {
			res[k] = nums[l] * nums[l]
			k--
			l++
		} else {
			res[k] = nums[r] * nums[r]
			k--
			r--
		}
	}
	return res
}

func minSubArrayLen(target int, nums []int) int {
	// 长度最小的子数组
	i := 0
	sum := 0
	res := len(nums) + 1
	for j := 0; j < len(nums); j++ {
		sum += nums[j]
		for sum >= target {
			res = min(res, j-i+1)
			sum -= nums[i]
			i++
		}
	}
	if res == len(nums)+1 {
		return 0
	}
	return res
}

// 最小覆盖子串
func minWindow(s string, t string) string {
	left, ansLeft, ansRight := 0, -1, len(s)
	sStr := make([]int, 128)
	tStr := make([]int, 128)
	for _, c := range t {
		tStr[c]++
	}
	for right := 0; right < len(s); right++ {
		sStr[s[right]]++
		for isCover(sStr, tStr) {
			if right-left < ansRight-ansLeft {
				ansLeft = left
				ansRight = right
			}
			sStr[s[left]]--
			left++
		}
	}

	if ansLeft == -1 {
		return ""
	}
	return s[ansLeft : ansRight+1]
}

func isCover(s []int, t []int) bool {
	for i := 'a'; i <= 'z'; i++ {
		if s[i] < t[i] {
			return false
		}
	}

	for i := 'A'; i <= 'Z'; i++ {
		if s[i] < t[i] {
			return false
		}
	}
	return true
}

// 螺旋矩阵
func generateMatrix(n int) [][]int {
	left, right := 0, n-1
	top, bottom := 0, n-1
	num := 1
	tar := n * n
	res := make([][]int, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
	}
	for num <= tar {
		for i := left; i <= right; i++ {
			res[top][i] = num
			num++
		}
		top += 1
		for i := top; i <= bottom; i++ {
			res[i][right] = num
			num++
		}
		right -= 1
		for i := right; i >= left; i-- {
			res[bottom][i] = num
			num++
		}
		bottom -= 1
		for i := bottom; i >= top; i-- {
			res[i][left] = num
			num++
		}
		left += 1
	}
	return res
}

// 顺时针旋转返回螺旋顺序
func spiralOrder(matrix [][]int) []int {
	m, n := len(matrix), len(matrix[0])
	if m == 0 || n == 0 {
		return []int{}
	}
	tar := m * n
	num := 0
	res := make([]int, m*n)
	left, right := 0, n-1
	top, bottom := 0, m-1
	for num < tar {
		for i := left; i <= right && num < tar; i++ {
			res[num] = matrix[top][i]
			num++
		}
		top++
		for i := top; i <= bottom && num < tar; i++ {
			res[num] = matrix[i][right]
			num++
		}
		right--
		for i := right; i >= left && num < tar; i-- {
			res[num] = matrix[bottom][i]
			num++
		}
		bottom--
		for i := bottom; i >= top && num < tar; i-- {
			res[num] = matrix[i][left]
			num++
		}
		left++
	}
	return res
}

func spiralArray(array [][]int) []int {
	if len(array) == 0 || len(array[0]) == 0 {
		return []int{}
	}
	m, n := len(array), len(array[0])
	res := make([]int, m*n)
	left, right := 0, n-1
	top, bottom := 0, m-1
	num := 0
	for num < m*n {
		for i := left; i <= right && num < m*n; i++ {
			res[num] = array[top][i]
			num++
		}
		top++
		for i := top; i <= bottom && num < m*n; i++ {
			res[num] = array[i][right]
			num++
		}
		right--
		for i := right; i >= left && num < m*n; i-- {
			res[num] = array[bottom][i]
			num++
		}
		bottom--
		for i := bottom; i >= top && num < m*n; i-- {
			res[num] = array[i][left]
			num++
		}
		left++
	}
	return res

}
