package protobuf

import (
	"reflect"
	"testing"

	testproto "github.com/mars9/protobuf/internal/proto"
)

func TestEncodeDecodeBuffer(t *testing.T) {
	t.Parallel()

	enc := NewEncodeBuffer()
	dec := NewDecodeBuffer()

	for _, v := range structMessages {
		enc.Reset()
		if err := enc.Encode(v); err != nil {
			t.Fatalf("encode buffer: %v", err)
		}

		m := &testproto.StructMessage{}

		dec.Write(enc.Bytes())
		if err := dec.Decode(m); err != nil {
			t.Fatalf("decode buffer: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("buffer enc/dec: expected %#v, got %#v", v, m)
		}
	}
}
