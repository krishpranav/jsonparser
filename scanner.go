package jsonparser

import (
	"io"
	"sync/atomic"
)

const (
	chunk    = 4095
	maxUint  = ^uint(0)
	maxInt   = int64(maxUint >> 1)
	nullByte = byte(0)
)

type scanner struct {
	pos       int64
	ipos      int64
	ifill     int64
	end       int64
	buf       [chunk + 1]byte
	nbuf      [chunk]byte
	fillReq   chan struct{}
	fillReady chan int64
	readerErr error
}

func newScanner(r io.Reader) *scanner {
	sr := &scanner{
		end:       maxInt,
		fillReq:   make(chan struct{}),
		fillReady: make(chan int64),
	}

	go func() {
		var rpos int64

		defer func() {
			atomic.StoreInt64(&sr.end, rpos)
			close(sr.fillReady)
		}()

		for range sr.fillReq {
		scan:
			n, err := r.Read(sr.nbuf[:])

			if n == 0 {
				switch err {
				case io.EOF:
					return
				case nil:
					goto scan
				default:
					sr.readerErr = err
					return
				}
			}

			rpos += int64(n)
			sr.fillReady <- int64(n)
		}
	}()

	sr.fillReq <- struct{}{}

	return sr
}

func (s *scanner) remaining() int64 {
	if atomic.LoadInt64(&s.end) == maxInt {
		return maxInt
	}
	return atomic.LoadInt64(&s.end) - s.pos
}

func (s *scanner) cur() byte {
	return s.buf[s.ipos]
}
