package timingwheel

import (
	"testing"
	"time"
)

func TestTimingWheel_AfterFunc(t *testing.T) {
	tw := NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			exitC := make(chan time.Time)

			start := time.Now()
			tw.AfterFunc(d, func() {
				exitC <- time.Now()
			})

			got := (<-exitC).Truncate(time.Millisecond)
			min := start.Add(d).Truncate(time.Millisecond)

			err := 5 * time.Millisecond
			if got.Before(min) || got.After(min.Add(err)) {
				t.Errorf("NewTimer(%s) want [%s, %s], got %s", d, min, min.Add(err), got)
			}
		})
	}
}

func TestTimingWheel_EveryFunc(t *testing.T) {
	tw := NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			ping := make(chan time.Time)

			start := time.Now()
			timer := tw.EveryFunc(d, func() {
				ping <- time.Now()
			})

			for i := 0; i < 10; i++ {
				got := (<-ping).Truncate(time.Millisecond)
				min := start.Add(d).Truncate(time.Millisecond)
				start = time.Now()

				err := 20 * time.Millisecond
				if got.Before(min) || got.After(min.Add(err)) {
					t.Errorf("NewTimer(%s) want [%s, %s], got %s", d, min, min.Add(err), got)
				}
			}
			if timer.Stop() != true {
				t.Fatal()
			}
		})
	}
}
