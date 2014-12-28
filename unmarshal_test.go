// +build goprotobuf

package protobuf

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mars9/protobuf/internal/testproto"
)

func TestTypesUnmarshal(t *testing.T) {
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
		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}

		var v = struct {
			Uint32  uint32
			Uint64  uint64
			Int32   int32
			Int64   int64
			Float32 float32
			Float64 float64
			Bool    bool
			String  string
		}{}
		if err = Unmarshal(pbData, &v); err != nil {
			t.Fatalf("unmarshal type: %v", err)
		}

		if !reflect.DeepEqual(m, v) {
			t.Fatalf("execpted type %#v, got %#v", m, v)
		}
	}
}

func TestPointerTypesUnmarshal(t *testing.T) {
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
		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}

		var v = struct {
			Uint32  *uint32
			Uint64  *uint64
			Int32   *int32
			Int64   *int64
			Float32 *float32
			Float64 *float64
			Bool    *bool
			String  *string
		}{}
		if err = Unmarshal(pbData, &v); err != nil {
			t.Fatalf("unmarshal type: %v", err)
		}

		if m.Uint32 != *v.Uint32 {
			t.Fatalf("expected pointer type %v, got %v", m.Uint32, v.Uint32)
		}
		if m.Uint64 != *v.Uint64 {
			t.Fatalf("expected pointer type %v, got %v", m.Uint64, v.Uint64)
		}
		if m.Int32 != *v.Int32 {
			t.Fatalf("expected pointer type %v, got %v", m.Int32, v.Int32)
		}
		if m.Int64 != *v.Int64 {
			t.Fatalf("expected pointer type %v, got %v", m.Int64, v.Int64)
		}
		if m.Float32 != *v.Float32 {
			t.Fatalf("expected pointer type %v, got %v", m.Float32, v.Float32)
		}
		if m.Float64 != *v.Float64 {
			t.Fatalf("expected pointer type %v, got %v", m.Float64, v.Float64)
		}
		if m.Bool != *v.Bool {
			t.Fatalf("expected pointer type %v, got %v", m.Bool, v.Bool)
		}
		if m.String != *v.String {
			t.Fatalf("expected pointer type %v, got %v", m.String, v.String)
		}
	}
}

func TestByteSliceUnmarshal(t *testing.T) {
	t.Parallel()

	for _, m := range protoBytes {
		pb := &testproto.ProtoBytes{
			Bytes:      m.Bytes,
			BytesSlice: m.BytesSlice,
		}
		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}

		var v = struct {
			Bytes      []byte
			BytesSlice [][]byte
		}{}
		if err = Unmarshal(pbData, &v); err != nil {
			t.Fatalf("unmarshal bytes: %v", err)
		}

		if len(m.BytesSlice) == 0 {
			if len(v.BytesSlice) != 0 {
				t.Fatalf("execpted bytes %#v, got %#v", m, v)
			}
			continue
		}
		if !reflect.DeepEqual(m, v) {
			t.Fatalf("execpted bytes %#v, got %#v", m, v)
		}
	}
}

func TestSliceUnmarshal(t *testing.T) {
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
	pbData, err := proto.Marshal(pb)
	if err != nil {
		t.Fatalf("protobuf marshal: %v", err)
	}

	var v = struct {
		Uint32Slice  []uint32
		Uint64Slice  []uint64
		Int32Slice   []int32
		Int64Slice   []int64
		Float32Slice []float32
		Float64Slice []float64
		BoolSlice    []bool
		StringSlice  []string
	}{}
	if err = Unmarshal(pbData, &v); err != nil {
		t.Fatalf("unmarshal bytes: %v", err)
	}

	if !reflect.DeepEqual(pb.Uint32Slice, v.Uint32Slice) {
		t.Fatalf("expected slice %v, got %v", pb.Uint32Slice, v.Uint32Slice)
	}
	if !reflect.DeepEqual(pb.Uint64Slice, v.Uint64Slice) {
		t.Fatalf("expected slice %v, got %v", pb.Uint64Slice, v.Uint64Slice)
	}
	if !reflect.DeepEqual(pb.Int32Slice, v.Int32Slice) {
		t.Fatalf("expected slice %v, got %v", pb.Int32Slice, v.Int32Slice)
	}
	if !reflect.DeepEqual(pb.Int64Slice, v.Int64Slice) {
		t.Fatalf("expected slice %v, got %v", pb.Int64Slice, v.Int64Slice)
	}
	if !reflect.DeepEqual(pb.Float32Slice, v.Float32Slice) {
		t.Fatalf("expected slice %v, got %v", pb.Float32Slice, v.Float32Slice)
	}
	if !reflect.DeepEqual(pb.Float64Slice, v.Float64Slice) {
		t.Fatalf("expected slice %v, got %v", pb.Float64Slice, v.Float64Slice)
	}
	if !reflect.DeepEqual(pb.BoolSlice, v.BoolSlice) {
		t.Fatalf("expected slice %v, got %v", pb.BoolSlice, v.BoolSlice)
	}
	if !reflect.DeepEqual(pb.StringSlice, v.StringSlice) {
		t.Fatalf("expected slice %v, got %v", pb.StringSlice, v.StringSlice)
	}
}

func TestStructUnmarshal(t *testing.T) {
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
		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}

		var v = struct {
			Struct1 pstruct
			Struct2 pstruct
		}{}
		if err = Unmarshal(pbData, &v); err != nil {
			t.Fatalf("unmarshal type: %v", err)
		}

		if !reflect.DeepEqual(m, v) {
			t.Fatalf("execpted struct %#v, got %#v", m, v)
		}
	}
}

