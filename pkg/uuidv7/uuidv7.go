package uuidv7

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	version      = 0x7
	variant      = 0x2
	subsecFactor = 16777216 // 2**24
)

type DecodedUuidv7 struct {
	// https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html#name-uuidv7-layout-and-bit-order
	Unixts        uint64  // 36-bit big-endian unsigned Unix Timestamp value
	SubsecA       uint16  // 12-bits allocated to sub-section precision values.
	Ver           uint8   // The 4 bit UUIDv8 version (0111)
	SubsecB       uint16  // 12-bits allocated to sub-section precision values.
	Variant       uint8   // 2-bit UUID variant (10)
	SubsecSeqNode uint64  // The remaining 62 bits which MAY be allocated to any combination of additional sub-section precision, sequence counter, or pseudo-random data.
	Ts            float64 // reconstructed timestamp (to ms precision)
}

func NewDecodedUuidv7(b [16]byte) (du DecodedUuidv7) {
	du.Unixts = binary.BigEndian.Uint64(b[:8]) >> 28
	du.SubsecA = binary.BigEndian.Uint16(b[4:6]) & 0x0fff
	du.Ver = uint8(b[6]) >> 4
	du.SubsecB = binary.BigEndian.Uint16(b[6:8]) & 0x0fff
	du.Variant = uint8(b[8]) >> 6
	du.SubsecSeqNode = binary.BigEndian.Uint64(b[8:]) & 0x3fffffffffffffff
	du.Ts = float64(uint64(du.SubsecA)<<12|uint64(du.SubsecB))/subsecFactor + float64(du.Unixts)
	return du
}

type Uuidv7 struct {
	ts uint64
	B  [16]byte
}

func (u *Uuidv7) String() string {
	// 8-4-4-4-12 in hex
	return fmt.Sprintf("%X-%X-%X-%X-%X", u.B[:4], u.B[4:6], u.B[6:8], u.B[8:10], u.B[10:])
}

func (u *Uuidv7) Ts() float64 {
	return float64(u.ts) / 1_000_000_000
}

type Uuidv7Source struct {
	rs randSource
}

func NewUuidv7Source() (us *Uuidv7Source, err error) {
	us = &Uuidv7Source{}
	us.rs, err = newRandSourceBufSwitch()
	return us, err
}

func (us *Uuidv7Source) New() (u Uuidv7, err error) {
	// RFC says unixts is big endian but doesn't say anything about the rest
	// here everything is big endian
	var r [8]byte // last 64 bytes -> 2bit variant + 62bit random

	u.ts = uint64(time.Now().UnixNano())

	// 36bits ts in second
	// 12bits subsec_a
	//  4bits version
	// 12bits subsec_b
	//  2bits variant
	// 62bits subsecseqnode  (here, random)

	// 36-bit big-endian unsigned Unix Timestamp value
	// convert to seconds, store the top 36 bits of the result in the top 36 bits of b
	var b uint64 = (u.ts / 1_000_000_000) << 28 & 0xfffffffff0000000

	// subsec encoding
	// https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html#name-uuidv7-decoding
	// ns total -> ns int -> ns frac
	var ns_frac float64 = float64(u.ts%1_000_000_000) / 1_000_000_000
	// ns_frac -> subsec encoding on 24bit uint
	var subsec uint64 = uint64(ns_frac * subsecFactor)

	// 12-bits subsec A
	// clear 12 lower subsec bits, shift to 36th bit from the left
	b = b | subsec>>12<<16

	// The 4 bit UUIDv7 version (0111)
	b = b | version<<12

	// lower 12-bits of subsec
	b = b | (subsec & 0x0000000000000fff)

	binary.BigEndian.PutUint64(u.B[:], b)

	// 64bits random
	_, err = us.rs.Read(r[:])
	if err != nil {
		return u, err
	}

	// Variant in top 2bits
	r[0] = (r[0] & 0x3f) | variant<<6

	binary.BigEndian.PutUint64(u.B[8:], binary.LittleEndian.Uint64(r[:]))

	return u, nil
}
