package uuidv7

import (
	"log"
	"testing"
)

var u Uuidv7

func benchmarkUuidv7(rs randSource, b *testing.B) {
	us := Uuidv7Source{}
	us.rs = rs
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		patate, err := us.New()
		if err != nil {
			log.Fatal(err)
		}
		u = patate
	}
}

func BenchmarkUuidv7RandNoRand(b *testing.B) {
	benchmarkUuidv7(&randSourceNoRand{}, b)
}

func BenchmarkUuidv7RandStraight(b *testing.B) {
	benchmarkUuidv7(&randSourceStraight{}, b)
}

func BenchmarkUuidv7RandSimpleBuf(b *testing.B) {
	benchmarkUuidv7(&randSourceBuf{}, b)
}

func BenchmarkUuidv7RandBufSwitch(b *testing.B) {
	rs, err := newRandSourceBufSwitch()
	if err != nil {
		log.Fatal(err)
	}
	benchmarkUuidv7(rs, b)
}
