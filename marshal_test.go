// +build goprotobuf

package protobuf

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mars9/protobuf/internal/testproto"
)

var protoSliceMarshal = struct {
	Uint32Slice  []uint32
	Uint64Slice  []uint64
	Int32Slice   []int32
	Int64Slice   []int64
	Float32Slice []float32
	Float64Slice []float64
	BoolSlice    []bool
	StringSlice  []string
}{}

var protoTagsSliceMarshal = struct {
	Sfixed32Slice []int32  `protobuf:"sfixed32,repeated"`
	Sfixed64Slice []int64  `protobuf:"sfixed64,repeated"`
	Fixed32Slice  []uint32 `protobuf:"fixed32,repeated"`
	Fixed64Slice  []uint64 `protobuf:"fixed64,repeated"`
}{}

func init() {
	for _, m := range protoTypes {
		protoSliceMarshal.Uint32Slice = append(protoSliceMarshal.Uint32Slice, m.Uint32)
		protoSliceMarshal.Uint64Slice = append(protoSliceMarshal.Uint64Slice, m.Uint64)
		protoSliceMarshal.Int32Slice = append(protoSliceMarshal.Int32Slice, m.Int32)
		protoSliceMarshal.Int64Slice = append(protoSliceMarshal.Int64Slice, m.Int64)
		protoSliceMarshal.Float32Slice = append(protoSliceMarshal.Float32Slice, m.Float32)
		protoSliceMarshal.Float64Slice = append(protoSliceMarshal.Float64Slice, m.Float64)
		protoSliceMarshal.BoolSlice = append(protoSliceMarshal.BoolSlice, m.Bool)
		protoSliceMarshal.StringSlice = append(protoSliceMarshal.StringSlice, m.String)
	}
	for _, m := range protoTags {
		protoTagsSliceMarshal.Sfixed32Slice = append(protoTagsSliceMarshal.Sfixed32Slice, m.Sfixed32)
		protoTagsSliceMarshal.Sfixed64Slice = append(protoTagsSliceMarshal.Sfixed64Slice, m.Sfixed64)
		protoTagsSliceMarshal.Fixed32Slice = append(protoTagsSliceMarshal.Fixed32Slice, m.Fixed32)
		protoTagsSliceMarshal.Fixed64Slice = append(protoTagsSliceMarshal.Fixed64Slice, m.Fixed64)
	}
}

func TestTypesMarshal(t *testing.T) {
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

		size := Size(&m)
		data := make([]byte, size)
		n, err := Marshal(data, &m)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if size != n {
			t.Fatalf("expected type size %d, got %d", size, n)
		}

		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}
		if n != len(pbData) {
			t.Fatalf("expected type size %d, got %d", size, n)
		}
		if bytes.Compare(data, pbData) != 0 {
			t.Fatalf("expected type bytes %#v, got %#v", pbData, data)
		}
	}
}

func TestPointerTypeMarshal(t *testing.T) {
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

		size := Size(&b)
		data := make([]byte, size)
		n, err := Marshal(data, &b)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if size != n {
			t.Fatalf("expected pointer type size %d, got %d", size, n)
		}

		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}
		if n != len(pbData) {
			t.Fatalf("expected pointer type size %d, got %d", size, n)
		}
		if bytes.Compare(data, pbData) != 0 {
			t.Fatalf("expected pointer type bytes %#v, got %#v", pbData, data)
		}
	}
}

func TestByteSliceMarshal(t *testing.T) {
	t.Parallel()

	for _, m := range protoBytes {
		pb := &testproto.ProtoBytes{
			Bytes:      m.Bytes,
			BytesSlice: m.BytesSlice,
		}

		size := Size(&m)
		data := make([]byte, size)
		n, err := Marshal(data, &m)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if size != n {
			t.Fatalf("expected bytes size %d, got %d", size, n)
		}

		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}
		if n != len(pbData) {
			t.Fatalf("expected bytes size %d, got %d", size, n)
		}
		if bytes.Compare(data, pbData) != 0 {
			t.Fatalf("expected bytes bytes %#v, got %#v", pbData, data)
		}
	}
}

