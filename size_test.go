package protobuf

import (
	"math"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mars9/protobuf/internal/testproto"
)

var protoTypes = []struct {
	Uint32  uint32
	Uint64  uint64
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64
	Bool    bool
	String  string
}{
	{Uint32: 0, Uint64: 0, Int32: 0, Int64: 0, Float32: 0, Float64: 0, Bool: false, String: ""},
	{Uint32: 1, Uint64: 1, Int32: 1, Int64: 1, Float32: 1, Float64: 1, Bool: true, String: "string"},
	{Uint32: 1, Uint64: 1, Int32: -1, Int64: -1, Float32: -1, Float64: -1, Bool: true, String: "string"},
	{Uint32: math.MaxUint32, Uint64: math.MaxUint64, Int32: math.MaxInt32, Int64: math.MaxInt64,
		Float32: math.MaxFloat32, Float64: math.MaxFloat64, Bool: true, String: "string"},
	{Uint32: 0, Uint64: 0, Int32: math.MinInt32, Int64: math.MinInt64,
		Float32: math.SmallestNonzeroFloat32, Float64: math.SmallestNonzeroFloat64, Bool: true,
		String: "string"},
}

var protoBytes = []struct {
	Bytes      []byte
	BytesSlice [][]byte
}{
	{Bytes: []byte{}, BytesSlice: [][]byte{}},
	{Bytes: []byte("bytes"), BytesSlice: [][]byte{[]byte("bytes"), []byte("bytes"), []byte("bytes")}},
}

var protoSliceSize = struct {
	Uint32Slice  []uint32
	Uint64Slice  []uint64
	Int32Slice   []int32
	Int64Slice   []int64
	Float32Slice []float32
	Float64Slice []float64
	BoolSlice    []bool
	StringSlice  []string
}{}

func init() {
	for _, m := range protoTypes {
		protoSliceSize.Uint32Slice = append(protoSliceSize.Uint32Slice, m.Uint32)
		protoSliceSize.Uint64Slice = append(protoSliceSize.Uint64Slice, m.Uint64)
		protoSliceSize.Int32Slice = append(protoSliceSize.Int32Slice, m.Int32)
		protoSliceSize.Int64Slice = append(protoSliceSize.Int64Slice, m.Int64)
		protoSliceSize.Float32Slice = append(protoSliceSize.Float32Slice, m.Float32)
		protoSliceSize.Float64Slice = append(protoSliceSize.Float64Slice, m.Float64)
		protoSliceSize.BoolSlice = append(protoSliceSize.BoolSlice, m.Bool)
		protoSliceSize.StringSlice = append(protoSliceSize.StringSlice, m.String)
	}
}

type pstruct struct {
	Uint32 uint32
	Uint64 uint64
}

var protoStruct = []struct {
	Struct1 pstruct
	Struct2 pstruct
}{
	{Struct1: pstruct{}, Struct2: pstruct{}},
	{Struct1: pstruct{Uint32: 1, Uint64: 1}, Struct2: pstruct{Uint32: 1, Uint64: 1}},
}

func TestTypesSize(t *testing.T) {
	t.Parallel()

	for _, m := range protoTypes {
		pb := &testproto.ProtoTypes{
			Uint32:  proto.Uint32(m.Uint32),
			Uint64:  proto.Uint64(m.Uint64),
			Int32:   proto.Int32(m.Int32),
			Int64:   proto.Int64(m.Int64),
			Float32: proto.Float32(m.Float32),
			Float64: proto.Float64(m.Float64),
			Bool:    proto.Bool(m.Bool),
			String_: proto.String(m.String),
		}
		pbSize := proto.Size(pb)
		size := Size(&m)
		if pbSize != size {
			t.Fatalf("expected type size %d, got %d", pbSize, size)
		}
	}
}

