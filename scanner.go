package jsonparser

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
