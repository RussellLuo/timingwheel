# timingwheel

Golang implementation of Hierarchical Timing Wheels.


## Installation

```bash
$ go get -u github.com/RussellLuo/timingwheel
```


## Design

`timingwheel` is ported from Kafka's [purgatory][1], which is designed based on [Hierarchical Timing Wheels][2].

中文博客：[层级时间轮的 Golang 实现][3]。


## Documentation

For usage and examples see the [Godoc][4].


## Benchmark

```
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/RussellLuo/timingwheel
BenchmarkTimingWheel_10kTimers_StartStop-8               5000000               411 ns/op
BenchmarkStandardTimer_10kTimers_StartStop-8             5000000               581 ns/op
BenchmarkTimingWheel_100kTimers_StartStop-8              5000000               333 ns/op
BenchmarkStandardTimer_100kTimers_StartStop-8            5000000               541 ns/op
BenchmarkTimingWheel_1mTimers_StartStop-8                5000000               372 ns/op
BenchmarkStandardTimer_1mTimers_StartStop-8              2000000              1302 ns/op
BenchmarkTimingWheel_10mTimers_StartStop-8               5000000               427 ns/op
BenchmarkStandardTimer_10mTimers_StartStop-8             1000000              2180 ns/op
PASS
ok      github.com/RussellLuo/timingwheel       106.640s
```


## License

[MIT][5]


[1]: https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/
[2]: http://www.cs.columbia.edu/~nahum/w6998/papers/ton97-timing-wheels.pdf
[3]: http://russellluo.com/2018/10/golang-implementation-of-hierarchical-timing-wheels.html
[4]: https://godoc.org/github.com/RussellLuo/timingwheel
[5]: http://opensource.org/licenses/MIT
