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
BenchmarkTimingWheel_StartStop_1millionTimers_WithSameDurations-8                1000000               314 ns/op
BenchmarkStandardTimer_StartStop_1millionTimers_WithSameDurations-8              1000000               233 ns/op
BenchmarkTimingWheel_StartStop_1millionTimers_WithDifferentDurations-8           1000000               323 ns/op
BenchmarkStandardTimer_StartStop_1millionTimers_WithDifferentDurations-8         1000000               239 ns/op
BenchmarkTimingWheel_StartStop_5millionsTimers_WithSameDurations-8               5000000               608 ns/op
BenchmarkStandardTimer_StartStop_5millionsTimers_WithSameDurations-8             5000000               288 ns/op
BenchmarkTimingWheel_StartStop_5millionsTimers_WithDifferentDurations-8          5000000               330 ns/op
BenchmarkStandardTimer_StartStop_5millionsTimers_WithDifferentDurations-8        5000000               465 ns/op
BenchmarkTimingWheel_StartStop_10millionsTimers_WithSameDurations-8             10000000               502 ns/op
BenchmarkStandardTimer_StartStop_10millionsTimers_WithSameDurations-8           10000000               376 ns/op
BenchmarkTimingWheel_StartStop_10millionsTimers_WithDifferentDurations-8        10000000               344 ns/op
BenchmarkStandardTimer_StartStop_10millionsTimers_WithDifferentDurations-8      10000000              1459 ns/op
```


## License

[MIT][5]


[1]: https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/
[2]: http://www.cs.columbia.edu/~nahum/w6998/papers/ton97-timing-wheels.pdf
[3]: http://russellluo.com/2018/10/golang-implementation-of-hierarchical-timing-wheels.html
[4]: https://godoc.org/github.com/RussellLuo/timingwheel
[5]: http://opensource.org/licenses/MIT
