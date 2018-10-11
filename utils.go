package timingwheel

import (
	"sync"
	"time"
)

// truncate returns the result of rounding x toward zero to a multiple of m.
// If m <= 0, Truncate returns x unchanged.
func truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}

// timeToMs returns an integer number, which represents t in milliseconds.
func timeToMs(t time.Time) int64 {
	return int64(time.Duration(t.UnixNano()) / time.Millisecond)
}

type waitGroupWrapper struct {
	sync.WaitGroup
}

func (w *waitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
