package jsonparser

import (
	"bytes"
	"testing"
)

func mkReader(s string) *bytes.Reader {
	return bytes.NewReader([]byte(s))
}

func TestDecoderSimple(t *testing.T) {
	var (
		counter int
		mv      *MetaValue
		body    = `[{"bio":"bada bing bada boom","id":1,"name":"Charles","falseVal":false}]`
	)

	decoder := NewDecoder(mkReader(body), 1)

	for mv = range decoder.Stream() {
		counter++
		t.Logf("depth=%d offset=%d len=%d (%v)", mv.Depth, mv.Offset, mv.Length, mv.Value)
	}

	if err := decoder.Err(); err != nil {
		t.Fatalf("decoder error: %s", err)
	}
}

func TestDecoderNested(t *testing.T) {
	var (
		counter int
		mv      *MetaValue
		body    = `{
  "1": {
    "bio": "bada bing bada boom",
    "id": 0,
    "name": "Roberto",
    "nested1": {
      "bio": "utf16 surrogate (\ud834\udcb2)\n\u201cutf 8\u201d",
      "id": 1.5,
      "name": "Roberto*Maestro",
      "nested2": { "nested2arr": [0,1,2], "nested3": {
        "nested4": { "depth": "recursion" }}
			}
		}
  },
  "2": {
    "nullfield": null,
    "id": -2
  }
}`
	)

	decoder := NewDecoder(mkReader(body), 2)

	for mv = range decoder.Stream() {
		counter++
		t.Logf("depth=%d offset=%d len=%d (%v)", mv.Depth, mv.Offset, mv.Length, mv.Value)
	}

	if err := decoder.Err(); err != nil {
		t.Fatalf("decoder error: %s", err)
	}
}

func TestDecoderFlat(t *testing.T) {
	var (
		counter int
		mv      *MetaValue
		body    = `[
  "1st test string",
  "Roberto*Maestro", "Charles",
  0, null, false,
  1, 2.5
]`
		expected = []struct {
			Value     interface{}
			ValueType ValueType
		}{
			{
				"1st test string",
				String,
			},
			{
				"Roberto*Maestro",
				String,
			},
			{
				"Charles",
				String,
			},
			{
				0.0,
				Number,
			},
			{
				nil,
				Null,
			},
			{
				false,
				Boolean,
			},
			{
				1.0,
				Number,
			},
			{
				2.5,
				Number,
			},
		}
	)

	decoder := NewDecoder(mkReader(body), 1)

	for mv = range decoder.Stream() {
		if mv.Value != expected[counter].Value {
			t.Fatalf("got %v, expected: %v", mv.Value, expected[counter])
		}
		if mv.ValueType != expected[counter].ValueType {
			t.Fatalf("got %v value type, expected: %v value type", mv.ValueType, expected[counter].ValueType)
		}
		counter++
		t.Logf("depth=%d offset=%d len=%d (%v)", mv.Depth, mv.Offset, mv.Length, mv.Value)
	}

	if err := decoder.Err(); err != nil {
		t.Fatalf("decoder error: %s", err)
	}
}

func TestDecoderReaderFailure(t *testing.T) {
	var (
		failAfter = 900
		mockData  = byte('[')
	)

	r := newMockReader(failAfter, mockData)
	decoder := NewDecoder(r, -1)

	for mv := range decoder.Stream() {
		t.Logf("depth=%d offset=%d len=%d (%v)", mv.Depth, mv.Offset, mv.Length, mv.Value)
	}

	err := decoder.Err()
	t.Logf("got error: %s", err)
	if err == nil {
		t.Fatalf("missing expected decoder error")
	}

	derr, ok := err.(DecoderError)
	if !ok {
		t.Fatalf("expected error of type DecoderError, got %T", err)
	}

	if derr.ReaderErr() == nil {
		t.Fatalf("missing expected underlying reader error")
	}
}
