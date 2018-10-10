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
$ go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/RussellLuo/timingwheel
BenchmarkTimingWheel_AddStop_SameDurations_1million-8                    1000000               425 ns/op
BenchmarkStandardTimer_StartStop_SameDurations_1million-8                1000000               258 ns/op
BenchmarkTimingWheel_AddStop_DifferentDurations_1million-8               1000000               337 ns/op
BenchmarkStandardTimer_StartStop_DifferentDurations_1million-8           1000000               237 ns/op
BenchmarkTimingWheel_AddStop_SameDurations_5million-8                    5000000               532 ns/op
BenchmarkStandardTimer_StartStop_SameDurations_5million-8                5000000               301 ns/op
BenchmarkTimingWheel_AddStop_DifferentDurations_5million-8               5000000               325 ns/op
BenchmarkStandardTimer_StartStop_DifferentDurations_5million-8           5000000               470 ns/op
BenchmarkTimingWheel_AddStop_SameDurations_10million-8                  10000000               460 ns/op
BenchmarkStandardTimer_StartStop_SameDurations_10million-8              10000000               353 ns/op
BenchmarkTimingWheel_AddStop_DifferentDurations_10million-8             10000000               333 ns/op
BenchmarkStandardTimer_StartStop_DifferentDurations_10million-8         10000000              1462 ns/op
```


## License

[MIT][5]


[1]: https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/
[2]: http://www.cs.columbia.edu/~nahum/w6998/papers/ton97-timing-wheels.pdf
[3]: http://russellluo.com/2018/10/golang-implementation-of-hierarchical-timing-wheels.html
[4]: https://godoc.org/github.com/RussellLuo/timingwheel
[5]: http://opensource.org/licenses/MIT
