# timingwheel example

example to use timingwhell 



## example code 
```
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


```

## run example OUTPUT
```
1 2018-10-09 20:29:51.4755922 +0000 UTC m=+1.003481501
2 2018-10-09 20:29:52.4777751 +0000 UTC m=+2.005663701
3 2018-10-09 20:29:53.4792432 +0000 UTC m=+3.007133101
4 2018-10-09 20:29:54.4812492 +0000 UTC m=+4.009164701
5 2018-10-09 20:29:55.482701 +0000 UTC m=+5.010588801
6 2018-10-09 20:29:56.4480973 +0000 UTC m=+6.011321301
7 2018-10-09 20:29:57.4488363 +0000 UTC m=+7.012058801
8 2018-10-09 20:29:58.4491704 +0000 UTC m=+8.012393101
9 2018-10-09 20:29:59.4503951 +0000 UTC m=+9.013618801
10 2018-10-09 20:30:00.4515239 +0000 UTC m=+10.014748601
11 2018-10-09 20:30:01.4526743 +0000 UTC m=+11.015896801
12 2018-10-09 20:30:02.454003 +0000 UTC m=+12.017229201
13 2018-10-09 20:30:03.4548702 +0000 UTC m=+13.018095001
14 2018-10-09 20:30:04.455348 +0000 UTC m=+14.018571001
15 2018-10-09 20:30:05.4562363 +0000 UTC m=+15.019460701
16 2018-10-09 20:30:06.456717 +0000 UTC m=+16.019941301
17 2018-10-09 20:30:07.4576972 +0000 UTC m=+17.020921601
18 2018-10-09 20:30:08.4583268 +0000 UTC m=+18.021552201
19 2018-10-09 20:30:09.4592617 +0000 UTC m=+19.022483601
******************> The go internal timer fires ( delay 20 second)
2018-10-09 20:30:10.4402801 +0000 UTC m=+20.003504401
20 2018-10-09 20:30:10.460425 +0000 UTC m=+20.023648101
21 2018-10-09 20:30:11.4609443 +0000 UTC m=+21.024180601
22 2018-10-09 20:30:12.4620006 +0000 UTC m=+22.025225001
---------------> The time wheel's timer fires (delay 22 second, send Signal to channel
--------------->> The time wheel's timer up (delay 22 second), DO real task
23 2018-10-09 20:30:13.4625539 +0000 UTC m=+23.025779001
24 2018-10-09 20:30:14.46396 +0000 UTC m=+24.027184801
25 2018-10-09 20:30:15.4648446 +0000 UTC m=+25.028068601
```