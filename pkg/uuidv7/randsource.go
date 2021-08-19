// This file is not used, replaced by `unpredictable`
// kept here for benchmarks

package uuidv7

import (
	"crypto/rand"
	"fmt"
)

type randSource interface {
	Read([]byte) (int, error)
}

const (
	randPoolSize    = 32768
	maxBufQueueSize = 32
)

type randBuf struct {
	buf    []byte
	origin []byte
}
type randSourceBufQueue struct {
	currentBuf *randBuf
	usedQ      chan *randBuf
	readyQ     chan randRefillMsg
}
type randRefillMsg struct {
	rbuf *randBuf
	err  error
}

func (rs *randSourceBufQueue) Read(b []byte) (n int, err error) {
	toread := cap(b)
	if toread > randPoolSize {
		return 0, fmt.Errorf("Invalid size")
	} else if toread == 0 {
		return 0, fmt.Errorf("nocap")
	}

	if toread > len(rs.currentBuf.buf) {
		// out of randomness, refill buf
		rs.usedQ <- rs.currentBuf

		msg := <-rs.readyQ
		rs.currentBuf = msg.rbuf
		if msg.err != nil {
			// on error currentBuf is likely still empty/unuseable, next read
			// will push it back to usedQ
			return 0, msg.err
		}
	}

	copy(b[:], rs.currentBuf.buf[:toread])
	rs.currentBuf.buf = rs.currentBuf.buf[toread:]

	return toread, nil
}

func (rs *randSourceBufQueue) refiller() {
	allocated := 0
	var incoming *randBuf
	for {
		if allocated < maxBufQueueSize {
			incoming = &randBuf{
				buf: make([]byte, randPoolSize),
			}
			incoming.origin = incoming.buf
			allocated++
		} else {
			// block receiving once we allocated all maxBufQueueSize buffers
			incoming = <-rs.usedQ
		}

		incoming.buf = incoming.origin

		n, err := rand.Read(incoming.buf)

		if err != nil {
			rs.readyQ <- randRefillMsg{incoming, err}
		}
		if n != randPoolSize {
			rs.readyQ <- randRefillMsg{incoming, fmt.Errorf("rand.Read short read: %d", n)}
		}
		rs.readyQ <- randRefillMsg{incoming, nil}
	}
}

func newRandSourceBufQueue() (rs *randSourceBufQueue, err error) {
	rs = &randSourceBufQueue{}
	rs.usedQ = make(chan *randBuf, maxBufQueueSize)
	rs.readyQ = make(chan randRefillMsg, maxBufQueueSize)

	go rs.refiller()
	msg := <-rs.readyQ
	if msg.err != nil {
		return nil, msg.err
	}
	rs.currentBuf = msg.rbuf

	return rs, nil
}

type randSourceBuf struct {
	buf randBuf
}

func (rs *randSourceBuf) Read(b []byte) (n int, err error) {
	toread := cap(b)
	if toread > randPoolSize {
		return 0, fmt.Errorf("Invalid size")
	} else if toread == 0 {
		return 0, fmt.Errorf("nocap")
	}

	if toread > len(rs.buf.buf) {
		rs.refill() // blocks
	}

	copy(b[:], rs.buf.buf[:toread])
	rs.buf.buf = rs.buf.buf[toread:]

	return toread, nil
}

func (rs *randSourceBuf) refill() (err error) {
	if rs.buf.buf == nil {
		b := make([]byte, randPoolSize)
		rs.buf.buf = b
		rs.buf.origin = rs.buf.buf
	} else {
		rs.buf.buf = rs.buf.origin
	}

	n, err := rand.Read(rs.buf.buf)
	if err != nil {
		return err
	}
	if n != randPoolSize {
		return fmt.Errorf("rand.Read short read: %d", n)
	}
	return nil
}

type randSourceStraight struct {
}

func (rs *randSourceStraight) Read(b []byte) (n int, err error) {
	return rand.Read(b)
}

type randSourceNoRand struct {
}

func (rs *randSourceNoRand) Read(b []byte) (n int, err error) {
	return len(b), nil
}
