package protobuf

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Uint32  uint32
	Uint64  uint64
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64
	Bool    bool
	String  string

	Uint32Slice  []uint32
	Uint64Slice  []uint64
	Int32Slice   []int32
	Int64Slice   []int64
	Float32Slice []float32
	Float64Slice []float64
	StringSlice  []string
}

type testTagStruct struct {
	Uint32  uint32  `protobuf:"varint,required"`
	Uint64  uint64  `protobuf:"varint,required"`
	Int32   int32   `protobuf:"varint,required"`
	Int64   int64   `protobuf:"varint,required"`
	Float32 float32 `protobuf:"varint,required"`
	Float64 float64 `protobuf:"varint,required"`
	Bool    bool    `protobuf:"varint,required"`

	Sfixed32 int32  `protobuf:"sfixed32,required"`
	Sfixed64 int64  `protobuf:"sfixed64,required"`
	Fixed32  uint32 `protobuf:"fixed32,required"`
	Fixed64  uint64 `protobuf:"fixed64,required"`

	Sfixed32Slice []int32  `protobuf:"sfixed32,repeated"`
	Sfixed64Slice []int64  `protobuf:"sfixed64,repeated"`
	Fixed32Slice  []uint32 `protobuf:"fixed32,repeated"`
	Fixed64Slice  []uint64 `protobuf:"fixed64,repeated"`
}

func TestMarshalUnmarshal(t *testing.T) {
	var v = testStruct{
		Uint32:  42,
		Uint64:  42,
		Int32:   -42,
		Int64:   -42,
		Float32: 42.0,
		Float64: 42.0,
		Bool:    true,
		String:  "string",

		Uint32Slice:  []uint32{40, 41, 42, 43, 44},
		Uint64Slice:  []uint64{40, 41, 42, 43, 44},
		Int32Slice:   []int32{-40, -41, -42, -43, -44},
		Int64Slice:   []int64{-40, -41, -42, -43, -44},
		Float32Slice: []float32{40.0, 41.0, 42.0, 43.0, 44.0},
		Float64Slice: []float64{40.0, 41.0, 42.0, 43.0, 44.0},
		StringSlice:  []string{"string1", "string2", "string3"},
	}

	size := Size(&v)
	data := make([]byte, size)
	n, err := Marshal(data, &v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if size != n {
		t.Fatalf("marshal bytes size: expected %d, got %d", size, n)
	}

	var x = testStruct{}
	if err := Unmarshal(data, &x); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(v, x) {
		t.Fatalf("marshal/unmarshal: expected %#v, got %#v", v, x)
	}
}

func TestTagMarshalUnmarshal(t *testing.T) {
	var v = testTagStruct{
		Uint32:  42,
		Uint64:  42,
		Int32:   -42,
		Int64:   -42,
		Float32: 42.0,
		Float64: 42.0,
		Bool:    true,

		Sfixed32: -42,
		Sfixed64: -42,
		Fixed32:  42,
		Fixed64:  42,

		Sfixed32Slice: []int32{-40, -41, -42, -43, -44},
		Sfixed64Slice: []int64{-40, -41, -42, -43, -44},
		Fixed32Slice:  []uint32{40, 41, 42, 43, 44},
		Fixed64Slice:  []uint64{40, 41, 42, 43, 44},
	}

	size := Size(&v)
	data := make([]byte, size)
	n, err := Marshal(data, &v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if size != n {
		t.Fatalf("marshal bytes size: expected %d, got %d", size, n)
	}

	var x = testTagStruct{}
	if err := Unmarshal(data, &x); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(v, x) {
		t.Fatalf("marshal/unmarshal: expected %#v, got %#v", v, x)
	}
}

func TestRequiredField(t *testing.T) {
	t.Parallel()

	var v = struct {
		Required *uint32 `protobuf:"varint,required"`
	}{}

	data := make([]byte, Size(&v))
	if _, err := Marshal(data, &v); err == nil {
		t.Fatalf("expected required field not set error")
	} else {
		if err.Error() != "required field not set" {
			t.Fatalf(`expected error "required field not set", got %q`, err.Error())
		}
	}

	x := uint32(42)
	v.Required = &x
	data = make([]byte, Size(&v))
	if _, err := Marshal(data, &v); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}

func TestRequiredTagField(t *testing.T) {
	t.Parallel()

	var v = struct {
		Required *uint32 `protobuf:"fixed32,required"`
	}{}

	data := make([]byte, Size(&v))
	if _, err := Marshal(data, &v); err == nil {
		t.Fatalf("expected required field not set error")
	} else {
		if err.Error() != "required field not set" {
			t.Fatalf(`expected error "required field not set", got %q`, err.Error())
		}
	}

	x := uint32(42)
	v.Required = &x
	data = make([]byte, Size(&v))
	if _, err := Marshal(data, &v); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}
