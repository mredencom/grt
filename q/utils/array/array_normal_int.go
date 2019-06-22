package array

import "grt/q/internal/rwmutex"
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
func (i *IntArray) Get(index int) int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	value := i.array[index]
	return value
}

// 增加指定index的value
func (i *IntArray) Set(index int, value int) *IntArray {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.array[index] = value
	return i
}

// SetArray使用给定的<array>设置基础切片数组。
func (i *IntArray) SetArray(array []int) *IntArray {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.array = array
	return i
}

//Replace从数组的开头用给定的<array>替换数组项
func (i *IntArray) Replace(array []int) *IntArray {
	i.mu.Lock()
	defer i.mu.Unlock()
	max := len(array)
	if max > len(i.array) {
		max = len(i.array)
	}
	for j := 0; j < max; j++ {
		i.array[j] = array[j]
	}
	return i
}