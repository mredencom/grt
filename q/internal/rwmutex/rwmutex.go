// 当前包是提供并发安全控制
package rwmutex

import "sync"

// RWMutex是一个sync.RWMutex，具有并发安全功能的切换.
type RWMutex struct {
	sync.RWMutex
	safe bool
}
// 建立一个安全并发
func New(unsafe ...bool) *RWMutex {
	mu := new(RWMutex)
	if len(unsafe) > 0 {
		mu.safe = !unsafe[0]
	} else {
		mu.safe = true
	}
	return mu
}
// 判断是否安全
func (mu *RWMutex) IsSafe() bool {
	return mu.safe
}
// 给加上锁保证安全
func (mu *RWMutex) Lock() {
	if mu.safe {
		mu.RWMutex.Lock()
	}
}
// 完成后去除掉锁
func (mu *RWMutex) Unlock() {
	if mu.safe {
		mu.RWMutex.Unlock()
	}
}

func (mu *RWMutex) RLock() {
	if mu.safe {
		mu.RWMutex.RLock()
	}
}

func (mu *RWMutex) RUnlock() {
	if mu.safe {
		mu.RWMutex.RUnlock()
	}
}
