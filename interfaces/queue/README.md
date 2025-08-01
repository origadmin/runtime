# Queue

```
Name count (26550976) Average time (202.1 ns/op) Average memory usage (7 B/op) Average allocation count (0 allocs/op)
BenchmarkQueue/LockFreeQueue/Offer-20           26550976               202.1 ns/op             7 B/op          0 allocs/op
```
# BenchmarkQueue/LockFreeQueue/Offer-20
- Name: indicates the name of the current benchmark.
- BenchmarkQueue is the top-level name of the benchmark.
- LockFreeQueue Indicates the queue type.
- Offer indicates the specific test operation (in this case, the team operation).
- -20 Indicates the GOMAXPROCS setting (parallelism) at run time.
# 26550976
- Count: Indicates how many times the benchmark has been run.
This number represents the number of operations that the test can perform in a given time period (the default is 1 second).
# 202.1 ns/op
- Average time: indicates the average time spent on each operation.
- ns stands for nanosecond (1 nanosecond = 10^-9 seconds).
- op Indicates an operation.
- Therefore, 202.1 ns/op means that the average time spent per queuing operation is 202.1 nanoseconds.
# 7 B/op
- Average memory usage: indicates the average number of bytes allocated per operation.
- B Indicates byte.
- op Indicates an operation.
- Therefore, 7 B/op means that an average of 7 bytes of memory are allocated for each queuing operation.
# 0 allocs/op
- Average number of memory allocations: indicates the number of memory allocations per operation.
- allocs Indicates the number of allocation times.
- op Indicates an operation.
- Therefore, 0 allocs/op means that no memory is allocated for each queueing operation.
# Summary
- Count (26550976) : Indicates how many times the test was run.
- Average Time (202.1 ns/op) : indicates the average time spent on each queuing operation.
- Average memory usage (7 B/op) : indicates the average number of bytes allocated for each queuing operation.
- Average number of memory allocations (0 allocs/op) : indicates the number of memory allocations per queued operation.
These metrics can help you understand the performance of different operations to further optimize your code. For example, if you find that the average time of an operation is long or the memory usage is high, you can optimize it accordingly.

# Benchmarking
```
$>go test -v -race -benchmem -cpu="8,16,32" -benchtime=5s -run=^$ -bench ^BenchmarkQueue$ -cpuprofile=cpu -memprofile=mem github.com/origadmin/toolkits/queue
goos: windows
goarch: amd64
pkg: github.com/origadmin/toolkits/queue
cpu: 12th Gen Intel(R) Core(TM) i7-12700H
BenchmarkQueue
BenchmarkQueue
BenchmarkQueue/LockFreeQueue/Offer
BenchmarkQueue/LockFreeQueue/Offer-8            31714936               199.1 ns/op             8 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Offer-16           26100540               205.9 ns/op             8 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Offer-32           31186132               226.5 ns/op             8 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Poll
BenchmarkQueue/LockFreeQueue/Poll-8             23684212               279.5 ns/op             0 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Poll-16            24003744               250.0 ns/op             0 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Poll-32            25795556               260.0 ns/op             0 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Mixed
BenchmarkQueue/LockFreeQueue/Mixed-8             7564692               858.1 ns/op             9 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Mixed-16            6900072               922.3 ns/op            12 B/op          0 allocs/op
BenchmarkQueue/LockFreeQueue/Mixed-32            6511975               956.4 ns/op            10 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Offer
BenchmarkQueue/MutexQueue/Offer-8               13904630               440.2 ns/op            41 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Offer-16              14502745               462.0 ns/op            37 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Offer-32              10905456               643.8 ns/op            77 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Poll
BenchmarkQueue/MutexQueue/Poll-8                21881016               774.1 ns/op             0 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Poll-16               18971287               668.8 ns/op             0 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Poll-32                9894604               874.3 ns/op             0 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Mixed
BenchmarkQueue/MutexQueue/Mixed-8                5919752              1240 ns/op              81 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Mixed-16               5668646              1077 ns/op               0 B/op          0 allocs/op
BenchmarkQueue/MutexQueue/Mixed-32               5773432               935.0 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/origadmin/toolkits/queue     220.216s
```

# Benchmarking
```
$>go test -v -race -cpu=8 -benchmem -run=^$ -benchtime=5s -bench ^BenchmarkQueueMix$ github.com/origadmin/toolkits/queue
goos: windows
goarch: amd64
pkg: github.com/origadmin/toolkits/queue
cpu: 12th Gen Intel(R) Core(TM) i7-12700H
BenchmarkQueueMix
BenchmarkQueueMix/DesignQueue/Mixed
BenchmarkQueueMix/DesignQueue/Mixed-8            4249044              1578 ns/op              20 B/op          1 allocs/op
BenchmarkQueueMix/WrapQueue/Mixed
BenchmarkQueueMix/WrapQueue/Mixed-8              1000000             39670 ns/op            1497 B/op         93 allocs/op
BenchmarkQueueMix/ChanQueue/Mixed
BenchmarkQueueMix/ChanQueue/Mixed-8              8007457               653.5 ns/op             0 B/op          0 allocs/op
BenchmarkQueueMix/LockFreeQueue/Mixed
BenchmarkQueueMix/LockFreeQueue/Mixed-8          5856188              1031 ns/op               9 B/op          0 allocs/op
BenchmarkQueueMix/MutexQueue/Mixed
BenchmarkQueueMix/MutexQueue/Mixed-8             4987773              1212 ns/op              20 B/op          0 allocs/op
PASS
ok      github.com/origadmin/toolkits/queue     69.761s
```

// Running tool: go test -v -race -count=5 -benchmem -run=^$ -bench ^BenchmarkQueue$ -parallel=32 github.com/origadmin/toolkits/queue
// Running tool: go test -v -race -benchmem -cpu="8,16,32" -benchtime=5s -run=^$ -bench ^BenchmarkQueue$ -cpuprofile=cpu -memprofile=mem github.com/origadmin/toolkits/queue
```shell
go test -v -race -benchmem -run=^$ -bench ^BenchmarkQueue$ -benchtime=30s -cpuprofile=cpu.out -memprofile=mem.out -parallel=32 github.com/origadmin/toolkits/queue
go test -v -race -benchmem -cpu="8" -run=^$ -bench ^BenchmarkQueueMix$ -benchtime=30s -cpuprofile=cpu -memprofile=mem github.com/origadmin/toolkits/queue
go test -v -race -benchmem -cpu="8" -run=^$ -count=5 -bench ^BenchmarkQueueMix$ -benchtime=15s -cpuprofile=cpu -memprofile=mem github.com/origadmin/toolkits/queue

go tool pprof --http=:8080 cpu.out
go tool pprof --http=:8081 mem.out
```


