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
BenchmarkUuidv7RandNoRand-8             15591750                75.49 ns/op            8 B/op          1 allocs/op
BenchmarkUuidv7RandStraight-8             536056              2178 ns/op               8 B/op          1 allocs/op
BenchmarkUuidv7RandSimpleBuf-8           4392721               275.6 ns/op             8 B/op          1 allocs/op
BenchmarkUuidv7RandBufSwitch-8           5999540               199.4 ns/op             8 B/op          1 allocs/op
PASS
ok      github.com/brittlesoft/go-uuidv7/pkg/uuidv7     6.349s
```

- `BenchmarkUuidv7RandNoRand-8`: Baseline, random bytes are all zeroes.
- `BenchmarkUuidv7RandStraight-8`: Straight read from `/dev/urandom` everytime.
- `BenchmarkUuidv7RandSimpleBuf`: Random bytes are taken from a 32k buffer, refilling it blocks.
- `BenchmarkUuidv7RandBufSwitch-8`: 2x32k buffers. Random bytes are read from the first, when empty,
switch to the other one and refill the first one in the background.