func TestSliceMarshal(t *testing.T) {
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

	size := Size(&protoSliceMarshal)
	data := make([]byte, size)
	n, err := Marshal(data, &protoSliceMarshal)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if size != n {
		t.Fatalf("expected slice size %d, got %d", size, n)
	}

	pbData, err := proto.Marshal(pb)
	if err != nil {
		t.Fatalf("protobuf marshal: %v", err)
	}
	if n != len(pbData) {
		t.Fatalf("expected slice size %d, got %d", size, n)
	}
	if bytes.Compare(data, pbData) != 0 {
		t.Fatalf("expected slice bytes %#v, got %#v", pbData, data)
	}
}

func TestStructMarshal(t *testing.T) {
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

		size := Size(&m)
		data := make([]byte, size)
		n, err := Marshal(data, &m)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if size != n {
			t.Fatalf("expected struct size %d, got %d", size, n)
		}

		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}
		if n != len(pbData) {
			t.Fatalf("expected struct size %d, got %d", size, n)
		}
		if bytes.Compare(data, pbData) != 0 {
			t.Fatalf("expected struct bytes %#v, got %#v", pbData, data)
		}
	}
}

func TestPointerStructMarshal(t *testing.T) {
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

		size := Size(&b)
		data := make([]byte, size)
		n, err := Marshal(data, &b)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if size != n {
			t.Fatalf("expected struct size %d, got %d", size, n)
		}

		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}
		if n != len(pbData) {
			t.Fatalf("expected struct size %d, got %d", size, n)
		}
		if bytes.Compare(data, pbData) != 0 {
			t.Fatalf("expected struct bytes %#v, got %#v", pbData, data)
		}
	}
}

func TestTagsMarshal(t *testing.T) {
	t.Parallel()

	for _, m := range protoTags {
		pb := &testproto.ProtoTags{
			Sfixed32: proto.Int32(m.Sfixed32),
			Sfixed64: proto.Int64(m.Sfixed64),
			Fixed32:  proto.Uint32(m.Fixed32),
			Fixed64:  proto.Uint64(m.Fixed64),
		}

		size := Size(&m)
		data := make([]byte, size)
		n, err := Marshal(data, &m)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if size != n {
			t.Fatalf("expected tag size %d, got %d", size, n)
		}

		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}
		if n != len(pbData) {
			t.Fatalf("expected tag size %d, got %d", size, n)
		}
		if bytes.Compare(data, pbData) != 0 {
			t.Fatalf("expected tag bytes %#v, got %#v", pbData, data)
		}

	}
}

func TestPointerTagsMarshal(t *testing.T) {
	t.Parallel()

	for _, m := range protoTags {
		pb := &testproto.ProtoTags{
			Sfixed32: proto.Int32(m.Sfixed32),
			Sfixed64: proto.Int64(m.Sfixed64),
			Fixed32:  proto.Uint32(m.Fixed32),
			Fixed64:  proto.Uint64(m.Fixed64),
		}

		var v = struct {
			Sfixed32 *int32  `protobuf:"sfixed32,required"`
			Sfixed64 *int64  `protobuf:"sfixed64,required"`
			Fixed32  *uint32 `protobuf:"fixed32,required"`
			Fixed64  *uint64 `protobuf:"fixed64,required"`
		}{
			Sfixed32: &m.Sfixed32,
			Sfixed64: &m.Sfixed64,
			Fixed32:  &m.Fixed32,
			Fixed64:  &m.Fixed64,
		}
		size := Size(&v)
		data := make([]byte, size)
		n, err := Marshal(data, &v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if size != n {
			t.Fatalf("expected tag size %d, got %d", size, n)
		}

		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}
		if n != len(pbData) {
			t.Fatalf("expected tag size %d, got %d", size, n)
		}
		if bytes.Compare(data, pbData) != 0 {
			t.Fatalf("expected tag bytes %#v, got %#v", pbData, data)
		}

	}
}

