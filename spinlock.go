package timingwheel

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/v2pro/plz/gls"
)

type spinLock struct {
	owner int64
	count int64
	lock  int64
}

func (sl *spinLock) Lock() {
	me := GetGoroutineId()

	if atomic.LoadInt64(&sl.owner) == me { // 如果当前线程已经获取到了锁，线程数增加一，然后返回
		sl.count++
		return
	}
	// 如果没获取到锁，则通过CAS自旋
	for !atomic.CompareAndSwapInt64(&sl.lock, 0, 1) {
		runtime.Gosched()
	}
	atomic.StoreInt64(&sl.owner, me)
}
func (sl *spinLock) Unlock() {
	if atomic.LoadInt64(&sl.owner) != GetGoroutineId() {
		panic("illegalMonitorStateError")
	}

	if sl.count > 0 { // 如果大于0，表示当前线程多次获取了该锁，释放锁通过count减一来模拟
		sl.count--
	} else { // 如果count==0，可以将锁释放，这样就能保证获取锁的次数与释放锁的次数是一致的了。
		sl.owner = -1
		atomic.StoreInt64(&sl.lock, 0)
	}
}

func GetGoroutineId() int64 {
	return gls.GoID()
}

func NewSpinLock() sync.Locker {
	lock := &spinLock{owner: -1}
	return lock
}