func TestPointerTypeSize(t *testing.T) {
	t.Parallel()

	for _, m := range protoTypes {
		pb := &testproto.ProtoTypes{
			Uint32:  proto.Uint32(m.Uint32),
			Uint64:  proto.Uint64(m.Uint64),
			Int32:   proto.Int32(m.Int32),
			Int64:   proto.Int64(m.Int64),
			Float32: proto.Float32(m.Float32),
			Float64: proto.Float64(m.Float64),
			Bool:    proto.Bool(m.Bool),
			String_: proto.String(m.String),
		}
		var b = struct {
			Uint32  *uint32
			Uint64  *uint64
			Int32   *int32
			Int64   *int64
			Float32 *float32
			Float64 *float64
			Bool    *bool
			String  *string
		}{
			Uint32:  &m.Uint32,
			Uint64:  &m.Uint64,
			Int32:   &m.Int32,
			Int64:   &m.Int64,
			Float32: &m.Float32,
			Float64: &m.Float64,
			Bool:    &m.Bool,
			String:  &m.String,
		}

		pbSize := proto.Size(pb)
		size := Size(&b)
		if pbSize != size {
			t.Fatalf("expected pointer type size %d, got %d", pbSize, size)
		}
	}
}

func TestByteSliceSize(t *testing.T) {
	t.Parallel()

	for _, m := range protoBytes {
		pb := &testproto.ProtoBytes{
			Bytes:      m.Bytes,
			BytesSlice: m.BytesSlice,
		}

		pbSize := proto.Size(pb)
		size := Size(&m)
		if pbSize != size {
			t.Fatalf("expected bytes size %d, got %d", pbSize, size)
		}
	}
}

func TestSliceSize(t *testing.T) {
	t.Parallel()

	pb := &testproto.ProtoSlice{
		Uint32Slice:  protoSliceSize.Uint32Slice,
		Uint64Slice:  protoSliceSize.Uint64Slice,
		Int32Slice:   protoSliceSize.Int32Slice,
		Int64Slice:   protoSliceSize.Int64Slice,
		Float32Slice: protoSliceSize.Float32Slice,
		Float64Slice: protoSliceSize.Float64Slice,
		BoolSlice:    protoSliceSize.BoolSlice,
		StringSlice:  protoSliceSize.StringSlice,
	}

	pbSize := proto.Size(pb)
	size := Size(&protoSliceSize)
	if pbSize != size {
		t.Fatalf("expected slice size %d, got %d", pbSize, size)
	}
}

func TestStructSize(t *testing.T) {
	t.Parallel()

	for _, m := range protoStruct {
		pb := &testproto.ProtoStruct{
			Struct1: &testproto.ProtoStruct_Struct{
				Uint32: proto.Uint32(m.Struct1.Uint32),
				Uint64: proto.Uint64(m.Struct1.Uint64),
			},
			Struct2: &testproto.ProtoStruct_Struct{
				Uint32: proto.Uint32(m.Struct2.Uint32),
				Uint64: proto.Uint64(m.Struct2.Uint64),
			},
		}

		pbSize := proto.Size(pb)
		size := Size(&m)
		if pbSize != size {
			t.Fatalf("expected struct size %d, got %d", pbSize, size)
		}
	}
}

func TestPointerStructSize(t *testing.T) {
	t.Parallel()

	for _, m := range protoStruct {
		pb := &testproto.ProtoStruct{
			Struct1: &testproto.ProtoStruct_Struct{
				Uint32: proto.Uint32(m.Struct1.Uint32),
				Uint64: proto.Uint64(m.Struct1.Uint64),
			},
			Struct2: &testproto.ProtoStruct_Struct{
				Uint32: proto.Uint32(m.Struct2.Uint32),
				Uint64: proto.Uint64(m.Struct2.Uint64),
			},
		}

		var b = struct {
			Struct1 *pstruct
			Struct2 *pstruct
		}{Struct1: &m.Struct1, Struct2: &m.Struct2}

		pbSize := proto.Size(pb)
		size := Size(&b)
		if pbSize != size {
			t.Fatalf("expected pointer struct size %d, got %d", pbSize, size)
		}
	}
}
