package jsonparser

import "fmt"

var (
	ErrSyntax        = DecoderError{msg: "invalid character"}
	ErrUnexpectedEOF = DecoderError{msg: "unexpected end of JSON input"}
)

type errPos [2]int

type DecoderError struct {
	msg     string
	context string
	pos     errPos
	atChar  byte
	readErr error
}

func (e DecoderError) ReaderErr() error {
	return e.readerErr
}

func (e DecoderError) Error() string {
	loc := fmt.Sprintf("%s [%d,%d]", quoteChar(e.atChar), e.pos[0], e.pos[1])
	s := fmt.Sprintf("%s %s: %s", e.msg, e.context, loc)
	if e.readerErr != nil {
		s += "\nreader error: " + e.readerErr.Error()
	}
	return s
}
