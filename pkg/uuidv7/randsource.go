package uuidv7

import (
	"crypto/rand"
	"fmt"
)

type randSource interface {
	Read([]byte) (int, error)
}

const randPoolSize = 32768

type randBuf struct {
	buf    []byte
	origin []byte
	ok     bool
}
type randSourceBufSwitch struct {
	bufs       [2]randBuf
	currentIdx uint
	refill     chan *randBuf
	refilled   chan randRefillMsg
}
type randRefillMsg struct {
	rbuf *randBuf
	err  error
}

func (rs *randSourceBufSwitch) Read(b []byte) (n int, err error) {
	currentBuf := &rs.bufs[rs.currentIdx]

	toread := cap(b)
	if toread > randPoolSize {
		return 0, fmt.Errorf("Invalid size")
	} else if toread == 0 {
		return 0, fmt.Errorf("nocap")
	}

	if toread > len(currentBuf.buf) {
		// out of randomness, refill buf and switch to next one
		currentBuf.ok = false
		rs.refill <- currentBuf

		rs.currentIdx = (rs.currentIdx + 1) % 2
		currentBuf = &rs.bufs[rs.currentIdx]

		if currentBuf.ok == false {
			// block waiting for refilled if not ready
			msg := <-rs.refilled
			if msg.err != nil {
				return 0, msg.err
			}
			currentBuf = msg.rbuf
			currentBuf.ok = true
		}
	}

	copy(b[:], currentBuf.buf[:toread])
	currentBuf.buf = currentBuf.buf[toread:]

	return toread, nil
}

func (rs *randSourceBufSwitch) refiller() {
	var incoming *randBuf
	for {
		incoming = <-rs.refill

		incoming.buf = incoming.origin
		n, err := rand.Read(incoming.buf)
		if err != nil {
			rs.refilled <- randRefillMsg{incoming, err}
		}
		if n != randPoolSize {
			rs.refilled <- randRefillMsg{incoming, fmt.Errorf("rand.Read short read: %d", n)}
		}
		rs.refilled <- randRefillMsg{incoming, nil}
	}
}

func newRandSourceBufSwitch() (rs *randSourceBufSwitch, err error) {
	rs = &randSourceBufSwitch{}
	rs.refill = make(chan *randBuf)
	rs.refilled = make(chan randRefillMsg, 1)

	go rs.refiller()
	for i, b := range rs.bufs {
		b.buf = make([]byte, randPoolSize)
		b.origin = b.buf

		rs.refill <- &b
		refilled := <-rs.refilled
		if refilled.err != nil {
			return rs, refilled.err
		}
		rs.bufs[i] = *refilled.rbuf
	}

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
