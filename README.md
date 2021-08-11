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
pkg: uuidv7
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkUuidv7RandNoRand-8              3223029               374.1 ns/op           152 B/op 7 allocs/op
BenchmarkUuidv7RandStraight-8             461305              2447 ns/op             152 B/op 7 allocs/op
BenchmarkUuidv7RandSimpleBuf-8           2055200               559.3 ns/op           152 B/op 7 allocs/op
BenchmarkUuidv7RandBufSwitch-8           3105789               399.1 ns/op           152 B/op 7 allocs/op
PASS
ok      uuidv7  6.122s
```

- `BenchmarkUuidv7RandNoRand-8`: Baseline, random bytes are all zeroes.
- `BenchmarkUuidv7RandStraight-8`: Straight read from `/dev/urandom` everytime.
- `BenchmarkUuidv7RandSimpleBuf`: Random bytes are taken from a 32k buffer, refilling it blocks.
- `BenchmarkUuidv7RandBufSwitch-8`: 2x32k buffers. Random bytes are read from the first, when empty,
switch to the other one and refill the first one in the background.
