package protobuf

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	testproto "github.com/mars9/protobuf/internal/proto"
)

func TestTypeEncode(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)

	for _, v := range typesMessages {
		buf.Reset()

		val := reflect.ValueOf(v)
		fields := val.Elem().NumField()
		if err := enc.encodeStruct(val.Elem(), fields); err != nil {
			t.Fatalf("encode type: %v", err)
		}

		data, err := proto.Marshal(v)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}
		if bytes.Compare(buf.Bytes(), data) != 0 {
			t.Fatalf("encode type: expected bytes %q, got %q", data, buf.Bytes())
		}

		m := &testproto.TypesMessage{}
		if err := proto.Unmarshal(buf.Bytes(), m); err != nil {
			t.Fatalf("unmarshal protobuf: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("encode type: expected %#v, got %#v", v, m)
		}
	}
}

func TestSliceEncode(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)

	for _, v := range sliceMessages {
		buf.Reset()

		val := reflect.ValueOf(v)
		fields := val.Elem().NumField()
		if err := enc.encodeStruct(val.Elem(), fields); err != nil {
			t.Fatalf("encode slice: %v", err)
		}

		data, err := proto.Marshal(v)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		if bytes.Compare(buf.Bytes(), data) != 0 {
			t.Fatalf("encode slice: expected bytes %q, got %q", data, buf.Bytes())
		}

		m := &testproto.SliceMessage{}
		if err := proto.Unmarshal(buf.Bytes(), m); err != nil {
			t.Fatalf("unmarshal protobuf: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("encode slice: expected %#v, got %#v", v, m)
		}
	}
}

func TestStructEncode(t *testing.T) {
	t.Parallel()

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)

	for _, v := range structMessages {
		buf.Reset()

		val := reflect.ValueOf(v)
		fields := val.Elem().NumField()
		if err := enc.encodeStruct(val.Elem(), fields); err != nil {
			t.Fatalf("encode struct: %v", err)
		}

		data, err := proto.Marshal(v)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		if bytes.Compare(buf.Bytes(), data) != 0 {
			t.Fatalf("encode struct: expected bytes %q, got %q", data, buf.Bytes())
		}

		m := &testproto.StructMessage{}
		if err := proto.Unmarshal(buf.Bytes(), m); err != nil {
			t.Fatalf("unmarshal protobuf: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("encode struct: expected %#v, got %#v", v, m)
		}
	}
}
