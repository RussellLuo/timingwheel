package timingwheel_test

import (
	"fmt"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func Example_AddTimer() {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	exitC := make(chan time.Time, 1)
	t := timingwheel.AfterFunc(time.Second, func() {
		fmt.Println("The timer fires")
		exitC <- time.Now()
	})
	tw.Add(t)

	<-exitC

	// Output:
	// The timer fires
}

func Example_StopTimer() {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	t := timingwheel.AfterFunc(time.Second, func() {
		fmt.Println("The timer fires")
	})
	tw.Add(t)

	<-time.After(900 * time.Millisecond)
	t.Stop()

	// Output:
	//
}
