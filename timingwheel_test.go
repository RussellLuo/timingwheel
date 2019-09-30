package timingwheel_test

import (
	"testing"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func TestTimingWheel_AfterFunc(t *testing.T) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
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

			start := time.Now().UTC()
			tw.AfterFunc(d, func() {
				exitC <- time.Now().UTC()
			})

			got := (<-exitC).Truncate(time.Millisecond)
			min := start.Add(d).Truncate(time.Millisecond)

			err := 5 * time.Millisecond
			if got.Before(min) || got.After(min.Add(err)) {
				t.Errorf("Timer(%s) expiration: want [%s, %s], got %s", d, min, min.Add(err), got)
			}
		})
	}
}

type scheduler struct {
	intervals []time.Duration
	current   int
}

func (s *scheduler) Next(prev time.Time) time.Time {
	if s.current >= len(s.intervals) {
		return time.Time{}
	}
	next := prev.Add(s.intervals[s.current])
	s.current += 1
	return next
}

func TestTimingWheel_ScheduleFunc(t *testing.T) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	s := &scheduler{intervals: []time.Duration{
		1 * time.Millisecond,   // start + 1ms
		4 * time.Millisecond,   // start + 5ms
		5 * time.Millisecond,   // start + 10ms
		40 * time.Millisecond,  // start + 50ms
		50 * time.Millisecond,  // start + 100ms
		400 * time.Millisecond, // start + 500ms
		500 * time.Millisecond, // start + 1s
	}}

	exitC := make(chan time.Time, len(s.intervals))

	start := time.Now().UTC()
	tw.ScheduleFunc(s, func() {
		exitC <- time.Now().UTC()
	})

	accum := time.Duration(0)
	for _, d := range s.intervals {
		got := (<-exitC).Truncate(time.Millisecond)
		accum += d
		min := start.Add(accum).Truncate(time.Millisecond)

		err := 5 * time.Millisecond
		if got.Before(min) || got.After(min.Add(err)) {
			t.Errorf("Timer(%s) expiration: want [%s, %s], got %s", accum, min, min.Add(err), got)
		}
	}
}
