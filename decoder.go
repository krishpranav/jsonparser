package jsonparser

import (
	"bytes"
	"encoding/json"
	"io"
	"sync/atomic"
)

type ValueType int

const (
	Unknown ValueType = iota
	Null
	String
	Number
	Boolean
	Array
	Object
)

type MetaValue struct {
	Offset    int
	Length    int
	Depth     int
	Value     interface{}
	ValueType ValueType
}

type KV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type KVS []KV

func (kvs KVS) MarshalJSON() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Write([]byte("{"))
	for i, kv := range kvs {
		b.Write([]byte("\"" + kv.Key + "\"" + ":"))
		valBuf, err := json.Marshal(kv.Value)
		if err != nil {
			return nil, err
		}
		b.Write(valBuf)
		if i < len(kvs)-1 {
			b.Write([]byte(","))
		}
	}
	b.Write([]byte("}"))
	return b.Bytes(), nil
}

type Decoder struct {
	*scanner
	emitDepth     int
	emitKV        bool
	emitRecursive bool
	objectAsKVS   bool
	depth         int
	scratch       *scratch
	metaCh        chan *MetaValue
	err           error
	lineNo        int
	lineStart     int64
}

func NewDecoder(r io.Reader, emitDepth int) *Decoder {
	d := &Decoder{
		scanner:   newScanner(r),
		emitDepth: emitDepth,
		scratch:   &scratch{data: make([]byte, 1024)},
		metaCh:    make(chan *MetaValue, 128),
	}
	if emitDepth < 0 {
		d.emitDepth = 0
		d.emitRecursive = true
	}
	return d
}

func (d *Decoder) ObjectAsKVS() *Decoder {
	d.objectAsKVS = true
	return d
}

func (d *Decoder) EmitKV() *Decoder {
	d.emitKV = true
	return d
}

func (d *Decoder) Recursive() *Decoder {
	d.emitRecursive = true
	return d
}

func (d *Decoder) Stream() chan *MetaValue {
	go d.decode()
	return d.metaCh
}

func (d *Decoder) Pos() int {
	return int(d.pos)
}

func (d *Decoder) Err() error {
	return d.err
}

func (d *Decoder) decode() {
	defer close(d.metaCh)
	d.skipSpaces()
	for d.remaining() > 0 {
		_, err := d.emitAny()
		if err != nil {
			d.err = err
			break
		}
		d.skipSpaces()
	}
}

func (d *Decoder) emitAny() (interface{}, error) {
	if d.pos >= atomic.LoadInt64(&d.end) {
		return nil, d.mkError(ErrUnexpectedEOF)
	}
	offset := d.pos - 1
	i, t, err := d.any()
	if d.willEmit() {
		d.metaCh <- &MetaValue{
			Offset:    int(offset),
			Length:    int(d.pos - offset),
			Depth:     d.depth,
			Value:     i,
			ValueType: t,
		}
	}
	return i, err
}

func (d *Decoder) willEmit() bool {
	if d.emitRecursive {
		return d.depth >= d.emitDepth
	}
	return d.depth == d.emitDepth
}
