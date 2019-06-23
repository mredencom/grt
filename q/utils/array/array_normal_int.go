package array

import (
	"grt/q/internal/rwmutex"
	"grt/q/utils/random"
	"sort"
)
// 定义一个数组的结构
type IntArray struct {
	mu    *rwmutex.RWMutex //使用互斥量控制并发
	array []int
}

// NewIntArray 返回一个空数组
// 创建一个数组
// unsafe 默认 为false
func NewIntArray(unsafe ...bool) *IntArray {
	return NewIntArraySize(0, 0, unsafe...)
}

// NewIntArraySize返回一个数组
// unsafe默认false 表示一个并非安全的数组
func NewIntArraySize(size int, cap int, unsafe ...bool) *IntArray {
	return &IntArray{
		mu:    rwmutex.New(unsafe...),
		array: make([]int, size, cap),
	}
}

// NewIntArrayFrom返回一个数组
// unsafe默认false 表示一个并非一个并发安全的数组
func NewIntArrayFrom(array []int, unsafe ...bool) *IntArray {
	return &IntArray{
		mu:    rwmutex.New(unsafe...),
		array: array,
	}
}

// NewIntArrayFromCopy 创建一个数组copy到另外一个数组上 返回一个数组
//参数<unsafe>用于指定是否在非并发安全中使用数组，
//默认为false
func NewIntArrayFromCopy(array []int, unsafe ...bool) *IntArray {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return &IntArray{
		mu:    rwmutex.New(unsafe...),
		array: newArray,
	}
}

// 获取指定索引上的值 并返回
// ⚠️不能数组下标越界调用
func (ai *IntArray) Get(index int) int {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	value := ai.array[index]
	return value
}

// 增加指定index的value
func (ai *IntArray) Set(index int, value int) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	ai.array[index] = value
	return ai
}

// SetArray使用给定的<array>设置基础切片数组。
func (ai *IntArray) SetArray(array []int) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	ai.array = array
	return ai
}

//Replace从数组的开头用给定的<array>替换数组项
func (ai *IntArray) Replace(array []int) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	max := len(array)
	if max > len(ai.array) {
		max = len(ai.array)
	}
	for j := 0; j < max; j++ {
		ai.array[j] = array[j]
	}
	return ai
}

//返回数组中的累计求和。
func (ai *IntArray) Sum() (sum int) {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	for _, v := range ai.array {
		sum += v
	}
	return
}

// Sort按递增顺序对数组进行排序。
// 参数<reverse>控制是否排序
// 按递增顺序（默认）或递减顺序
// 这里使用的快速排序。
func (ai *IntArray) Sort(reverse ...bool) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	if len(reverse) > 0 && reverse[0] {
		sort.Slice(ai.array, func(a, j int) bool {
			if ai.array[a] < ai.array[j] {
				return false
			}
			return true
		})
	} else {
		// 递增排序.
		sort.Ints(ai.array)
	}
	return ai
}

//SortFunc按自定义函数<less>对数组进行排序
func (ai *IntArray) SortFunc(less func(v1,v2 int) bool) *IntArray  {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	sort.Slice(ai.array, func(i, j int) bool {
		return less(ai.array[i], ai.array[j])
	})
	return ai
}

// InsertBefore将<value>插入<index>的前面。
func (ai *IntArray) InsertBefore(index int, value int) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	rear := append([]int{}, ai.array[index:]...)
	ai.array = append(ai.array[0:index], value)
	ai.array = append(ai.array, rear...)
	return ai
}

// InsertAfter将<value>插入<index>的后面。
func (ai *IntArray) InsertAfter(index int, value int) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	rear := append([]int{}, ai.array[index+1:]...)
	ai.array = append(ai.array[0:index+1], value)
	ai.array = append(ai.array, rear...)
	return ai
}

// 删除指定索引上的值 返回删除的值
func (ai *IntArray) Remove(index int) int {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	// 删除时确定数组边界以提高删除效率.
	if index == 0 {
		value := ai.array[0]
		ai.array = ai.array[1:]
		return value
	} else if index == len(ai.array)-1 {
		value := ai.array[index]
		ai.array = ai.array[:index]
		return value
	}
	// 如果是非边界删除，
	// 它将涉及创建一个数组，
	// 然后删除效率较低。
	value := ai.array[index]
	ai.array = append(ai.array[:index], ai.array[index+1:]...)
	return value
}

// PushLeft将一个或多个项目推送到数组的开头。
func (ai *IntArray) PushLeft(value ...int) *IntArray {
	ai.mu.Lock()
	ai.array = append(value, ai.array...)
	ai.mu.Unlock()
	return ai
}

// PushRight将一个或多个项目推送到数组的尾部。
func (ai *IntArray) PushRight(value ...int) *IntArray {
	ai.mu.Lock()
	ai.array = append(ai.array, value...)
	ai.mu.Unlock()
	return ai
}

// PopLeft去除开头元素并返回数组开头的元素。
func (ai *IntArray) PopLeft() int {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	value := ai.array[0]
	ai.array = ai.array[1:]
	return value
}

// PopRight去除尾部元素并返回数组尾部的元素。
func (ai *IntArray) PopRight() int {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	index := len(ai.array) - 1
	value := ai.array[index]
	ai.array = ai.array[:index]
	return value
}

// PopRand随机删除并从数组中返回一个删除的值。
func (ai *IntArray) PopRand() int {
	return ai.Remove(random.Intn(len(ai.array)))
}

//PopRands随机删除并从数组中返回<size>项。例如：size=2 就在数组中随机删除两个元素。
func (ai *IntArray) PopRands(size int) []int {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	if size > len(ai.array) {
		size = len(ai.array)
	}
	array := make([]int, size)
	for i := 0; i < size; i++ {
		index := random.Intn(len(ai.array))
		array[i] = ai.array[index]
		ai.array = append(ai.array[:index], ai.array[index+1:]...)
	}
	return array
}

// PopLefts删除并返回数组开头的<size>项。
func (ai *IntArray) PopLefts(size int) []int {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	length := len(ai.array)
	if size > length {
		size = length
	}
	value := ai.array[0:size]
	//更新元素.
	ai.array = ai.array[size:]
	return value
}

// PopRights删除并返回数组开头的<size>项。
func (ai *IntArray) PopRights(size int) []int {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	index := len(ai.array) - size
	if index < 0 {
		index = 0
	}
	value := ai.array[index:]
	ai.array = ai.array[:index]
	return value
}