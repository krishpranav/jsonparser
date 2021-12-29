package jsonparser

import (
	"bytes"
	"encoding/json"
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
