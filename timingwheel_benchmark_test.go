package timingwheel_test

import (
	"testing"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func benchmarkTimingWheel_AddStop(b *testing.B, genD func(int) time.Duration) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	timers := make([]*timingwheel.Timer, b.N)
	for i := 0; i < b.N; i++ {
		timers[i] = timingwheel.AfterFunc(genD(i), func() {})
		tw.Add(timers[i])
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

func BenchmarkTimingWheel_AddStop_SameDurations_1million(b *testing.B) {
	b.N = 1000000
	benchmarkTimingWheel_AddStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkStandardTimer_StartStop_SameDurations_1million(b *testing.B) {
	b.N = 1000000
	benchmarkStandardTimer_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkTimingWheel_AddStop_DifferentDurations_1million(b *testing.B) {
	b.N = 1000000
	benchmarkTimingWheel_AddStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkStandardTimer_StartStop_DifferentDurations_1million(b *testing.B) {
	b.N = 1000000
	benchmarkStandardTimer_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkTimingWheel_AddStop_SameDurations_5million(b *testing.B) {
	b.N = 5000000
	benchmarkTimingWheel_AddStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkStandardTimer_StartStop_SameDurations_5million(b *testing.B) {
	b.N = 5000000
	benchmarkStandardTimer_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkTimingWheel_AddStop_DifferentDurations_5million(b *testing.B) {
	b.N = 5000000
	benchmarkTimingWheel_AddStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkStandardTimer_StartStop_DifferentDurations_5million(b *testing.B) {
	b.N = 5000000
	benchmarkStandardTimer_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkTimingWheel_AddStop_SameDurations_10million(b *testing.B) {
	b.N = 10000000
	benchmarkTimingWheel_AddStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkStandardTimer_StartStop_SameDurations_10million(b *testing.B) {
	b.N = 10000000
	benchmarkStandardTimer_StartStop(b, func(int) time.Duration {
		return time.Second
	})
}

func BenchmarkTimingWheel_AddStop_DifferentDurations_10million(b *testing.B) {
	b.N = 10000000
	benchmarkTimingWheel_AddStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}

func BenchmarkStandardTimer_StartStop_DifferentDurations_10million(b *testing.B) {
	b.N = 10000000
	benchmarkStandardTimer_StartStop(b, func(i int) time.Duration {
		return time.Duration(i%10000) * time.Millisecond
	})
}