func TestPointerStructUnmarshal(t *testing.T) {
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
		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}

		var v = struct {
			Struct1 *pstruct
			Struct2 *pstruct
		}{}
		if err = Unmarshal(pbData, &v); err != nil {
			t.Fatalf("unmarshal type: %v", err)
		}

		if !reflect.DeepEqual(&m.Struct1, v.Struct1) {
			t.Fatalf("execpted struct %#v, got %#v", &m.Struct1, v.Struct1)
		}
		if !reflect.DeepEqual(&m.Struct2, v.Struct2) {
			t.Fatalf("execpted struct %#v, got %#v", &m.Struct2, v.Struct2)
		}
	}
}

func TestTagsUnmarshal(t *testing.T) {
	t.Parallel()

	for _, m := range protoTags {
		pb := &testproto.ProtoTags{
			Sfixed32: proto.Int32(m.Sfixed32),
			Sfixed64: proto.Int64(m.Sfixed64),
			Fixed32:  proto.Uint32(m.Fixed32),
			Fixed64:  proto.Uint64(m.Fixed64),
		}
		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}

		var v = struct {
			Sfixed32 int32  `protobuf:"sfixed32,required"`
			Sfixed64 int64  `protobuf:"sfixed64,required"`
			Fixed32  uint32 `protobuf:"fixed32,required"`
			Fixed64  uint64 `protobuf:"fixed64,required"`
		}{}
		if err = Unmarshal(pbData, &v); err != nil {
			t.Fatalf("unmarshal tag: %v", err)
		}
		if !reflect.DeepEqual(m, v) {
			t.Fatalf("execpted tag %#v, got %#v", m, v)
		}
	}
}

func TestPointerTagsUnmarshal(t *testing.T) {
	t.Parallel()

	for _, m := range protoTags {
		pb := &testproto.ProtoTags{
			Sfixed32: proto.Int32(m.Sfixed32),
			Sfixed64: proto.Int64(m.Sfixed64),
			Fixed32:  proto.Uint32(m.Fixed32),
			Fixed64:  proto.Uint64(m.Fixed64),
		}
		pbData, err := proto.Marshal(pb)
		if err != nil {
			t.Fatalf("protobuf marshal: %v", err)
		}

		var v = struct {
			Sfixed32 *int32  `protobuf:"sfixed32,required"`
			Sfixed64 *int64  `protobuf:"sfixed64,required"`
			Fixed32  *uint32 `protobuf:"fixed32,required"`
			Fixed64  *uint64 `protobuf:"fixed64,required"`
		}{}
		if err = Unmarshal(pbData, &v); err != nil {
			t.Fatalf("unmarshal tag: %v", err)
		}
		if m.Sfixed32 != *v.Sfixed32 {
			t.Fatalf("expected pointer tag %#v, got %#v", m.Sfixed32, v.Sfixed32)
		}
		if m.Sfixed64 != *v.Sfixed64 {
			t.Fatalf("expected pointer tag %#v, got %#v", m.Sfixed64, v.Sfixed64)
		}
		if m.Fixed32 != *v.Fixed32 {
			t.Fatalf("expected pointer tag %#v, got %#v", m.Fixed32, v.Fixed32)
		}
		if m.Fixed64 != *v.Fixed64 {
			t.Fatalf("expected pointer tag %#v, got %#v", m.Fixed64, v.Fixed64)
		}
	}
}

func TestTagsSliceUnmarshal(t *testing.T) {
	t.Parallel()

	pb := &testproto.ProtoTagsSlice{
		Sfixed32Slice: protoTagsSliceMarshal.Sfixed32Slice,
		Sfixed64Slice: protoTagsSliceMarshal.Sfixed64Slice,
		Fixed32Slice:  protoTagsSliceMarshal.Fixed32Slice,
		Fixed64Slice:  protoTagsSliceMarshal.Fixed64Slice,
	}
	pbData, err := proto.Marshal(pb)
	if err != nil {
		t.Fatalf("protobuf marshal: %v", err)
	}

	var v = struct {
		Sfixed32Slice []int32  `protobuf:"sfixed32,repeated"`
		Sfixed64Slice []int64  `protobuf:"sfixed64,repeated"`
		Fixed32Slice  []uint32 `protobuf:"fixed32,repeated"`
		Fixed64Slice  []uint64 `protobuf:"fixed64,repeated"`
	}{}
	if err = Unmarshal(pbData, &v); err != nil {
		t.Fatalf("unmarshal bytes: %v", err)
	}

	if !reflect.DeepEqual(pb.Sfixed32Slice, v.Sfixed32Slice) {
		t.Fatalf("execpted slice tag %#v, got %#v", pb.Sfixed32Slice, v.Sfixed32Slice)
	}
	if !reflect.DeepEqual(pb.Sfixed64Slice, v.Sfixed64Slice) {
		t.Fatalf("execpted slice tag %#v, got %#v", pb.Sfixed64Slice, v.Sfixed64Slice)
	}
	if !reflect.DeepEqual(pb.Fixed32Slice, v.Fixed32Slice) {
		t.Fatalf("execpted slice tag %#v, got %#v", pb.Fixed32Slice, v.Fixed32Slice)
	}
	if !reflect.DeepEqual(pb.Fixed64Slice, v.Fixed64Slice) {
		t.Fatalf("execpted slice tag %#v, got %#v", pb.Fixed64Slice, v.Fixed64Slice)
	}
}
