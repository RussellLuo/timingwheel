package timingwheel

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type TimingWheel struct {
	tick      int64 // in milliseconds
	wheelSize int64

	interval    int64 // in milliseconds
	currentTime int64 // in milliseconds
	buckets     []*Bucket
	queue       *DelayQueue

	// The parent overflow wheel.
	//
	// NOTE: This field may be updated and read concurrently, through Add().
	overflowWheel unsafe.Pointer // type: *TimingWheel

	workerPoolSize int64

	exitC     chan struct{}
	waitGroup WaitGroupWrapper
}

func NewTimingWheel(tick time.Duration, wheelSize int64) *TimingWheel {
	tickMs := int64(tick / time.Millisecond)
	if tickMs <= 0 {
		panic("tick must be greater than or equal to 1ms")
	}

	startMs := timeToMs(time.Now())

	return newTimingWheel(
		tickMs,
		wheelSize,
		startMs,
		NewDelayQueue(int(wheelSize)),
	)
}

func newTimingWheel(tickMs int64, wheelSize int64, startMs int64, queue *DelayQueue) *TimingWheel {
	buckets := make([]*Bucket, wheelSize)
	for i := range buckets {
		buckets[i] = NewBucket()
	}
	return &TimingWheel{
		tick:        tickMs,
		wheelSize:   wheelSize,
		currentTime: truncate(startMs, tickMs),
		interval:    tickMs * wheelSize,
		buckets:     buckets,
		queue:       queue,
		exitC:       make(chan struct{}),
	}
}

func (tw *TimingWheel) add(t *Timer) bool {
	if t.Expiration < tw.currentTime+tw.tick {
		// Already expired
		return false
	} else if t.Expiration < tw.currentTime+tw.interval {
		virtualID := t.Expiration / tw.tick
		bucket := tw.buckets[virtualID%tw.wheelSize]
		bucket.Add(t)

		if bucket.SetExpiration(virtualID * tw.tick) {
			tw.queue.Offer(bucket)
		}

		return true
	} else {
		// Out of the interval. Put it into the overflow wheel
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel == nil {
			atomic.CompareAndSwapPointer(
				&tw.overflowWheel,
				nil,
				unsafe.Pointer(newTimingWheel(
					tw.interval,
					tw.wheelSize,
					tw.currentTime,
					tw.queue,
				)),
			)
			overflowWheel = atomic.LoadPointer(&tw.overflowWheel)
		}
		return (*TimingWheel)(overflowWheel).add(t)
	}
}

func (tw *TimingWheel) Add(t *Timer) {
	if !tw.add(t) {
		// Already expired
		// TODO: Execute the timer task in a fixed-sized goroutine pool
		tw.waitGroup.Wrap(func() {
			t.Task()
		})
	}
}

func (tw *TimingWheel) advanceClock(expiration int64) {
	if expiration >= tw.currentTime+tw.tick {
		tw.currentTime = truncate(expiration, tw.tick)

		// Try to advance the clock of the overflow wheel if present
		overflowWheel := atomic.LoadPointer(&tw.overflowWheel)
		if overflowWheel != nil {
			(*TimingWheel)(overflowWheel).advanceClock(tw.currentTime)
		}
	}
}

func (tw *TimingWheel) Start() {
	tw.waitGroup.Wrap(func() {
		tw.queue.Poll(tw.exitC)
	})

	tw.waitGroup.Wrap(func() {
		for {
			select {
			case bucket := <-tw.queue.C:
				tw.advanceClock(bucket.Expiration())
				bucket.Flush(tw.Add)
			case <-tw.exitC:
				return
			}
		}
	})
}

func (tw *TimingWheel) Stop() {
	close(tw.exitC)
	tw.waitGroup.Wait()
}
