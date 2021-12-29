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
}
