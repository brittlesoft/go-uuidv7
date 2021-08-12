# go-uuidv7

Exercise / napkin quality code that implements the
[UUIDv7](https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html)
Internet-Draft.

Mostly as a brain tease.
Some experiments have been done w.r.t performance since getting random bytes was the slowest part of
the program in the first version.


## Performance

```
goos: linux
goarch: amd64
pkg: github.com/brittlesoft/go-uuidv7/pkg/uuidv7
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkUuidv7RandNoRand-8             15836040                76.84 ns/op            8 B/op          1 allocs/op
BenchmarkUuidv7RandStraight-8             467644              2212 ns/op               8 B/op          1 allocs/op
BenchmarkUuidv7RandSimpleBuf-8           4506566               269.8 ns/op             8 B/op          1 allocs/op
BenchmarkUuidv7RandBufSwitch-8           5953366               198.6 ns/op             8 B/op          1 allocs/op
BenchmarkUuidv7RandUnpredictable-8      12598537                93.30 ns/op            8 B/op          1 allocs/op
BenchmarkUnpredictableUuidv4-8          19120903                66.36 ns/op           16 B/op          1 allocs/op
PASS
ok      github.com/brittlesoft/go-uuidv7/pkg/uuidv7     8.856s
```

- `BenchmarkUuidv7RandNoRand-8`: Baseline, random bytes are all zeroes.
- `BenchmarkUuidv7RandStraight-8`: Straight read from `/dev/urandom` everytime.
- `BenchmarkUuidv7RandSimpleBuf`: Random bytes are taken from a 32k buffer, refilling it blocks.
- `BenchmarkUuidv7RandBufSwitch-8`: 2x32k buffers. Random bytes are read from the first, when empty,
- `BenchmarkUuidv7RandUnpredictable-8`: Random bytes are sourced from [unpredictable](https://github.com/art4711/unpredictable).
  (See [arc4random(3)](https://man.openbsd.org/arc4random))
- `BenchmarkUnpredictableUuidv4-8`: Fully random uuidv4 sourced from unpredictable.
