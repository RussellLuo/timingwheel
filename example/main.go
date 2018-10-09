// simple example  to use timingwheel
package main

import (
	"fmt"
	"time"

	"github.com/RussellLuo/timingwheel"
)

const (
	WheelSize          = 20 // timingwheel's size ( DelayQueue Size )
	TimingWhellDelay   = 22 // example task delay  setting
	InternalTimerDelay = 20 // go internal timer delay setting
	MainGoroutineWait  = 25 //
)

var (
	// global instance of TimingWheel
	tw *timingwheel.TimingWheel
	// global channel for time wheel signal
	taskChannel chan time.Time
)

// initial the global TimingWheel instance
func init() {
	tw = timingwheel.NewTimingWheel(time.Millisecond, WheelSize)
	taskChannel = make(chan time.Time)
}

// timingwheelTaskTriggerFunc   a time task trigger to send singal to channel when task's in time
func timingwheelTaskTriggerFunc() func() {
	return func() {
		// record log or something...
		fmt.Println("---------------> The time wheel's timer fires (delay 22  second, send Signal to channel")
		// send signal to channel when time's up
		taskChannel <- time.Now()
	}
}

// timingWheelTaskRunner  a task goroutine to run real task when receive the time's up signal
func timingWheelTaskRunner() {
	// wait for time's up signal from time wheel channel
	<-taskChannel
	// do the real task
	fmt.Println("---------------> The time wheel's REAL TASK,  DO real task")
}

// use timer
func goInternalTimerTaskRunner() {
	go func() {
		fmt.Println("******************> The go internal timer fires ( delay 20 second)")
		fmt.Println(time.Now().String())
	}()

}

// go internal timer for delay 20 second
func goInternalTimer() {
	internalTimerChannel := time.NewTicker(InternalTimerDelay * time.Second)
	select {
	case <-internalTimerChannel.C:
		goInternalTimerTaskRunner()
	}
}

func main() {
	// time wheel start
	tw.Start()
	defer tw.Stop()

	// add time trigger for a new task to delay 22 second
	t := timingwheel.AfterFunc(time.Second*TimingWhellDelay, timingwheelTaskTriggerFunc())
	tw.Add(t)

	// running task goroutine that wait for signal
	go timingWheelTaskRunner()

	// go internal timer
	go goInternalTimer()

	//  for testing only
	for i := 1; i <= MainGoroutineWait; i++ {
		time.Sleep(time.Second * 1)
		fmt.Printf("%d ", i)
		fmt.Println(time.Now().String())
	}

	// don't quit main goroutine
	// select {}

	// wait for task goroutine finish
	time.Sleep(MainGoroutineWait * time.Second)
}
