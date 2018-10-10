package timingwheel_test

import (
	"testing"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func benchmarkTimingWheel_StartStop(b *testing.B, genD func(int) time.Duration) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	timers := make([]*timingwheel.Timer, b.N)
	for i := 0; i < b.N; i++ {
		timers[i] = tw.AfterFunc(genD(i), func() {})
	}

	for i := 0; i < b.N; i++ {
		timers[i].Stop()
	}
}

func benchmarkStandardTimer_StartStop(b *testing.B, genD func(int) time.Duration) {
	timers := make([]*time.Timer, b.N)
	for i := 0; i < b.N; i++ {
		timers[i] = time.AfterFunc(genD(i), func() {})
	}

	for i := 0; i < b.N; i++ {
		timers[i].Stop()
	}
}

func BenchmarkTimingWheel_StartStop_1millionTimers_WithSameDurations(b *testing.B) {
	b.N = 1000000
	benchmarkTimingWheel_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkStandardTimer_StartStop_1millionTimers_WithSameDurations(b *testing.B) {
	b.N = 1000000
	benchmarkStandardTimer_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkTimingWheel_StartStop_1millionTimers_WithDifferentDurations(b *testing.B) {
	b.N = 1000000
	benchmarkTimingWheel_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkStandardTimer_StartStop_1millionTimers_WithDifferentDurations(b *testing.B) {
	b.N = 1000000
	benchmarkStandardTimer_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkTimingWheel_StartStop_5millionsTimers_WithSameDurations(b *testing.B) {
	b.N = 5000000
	benchmarkTimingWheel_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkStandardTimer_StartStop_5millionsTimers_WithSameDurations(b *testing.B) {
	b.N = 5000000
	benchmarkStandardTimer_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkTimingWheel_StartStop_5millionsTimers_WithDifferentDurations(b *testing.B) {
	b.N = 5000000
	benchmarkTimingWheel_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkStandardTimer_StartStop_5millionsTimers_WithDifferentDurations(b *testing.B) {
	b.N = 5000000
	benchmarkStandardTimer_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkTimingWheel_StartStop_10millionsTimers_WithSameDurations(b *testing.B) {
	b.N = 10000000
	benchmarkTimingWheel_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkStandardTimer_StartStop_10millionsTimers_WithSameDurations(b *testing.B) {
	b.N = 10000000
	benchmarkStandardTimer_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkTimingWheel_StartStop_10millionsTimers_WithDifferentDurations(b *testing.B) {
	b.N = 10000000
	benchmarkTimingWheel_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkStandardTimer_StartStop_10millionsTimers_WithDifferentDurations(b *testing.B) {
	b.N = 10000000
	benchmarkStandardTimer_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}
