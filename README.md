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
$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/RussellLuo/timingwheel
BenchmarkTimingWheel_StartStop/N-1m-8            5000000               329 ns/op              83 B/op          2 allocs/op
BenchmarkTimingWheel_StartStop/N-5m-8            5000000               363 ns/op              95 B/op          2 allocs/op
BenchmarkTimingWheel_StartStop/N-10m-8           5000000               440 ns/op              37 B/op          1 allocs/op
BenchmarkStandardTimer_StartStop/N-1m-8         10000000               199 ns/op              64 B/op          1 allocs/op
BenchmarkStandardTimer_StartStop/N-5m-8          2000000               644 ns/op              64 B/op          1 allocs/op
BenchmarkStandardTimer_StartStop/N-10m-8          500000              2434 ns/op              64 B/op          1 allocs/op
PASS
ok      github.com/RussellLuo/timingwheel       116.977s
```


## License

[MIT][5]


[1]: https://www.confluent.io/blog/apache-kafka-purgatory-hierarchical-timing-wheels/
[2]: http://www.cs.columbia.edu/~nahum/w6998/papers/ton97-timing-wheels.pdf
[3]: http://russellluo.com/2018/10/golang-implementation-of-hierarchical-timing-wheels.html
[4]: https://godoc.org/github.com/RussellLuo/timingwheel
[5]: http://opensource.org/licenses/MIT
