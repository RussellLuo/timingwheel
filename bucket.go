package timingwheel

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Timer represents a single event. When the Timer expires, the given
// task will be executed.
type Timer struct {
	Expiration int64 // in milliseconds
	Task       func()

	// The list to which this timer's element belongs.
	//
	// NOTE: This field may be updated and read concurrently,
	// through Timer.Stop() and Bucket.Flush().
	bucket unsafe.Pointer // type: *Bucket

	// The timer's element.
	element *list.Element
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
func AfterFunc(d time.Duration, f func()) *Timer {
	return &Timer{
		Expiration: timeToMs(time.Now().Add(d)),
		Task:       f,
	}
}

func (t *Timer) getBucket() *Bucket {
	return (*Bucket)(atomic.LoadPointer(&t.bucket))
}

func (t *Timer) setBucket(b *Bucket) {
	atomic.StorePointer(&t.bucket, unsafe.Pointer(b))
}

// Stop prevents the Timer from firing. It returns true if the call
// stops the timer, false if the timer has already expired or been stopped.
//
// If the timer has already expired and the function Task has been started in its
// own goroutine; Stop does not wait for Task to complete before returning. If the caller
// needs to know whether Task is completed, it must coordinate with Task explicitly.
func (t *Timer) Stop() bool {
	stopped := false
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		// If b.Remove is called just after the timing wheel's goroutine has:
		//     1. removed t from b (through b.Flush -> b.remove)
		//     2. moved t from b to another bucket ab (through b.Flush -> b.remove and ab.Add)
		// this may fail to remove t due to the change of t's bucket.
		stopped = b.Remove(t)

		// Thus, here we re-get t's possibly new bucket (nil for case 1, or ab (non-nil) for case 2),
		// and retry until the bucket becomes nil, which indicates that t has finally been removed.
	}
	return stopped
}

type Bucket struct {
	mu     sync.Mutex
	timers *list.List

	expiration int64
}

func NewBucket() *Bucket {
	return &Bucket{
		timers:     list.New(),
		expiration: -1,
	}
}

func (b *Bucket) Expiration() int64 {
	return atomic.LoadInt64(&b.expiration)
}

func (b *Bucket) SetExpiration(expiration int64) bool {
	return atomic.SwapInt64(&b.expiration, expiration) != expiration
}

func (b *Bucket) Add(t *Timer) {
	b.mu.Lock()

	e := b.timers.PushBack(t)
	t.setBucket(b)
	t.element = e

	b.mu.Unlock()
}

func (b *Bucket) remove(t *Timer) bool {
	if t.getBucket() != b {
		// If remove is called from t.Stop, and this happens just after the timing wheel's goroutine has:
		//     1. removed t from b (through b.Flush -> b.remove)
		//     2. moved t from b to another bucket ab (through b.Flush -> b.remove and ab.Add)
		// then t.getBucket will return nil for case 1, or ab (non-nil) for case 2.
		// In either case, the returned value does not equal to b.
		return false
	}
	b.timers.Remove(t.element)
	t.setBucket(nil)
	t.element = nil
	return true
}

func (b *Bucket) Remove(t *Timer) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.remove(t)
}

func (b *Bucket) Flush(reinsert func(*Timer)) {
	b.mu.Lock()
	e := b.timers.Front()
	for e != nil {
		next := e.Next()
		t := e.Value.(*Timer)
		b.remove(t)
		reinsert(t)
		e = next
	}
	b.mu.Unlock()

	b.SetExpiration(-1)
}
