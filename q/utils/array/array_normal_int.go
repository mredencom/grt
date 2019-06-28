package array

import (
	"grt/q/internal/rwmutex"
	"grt/q/utils/random"
	"math"
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

// 返回指定数组的位置数据
// 范围选择并按范围返回项目，如数组[start：end]。
// 注意，如果在并发安全使用中，它返回slice的副本;
// 否则是指向底层数据的指针。
func (ai *IntArray) Range(start, end int) []int {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	length := len(ai.array)
	if start > length || start > end {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	array := ([]int)(nil)
	if ai.mu.IsSafe() {
		ai.mu.RLock()
		defer ai.mu.RUnlock()
		array = make([]int, end-start)
		copy(array, ai.array[start:end])
	} else {
		array = ai.array[start:end]
	}
	return array
}

//将值插入数组的右边。
func (ai *IntArray) Append(value ...int) *IntArray{
	ai.mu.Lock()
	ai.array = append(ai.array,value...)
	ai.mu.Unlock()
	return ai
}

//获取数组的长度
func (ai *IntArray) Len() int {
	ai.mu.Lock()
	length := len(ai.array)
	ai.mu.Unlock()
	return length
}

// Slice返回数组的基础数据。
// 注意，如果在并发安全使用中，它返回slice的副本;
// 否则是指向底层数据的指针。
func (ai *IntArray) Slice() []int {
	array:=([]int)(nil)
	if ai.mu.IsSafe() {
		//这里如果是安全的需要使用copy 只允许读取
		ai.mu.RLock()
		defer ai.mu.RUnlock()
		array = make([]int,len(ai.array))
		//使用copy
		copy(array, ai.array)
	}else {
		array = ai.array
	}
	return array
}

//克隆当前的数组，并返回克隆后的数组
func (ai *IntArray) Clone() []int {
	ai.mu.RLock()
	array:=make([]int, len(ai.array))
	copy(array, ai.array)
	ai.mu.RUnlock()
	return array
}

//清空数组中的元素 并返回数组
func (ai *IntArray) Clear() *IntArray{
	ai.mu.Lock()
	if len(ai.array) > 0{
		ai.array = make([]int, 0)
	}
	ai.mu.Unlock()
	return ai
}

//查询指定的值 然后返回索引值 如果是空数组返回-1
func (ai *IntArray) Search(value int) int {
	//初始化空索引
	result :=-1
	if len(ai.array)==0{
		return result
	}
	ai.mu.RLock()
	for index, v :=range ai.array{
		if value==v {
			result =index
			break
		}
	}
	ai.mu.RUnlock()
	return result
}

//包含检查数组中是否存在值。返回bool
func (ai *IntArray) Contains(value int) bool {
	return ai.Search(value) != -1
}

// 去除重复项，返回去除重复项的数组
func (ai *IntArray) Unique() *IntArray {
	ai.mu.Lock()
	for i := 0; i < len(ai.array)-1; i++ {
		for j := i + 1; j < len(ai.array); j++ {
			if ai.array[i] == ai.array[j] {
				//跳掉重复项的数据
				ai.array = append(ai.array[:j], ai.array[j+1:]...)
			}
		}
	}
	ai.mu.Unlock()
	return ai
}

// LockFunc通过回调函数<f>锁定写入。
func (ai *IntArray) LockFunc(f func(array []int)) *IntArray{
	ai.mu.Lock()
	defer ai.mu.Unlock()
	f(ai.array)
	return ai
}

// RLockFunc通过回调函数<f>锁定读取。
func (ai *IntArray) RLockFunc(f func(array []int)) *IntArray {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	f(ai.array)
	return ai
}

// Merge将<array>合并到当前数组中。
// 参数<array>可以是任何array或slice类型。
// Merge和Append之间的区别是Append仅支持指定的切片类型，
// 但Merge支持更多参数类型。
func (ai *IntArray) Merge(array interface{}) {
	//@todo 待实现
	//switch array.(type) {
	//default:
	//	ai.Append(array)
	//
	//}
}

// 从<startIndex>参数开始填充一个数组，其中包含值<value>，
// 键的num个条目。
func (ai *IntArray) Fill(startIndex, num, value int) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	if startIndex < 0 {
		startIndex = 0
	}
	for i := startIndex; i < startIndex+num; i++ {
		if i > len(ai.array)-1 {
			ai.array = append(ai.array, value)
		} else {
			ai.array[i] = value
		}
	}
	return ai
}
// Chunk将数组拆分为多个数组
// 每个数组的大小由<size>决定。
// 最后一个块可能包含少于size的元素。
func (ai *IntArray) Chunk(size int) [][]int  {
	if size < 1 {
		return nil
	}
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	//取数组的长度
	length := len(ai.array)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var result [][]int
	//循环没有个块数组
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		result = append(result, ai.array[i*size:end])
		i++
	}
	return result
}

// 使用<value>将pad pad数组填充到指定的长度。
// 如果大小为正，则数组在右侧填充，或在左侧填充。
// 如果<size>的绝对值小于或等于数组的长度
// 然后没有填充。
func (ai *IntArray) Pad(size int, value int) *IntArray {
	ai.mu.Lock()
	defer ai.mu.Unlock()
	if size == 0 || (size > 0 && size < len(ai.array)) || (size < 0 && size > -len(ai.array)) {
		return ai
	}
	//if size < 0 {
	//	n = size
	//}
	n := len(ai.array)
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		tmp[i] = value
	}
	if size > 0 {
		ai.array = append(ai.array, tmp...)
	} else {
		ai.array = append(tmp, ai.array...)
	}
	return ai
}