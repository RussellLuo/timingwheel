// simple example  to use timingwheel
package main

import (
	"fmt"
	"time"

	"github.com/RussellLuo/timingwheel"
)

const (
	WheelSize          = 20              // timingwheel's size ( DelayQueue Size )
	TimingWhellDelay   = 22000           // example task delay  setting, in time.Millisecond
	InternalTimerDelay = 20              // go internal timer delay setting
	MainGoroutineWait  = 25              //
	MaxTaskChannelSize = 100             //
	FixedTimeZone      = "Asia/Shanghai" // 中国上海, 当地时区 China TimeZone 
)

var (
	// global instance of TimingWheel
	tw *timingwheel.TimingWheel
	// global channel for time wheel signal
	taskChannel chan time.Time
	// time location 设置时区
	tl *time.Location
)

// initial the global TimingWheel instance
func init() {
	tw = timingwheel.NewTimingWheel(time.Millisecond, WheelSize)
	taskChannel = make(chan time.Time, MaxTaskChannelSize)
	tl, _ = time.LoadLocation(FixedTimeZone)
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
		fmt.Println(time.Now().In(tl).Format("2006-01-02 15:04:05.999999999 -0700 MST"))
	}()

}

// go internal timer for delay 20 second
func goInternalTimer() {
	internalTimerChannel := time.NewTicker(InternalTimerDelay * time.Second)
	defer internalTimerChannel.Stop()

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
	t := timingwheel.AfterFunc(time.Millisecond*TimingWhellDelay, timingwheelTaskTriggerFunc())
	tw.Add(t)

	// running task goroutine that wait for signal
	go timingWheelTaskRunner()

	// go internal timer
	go goInternalTimer()

	//  for testing only
	for i := 1; i <= MainGoroutineWait; i++ {

		time.Sleep(time.Millisecond * 999)
		fmt.Printf("%d ", i)
		fmt.Println(time.Now().In(tl).Format("2006-01-02 15:04:05.999999999 -0700 MST"))
	}

	// don't quit main goroutine
	// select {}

	// wait for task goroutine finish
	time.Sleep(MainGoroutineWait * time.Second)
}
