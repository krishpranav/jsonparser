package jsonparser

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
