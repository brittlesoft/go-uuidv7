package uuidv7

import (
	"io"
	"log"
	"testing"

	"github.com/art4711/unpredictable"
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

func BenchmarkUuidv7RandUnpredictable(b *testing.B) {
	rs := unpredictable.NewReader()
	benchmarkUuidv7(rs, b)
}

func uuidv4Unpredictable(rs io.Reader) [16]byte {
	var r [16]byte = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
	rs.Read(r[:])
	return r
}

func BenchmarkUnpredictableUuidv4(b *testing.B) {
	b.ReportAllocs()

	rs := unpredictable.NewReader()
	for n := 0; n < b.N; n++ {
		r := uuidv4Unpredictable(rs)
		u.B = r
	}
}
