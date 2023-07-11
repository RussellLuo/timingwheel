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

func TestTimingWheel_IsStopped(t *testing.T) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	if tw.IsStopped() {
		t.Errorf("IsStopped() = true before stop")
	}
	tw.Stop()
	if !tw.IsStopped() {
		t.Errorf("IsStopped() = false after stop")
	}
	// test stop 2 times
	tw.Stop()
}

func TestTimingWheel_Len(t *testing.T) {
	type fields struct {
		tw  *timingwheel.TimingWheel
		len int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "",
			fields: fields{
				tw:  timingwheel.NewTimingWheel(1*time.Millisecond, 20),
				len: 0,
			},
			want: 0,
		},
		{
			name: "",
			fields: fields{
				tw:  timingwheel.NewTimingWheel(1*time.Second, 20),
				len: 100,
			},
			want: 100,
		},
		{
			name: "",
			fields: fields{
				tw:  timingwheel.NewTimingWheel(1*time.Minute, 20),
				len: 100,
			},
			want: 100,
		},
		{
			name: "",
			fields: fields{
				tw:  timingwheel.NewTimingWheel(1*time.Minute, 20),
				len: 10000,
			},
			want: 10000,
		},
		{
			name: "",
			fields: fields{
				tw:  timingwheel.NewTimingWheel(1*time.Minute, 200),
				len: 100,
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := tt.fields.tw
			tw.Start()
			defer tw.Stop()
			for i := 0; i < tt.fields.len; i++ {
				tw.AfterFunc(time.Duration(i+1)*time.Minute, func() {
				})
			}
			if got := tw.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimingWheel_clear(t *testing.T) {
	tw := timingwheel.NewTimingWheel(1*time.Minute, 20)
	tw.Start()
	l := 10000
	for i := 0; i < l; i++ {
		tw.AfterFunc(time.Duration(i+1)*time.Minute, func() {
		})
	}
	if tw.Len() != l {
		t.Errorf("add events fail")
	}
	tw.Stop()
	if tw.Len() != 0 {
		t.Errorf("clear events fail. tw.Len(): %d", tw.Len())
	}
}
