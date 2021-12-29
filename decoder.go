package jsonparser

import (
	"bytes"
	"encoding/json"
	"io"
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
