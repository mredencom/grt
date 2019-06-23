// Package rand提供高性能的随机字符串生成功能.
package random

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	digits  = []rune("0123456789")
)

// 满足随机计算是否满足给定的概率<num> / <total>。
func Meet(num, total int) bool {
	return Intn(total) < num
}

// }MeetProb随机计算是否满足给定的概率。
func MeetProb(prob float32) bool {
	return Intn(1e7) < int(prob*1e7)
}

// N在min和max之间返回一个随机int  -  [min，max]。
func N(min, max int) int {
	if min >= max {
		return min
	}
	if min >= 0 {
		// 因为In Tn不支持负数，
		// 所以我们应该先将值移到左边，
		// 然后调用Intn产生随机数，
		// 最后将结果移到右边。
		return Intn(max-(min-0)+1) + (min - 0)
	}
	if min < 0 {
		// 因为In Tn不支持负数，
		// 所以我们应该先将值移到右边，
		// 然后调用Intn产生随机数，
		// 最后将结果移到左边。
		return Intn(max+(0-min)+1) - (0 - min)
	}
	return 0
}

// 已弃用
// N的别名
func Rand(min, max int) int {
	return N(min, max)
}

// Str返回包含数字和字母的随机字符串，其长度为<n>。
func Str(n int) string {
	b := make([]rune, n)
	for i := range b {
		if Intn(2) == 1 {
			b[i] = digits[Intn(10)]
		} else {
			b[i] = letters[Intn(52)]
		}
	}
	return string(b)
}

// 已弃用
// Str的别名
func RandStr(n int) string {
	return Str(n)
}

// 数字返回一个只包含数字的随机字符串，其长度为<n>。
func Digits(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = digits[Intn(10)]
	}
	return string(b)

}

//已弃用
// 数字别名
func RandDigits(n int) string {
	return Digits(n)
}

// 字母返回一个只包含字母的随机字符串，其长度为<n>。
func Letters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[Intn(52)]
	}
	return string(b)

}

//已弃用
// 字母别名。
func RandLetters(n int) string {
	return Letters(n)
}

// Perm作为n个int的片段返回整数[0，n]的伪随机置换。
func Perm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}