func TestTagsSliceMarshal(t *testing.T) {
	t.Parallel()

	pb := &testproto.ProtoTagsSlice{
		Sfixed32Slice: protoTagsSliceMarshal.Sfixed32Slice,
		Sfixed64Slice: protoTagsSliceMarshal.Sfixed64Slice,
		Fixed32Slice:  protoTagsSliceMarshal.Fixed32Slice,
		Fixed64Slice:  protoTagsSliceMarshal.Fixed64Slice,
	}

	size := Size(&protoTagsSliceMarshal)
	data := make([]byte, size)
	n, err := Marshal(data, &protoTagsSliceMarshal)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if size != n {
		t.Fatalf("expected tags slice size %d, got %d", size, n)
	}

	pbData, err := proto.Marshal(pb)
	if err != nil {
		t.Fatalf("protobuf marshal: %v", err)
	}
	if n != len(pbData) {
		t.Fatalf("expected tags slice size %d, got %d", size, n)
	}
	if bytes.Compare(data, pbData) != 0 {
		t.Fatalf("expected tags slice bytes %#v, got %#v", pbData, data)
	}
}

func TestStructSliceMarshal(t *testing.T) {
	t.Parallel()

	pb := &testproto.ProtoStructSlice{
		StructSlice: []*testproto.ProtoStructSlice_Struct{
			&testproto.ProtoStructSlice_Struct{Uint32: proto.Uint32(1<<32 - 1), Uint64: proto.Uint64(1<<64 - 1)},
			&testproto.ProtoStructSlice_Struct{Uint32: proto.Uint32(0), Uint64: proto.Uint64(0)},
			&testproto.ProtoStructSlice_Struct{Uint32: proto.Uint32(42), Uint64: proto.Uint64(42)},
		},
	}

	size := Size(&protoStructSlice)
	data := make([]byte, size)
	n, err := Marshal(data, &protoStructSlice)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if size != n {
		t.Fatalf("expected struct slice size %d, got %d", size, n)
	}

	pbData, err := proto.Marshal(pb)
	if err != nil {
		t.Fatalf("protobuf marshal: %v", err)
	}
	if n != len(pbData) {
		t.Fatalf("expected struct slice size %d, got %d", size, n)
	}
	if bytes.Compare(data, pbData) != 0 {
		t.Fatalf("expected struct slice bytes %#v, got %#v", pbData, data)
	}
}

func TestStructPtrSliceMarshal(t *testing.T) {
	t.Parallel()

	pb := &testproto.ProtoStructSlice{
		StructSlice: []*testproto.ProtoStructSlice_Struct{
			&testproto.ProtoStructSlice_Struct{Uint32: proto.Uint32(1<<32 - 1), Uint64: proto.Uint64(1<<64 - 1)},
			&testproto.ProtoStructSlice_Struct{Uint32: proto.Uint32(0), Uint64: proto.Uint64(0)},
			&testproto.ProtoStructSlice_Struct{Uint32: proto.Uint32(42), Uint64: proto.Uint64(42)},
		},
	}

	size := Size(&protoStructPtrSlice)
	data := make([]byte, size)
	n, err := Marshal(data, &protoStructPtrSlice)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if size != n {
		t.Fatalf("expected struct slice size %d, got %d", size, n)
	}

	pbData, err := proto.Marshal(pb)
	if err != nil {
		t.Fatalf("protobuf marshal: %v", err)
	}
	if n != len(pbData) {
		t.Fatalf("expected struct slice size %d, got %d", size, n)
	}
	if bytes.Compare(data, pbData) != 0 {
		t.Fatalf("expected struct slice bytes %#v, got %#v", pbData, data)
	}
}
