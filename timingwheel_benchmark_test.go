package timingwheel_test

import (
	"testing"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func genD(i int) time.Duration {
	return time.Duration(i%10000) * time.Millisecond
}

func benchmarkTimingWheel_StartStop(b *testing.B, n int) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	for i := 0; i < n; i++ {
		tw.AfterFunc(genD(i), func() {})
	}
	b.ResetTimer()

	timers := make([]*timingwheel.Timer, b.N)
	for i := 0; i < b.N; i++ {
		timers[i] = tw.AfterFunc(genD(i), func() {})
	}

	for i := 0; i < b.N; i++ {
		timers[i].Stop()
	}
}

func benchmarkStandardTimer_StartStop(b *testing.B, n int) {
	for i := 0; i < n; i++ {
		time.AfterFunc(genD(i), func() {})
	}
	b.ResetTimer()

	timers := make([]*time.Timer, b.N)
	for i := 0; i < b.N; i++ {
		timers[i] = time.AfterFunc(genD(i), func() {})
	}

	for i := 0; i < b.N; i++ {
		timers[i].Stop()
	}
}

func BenchmarkTimingWheel_10kTimers_StartStop(b *testing.B) {
	benchmarkTimingWheel_StartStop(b, 10000)
}

func BenchmarkStandardTimer_10kTimers_StartStop(b *testing.B) {
	benchmarkStandardTimer_StartStop(b, 10000)
}

func BenchmarkTimingWheel_100kTimers_StartStop(b *testing.B) {
	benchmarkTimingWheel_StartStop(b, 100000)
}

func BenchmarkStandardTimer_100kTimers_StartStop(b *testing.B) {
	benchmarkStandardTimer_StartStop(b, 100000)
}

func BenchmarkTimingWheel_1mTimers_StartStop(b *testing.B) {
	benchmarkTimingWheel_StartStop(b, 1000000)
}

func BenchmarkStandardTimer_1mTimers_StartStop(b *testing.B) {
	benchmarkStandardTimer_StartStop(b, 1000000)
}

func BenchmarkTimingWheel_10mTimers_StartStop(b *testing.B) {
	benchmarkTimingWheel_StartStop(b, 10000000)
}

func BenchmarkStandardTimer_10mTimers_StartStop(b *testing.B) {
	benchmarkStandardTimer_StartStop(b, 10000000)
}
