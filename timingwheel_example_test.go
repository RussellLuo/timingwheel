package timingwheel_test

import (
	"fmt"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func Example() {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	exitC := make(chan time.Time)
	t := timingwheel.AfterFunc(time.Second, func() {
		exitC <- time.Now()
	})
	tw.Add(t)

	<-exitC
	fmt.Println("The timer fires")

	// Output:
	// The timer fires
}
