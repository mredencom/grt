package rwmutex_test

import (
	"grt/q/container/garray"
	"grt/q/internal/rwmutex"
	"gf/g/test/gtest"
	"testing"
	"time"
)

func TestRwmutexIsSafe(t *testing.T) {
	gtest.Case(t, func() {
		lock := rwmutex.New()
		gtest.Assert(lock.IsSafe(), true)

		lock = rwmutex.New(false)
		gtest.Assert(lock.IsSafe(), true)

		lock = rwmutex.New(false, false)
		gtest.Assert(lock.IsSafe(), true)

		lock = rwmutex.New(true, false)
		gtest.Assert(lock.IsSafe(), false)

		lock = rwmutex.New(true, true)
		gtest.Assert(lock.IsSafe(), false)

		lock = rwmutex.New(true)
		gtest.Assert(lock.IsSafe(), false)
	})
}

func TestSafeRwmutex(t *testing.T) {
	gtest.Case(t, func() {
		safeLock := rwmutex.New()
		array := garray.New()

		go func() {
			safeLock.Lock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			safeLock.Unlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			safeLock.Lock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			safeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(80 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
	})
}

func TestSafeReaderRwmutex(t *testing.T) {
	gtest.Case(t, func() {
		safeLock := rwmutex.New()
		array := garray.New()

		go func() {
			safeLock.RLock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			safeLock.RUnlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			safeLock.RLock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			safeLock.RUnlock()
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			safeLock.Lock()
			array.Append(1)
			safeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 6)
	})
}

func TestUnsafeRwmutex(t *testing.T) {
	gtest.Case(t, func() {
		unsafeLock := rwmutex.New(true)
		array := garray.New()

		go func() {
			unsafeLock.Lock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			unsafeLock.Unlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			unsafeLock.Lock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			unsafeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
	})
}
