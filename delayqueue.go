package timingwheel

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

// The start of PriorityQueue implementation.
// Borrowed from https://github.com/nsqio/nsq/blob/master/internal/pqueue/pqueue.go

type Item struct {
	Value    interface{}
	Priority int64
	Index    int
}

// this is a priority queue as implemented by a min heap
// ie. the 0th element is the *lowest* value
type PriorityQueue []*Item

func NewPriorityQueue(capacity int) PriorityQueue {
	return make(PriorityQueue, 0, capacity)
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c {
		npq := make(PriorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	item := x.(*Item)
	item.Index = n
	(*pq)[n] = item
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)
	if n < (c/2) && c > 25 {
		npq := make(PriorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

func (pq *PriorityQueue) PeekAndShift(max int64) (*Item, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	if item.Priority > max {
		return nil, item.Priority - max
	}
	heap.Remove(pq, 0)

	return item, 0
}

// The end of PriorityQueue implementation.

// DelayQueue is an unbounded blocking queue of *Delayed* elements, in which
// an element can only be taken when its delay has expired. The head of the
// queue is the *Delayed* element whose delay expired furthest in the past.
type DelayQueue struct {
	C chan *Bucket

	mu sync.Mutex
	pq PriorityQueue

	// Similar to the sleeping state of runtime.timers.
	sleeping int32
	wakeupC  chan struct{}

	// Similar to the rescheduling state of runtime.timers.
	rescheduling int32
	readyC       chan struct{}
}

// NewDelayQueue creates an instance of DelayQueue with the specified size.
func NewDelayQueue(size int) *DelayQueue {
	return &DelayQueue{
		C:       make(chan *Bucket),
		pq:      NewPriorityQueue(size),
		wakeupC: make(chan struct{}),
		readyC:  make(chan struct{}),
	}
}

// Offer inserts the bucket into the current queue.
func (dq *DelayQueue) Offer(bucket *Bucket) {
	item := &Item{Value: bucket, Priority: bucket.Expiration()}

	dq.mu.Lock()
	heap.Push(&dq.pq, item)
	dq.mu.Unlock()

	if item.Index == 0 {
		// A new item with the earliest expiration is added.
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
		if atomic.CompareAndSwapInt32(&dq.rescheduling, 1, 0) {
			dq.readyC <- struct{}{}
		}
	}
}

// Poll starts an infinite loop, in which it continually waits for an bucket to
// expire and then send the expired bucket to the timing wheel via the channel C.
func (dq *DelayQueue) Poll(exitC chan struct{}) {
	for {
		now := timeToMs(time.Now())

		dq.mu.Lock()
		item, delta := dq.pq.PeekAndShift(now)
		dq.mu.Unlock()

		if item == nil {
			if delta == 0 {
				// No items left.
				atomic.StoreInt32(&dq.rescheduling, 1)
				// Wait until a new item is added.
				select {
				case <-dq.readyC:
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 {
				// At least one item is pending.
				atomic.StoreInt32(&dq.sleeping, 1)
				select {
				case <-dq.wakeupC:
					// A new item with an "earlier" expiration than the current "earliest" one is added.
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// The current "earliest" item expires.

					// Reset the sleeping state since there's no need to receive from wakeupC.
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						// A caller of Offer() is being blocked on sending to wakeupC,
						// drain wakeupC to unblock the caller.
						<-dq.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		}

		bucket := item.Value.(*Bucket)
		select {
		case dq.C <- bucket:
			// Send the expired bucket to the timing wheel.
		case <-exitC:
			goto exit
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
	atomic.StoreInt32(&dq.rescheduling, 0)
}
