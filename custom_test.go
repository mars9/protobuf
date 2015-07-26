package protobuf

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

type testCustom struct {
	Uint32 uint32
	Time   time.Time
	Uint64 uint64
}

func TestCustomTypes(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)
	dec := NewDecoder(buf, 0)

	for i, v := range []testCustom{
		testCustom{},
		testCustom{Time: time.Now()},
		testCustom{Uint32: 42, Time: time.Now(), Uint64: 42},
	} {
		val := reflect.ValueOf(&v)
		n := sizeStruct(val.Elem(), val.Elem().NumField())
		if i == 0 && n != 0 {
			t.Fatalf("custom size: expected empty message: got %d", n)
		}

		if err := enc.Encode(&v); err != nil {
			t.Fatalf("custom encode: %v", err)
		}

		m := testCustom{}
		if err := dec.Decode(&m); err != nil {
			t.Fatalf("custom decode: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("custom decode: expected %#v, got %#v", v, m)
		}
	}
}

type customInt32 int32

type testCustomType struct {
	Type customInt32
}

func TestCustomType(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)
	dec := NewDecoder(buf, 0)

	if err := enc.Encode(&testCustomType{customInt32(42)}); err != nil {
		t.Fatalf("encode: %v", err)
	}

	m := testCustomType{}
	if err := dec.Decode(&m); err != nil {
		t.Fatalf("decode nil: %v", err)
	}

	if m.Type != 42 {
		t.Fatalf("decode custom type: expected 42, got %d", m.Type)
	}
}

type testNil struct {
	Data []byte
}

func TestNilDecode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)
	dec := NewDecoder(buf, 0)

	if err := enc.Encode(&testNil{testBytes[:]}); err != nil {
		t.Fatalf("encode: %v", err)
	}

	if err := dec.Decode(nil); err != nil {
		t.Fatalf("decode nil: %v", err)
	}

	if buf.Len() != 0 {
		t.Fatalf("decode nil: expected empty buffer, got %d", buf.Len())
	}
}
