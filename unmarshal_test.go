package protobuf

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mars9/protobuf/internal/testproto"
)

func TestUintUnmarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range uintMessages {
		pbMsg := testproto.TestUint{
			Uint32: proto.Uint32(msg.Uint32),
			Uint64: proto.Uint64(msg.Uint64),
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		var m = struct {
			Uint32 uint32
			Uint64 uint64
		}{}
		if err = Unmarshal(pb, &m); err != nil {
			t.Fatalf("unmarshal uint: %v", err)
		}
		if msg.Uint32 != m.Uint32 {
			t.Fatalf("unmarshal uint: expected uint32 %d, got %d", msg.Uint32, m.Uint32)
		}
		if msg.Uint64 != m.Uint64 {
			t.Fatalf("unmarshal uint: expected uint64 %d, got %d", msg.Uint64, m.Uint64)
		}
	}
}

func TestIntUnmarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range intMessages {
		pbMsg := testproto.TestInt{
			Int32: proto.Int32(msg.Int32),
			Int64: proto.Int64(msg.Int64),
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		var m = struct {
			Int32 int32
			Int64 int64
		}{}
		if err = Unmarshal(pb, &m); err != nil {
			t.Fatalf("unmarshal int: %v", err)
		}
		if msg.Int32 != m.Int32 {
			t.Fatalf("unmarshal int: expected int32 %d, got %d", msg.Int32, m.Int32)
		}
		if msg.Int64 != m.Int64 {
			t.Fatalf("unmarshal int: expected int64 %d, got %d", msg.Int64, m.Int64)
		}
	}
}

func TestBoolUnmarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range boolMessages {
		pbMsg := testproto.TestBool{
			Bool: proto.Bool(msg.Bool),
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		var m = struct {
			Bool bool
		}{}
		if err = Unmarshal(pb, &m); err != nil {
			t.Fatalf("unmarshal bool: %v", err)
		}
		if msg.Bool != m.Bool {
			t.Fatalf("unmarshal bool: expected bool %v, got %v", msg.Bool, m.Bool)
		}
	}
}

func TestBytesUnmarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range bytesMessages {
		pbMsg := testproto.TestBytes{
			String_: proto.String(msg.String),
			Bytes:   msg.Bytes,
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		var m = struct {
			String string
			Bytes  []byte
		}{}
		if err = Unmarshal(pb, &m); err != nil {
			t.Fatalf("unmarshal bytes: %v", err)
		}
		if msg.String != m.String {
			t.Fatalf("unmarshal bytes: expected string %q, got %q", msg.String, m.String)
		}
		if bytes.Compare(msg.Bytes, m.Bytes) != 0 {
			t.Fatalf("unmarshal bytes: expected bytes %q, got %q", msg.Bytes, m.Bytes)
		}
	}
}

func TestSliceUnmarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range sliceMessages {
		pbMsg := testproto.TestSlice{
			Uint32Slice: msg.Uint32Slice,
			Uint64Slice: msg.Uint64Slice,
			Int32Slice:  msg.Int32Slice,
			Int64Slice:  msg.Int64Slice,
			BoolSlice:   msg.BoolSlice,
			StringSlice: msg.StringSlice,
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		var m = struct {
			Uint32Slice []uint32
			Uint64Slice []uint64
			Int32Slice  []int32
			Int64Slice  []int64
			BoolSlice   []bool
			StringSlice []string
		}{}
		if err = Unmarshal(pb, &m); err != nil {
			t.Fatalf("unmarshal slice: %v", err)
		}

		if !reflect.DeepEqual(msg.Uint32Slice, m.Uint32Slice) {
			t.Fatalf("unmarshal uint32 slice: expected %#v, got %#v", msg.Uint32Slice, m.Uint32Slice)
		}
		if !reflect.DeepEqual(msg.Uint64Slice, m.Uint64Slice) {
			t.Fatalf("unmarshal uint64 slice: expected %#v, got %#v", msg.Uint64Slice, m.Uint64Slice)
		}
		if !reflect.DeepEqual(msg.Int32Slice, m.Int32Slice) {
			t.Fatalf("unmarshal int32 slice: expected %#v, got %#v", msg.Int32Slice, m.Int32Slice)
		}
		if !reflect.DeepEqual(msg.Int64Slice, m.Int64Slice) {
			t.Fatalf("unmarshal int64 slice: expected %#v, got %#v", msg.Int64Slice, m.Int64Slice)
		}
		if !reflect.DeepEqual(msg.BoolSlice, m.BoolSlice) {
			t.Fatalf("unmarshal bool slice: expected %#v, got %#v", msg.BoolSlice, m.BoolSlice)
		}
		if !reflect.DeepEqual(msg.StringSlice, m.StringSlice) {
			t.Fatalf("unmarshal string slice: expected %#v, got %#v", msg.StringSlice, m.StringSlice)
		}
	}
}

func TestFixedSliceUnmarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range fixedSliceMessages {
		pbMsg := testproto.TestFixedSlice{
			Float32Slice: msg.Float32Slice,
			Float64Slice: msg.Float64Slice,
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		var m = struct {
			Float32Slice []float32
			Float64Slice []float64
		}{}
		if err = Unmarshal(pb, &m); err != nil {
			t.Fatalf("unmarshal fixed slice: %v", err)
		}

		if !reflect.DeepEqual(msg.Float32Slice, m.Float32Slice) {
			t.Fatalf("unmarshal float32 slice: expected %#v, got %#v", msg.Float32Slice, m.Float32Slice)
		}
		if !reflect.DeepEqual(msg.Float64Slice, m.Float64Slice) {
			t.Fatalf("unmarshal float64 slice: expected %#v, got %#v", msg.Float64Slice, m.Float64Slice)
		}
	}
}

func TestEmbeddedUnmarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range embeddedMessages {
		pbMsg := testproto.TestEmbedded{
			Embedded1: &testproto.TestEmbedded_Embedded{
				Uint32: &msg.Embedded1.Uint32,
				Uint64: &msg.Embedded1.Uint64,
			},
			Embedded2: &testproto.TestEmbedded_Embedded{
				Uint32: &msg.Embedded2.Uint32,
				Uint64: &msg.Embedded2.Uint64,
			},
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		var m = struct {
			Embedded1 embeddedMessage
			Embedded2 *embeddedMessage
		}{}
		if err = Unmarshal(pb, &m); err != nil {
			t.Fatalf("unmarshal embedded: %v", err)
		}
		if msg.Embedded1.Uint32 != m.Embedded1.Uint32 {
			t.Fatalf("unmarshal embedded: expected %#v, got %#v", msg.Embedded1.Uint32, m.Embedded1.Uint32)
		}
		if msg.Embedded1.Uint64 != m.Embedded1.Uint64 {
			t.Fatalf("unmarshal embedded: expected %#v, got %#v", msg.Embedded1.Uint64, m.Embedded1.Uint64)
		}
		if msg.Embedded2.Uint32 != m.Embedded2.Uint32 {
			t.Fatalf("unmarshal embedded: expected %#v, got %#v", msg.Embedded2.Uint32, m.Embedded2.Uint32)
		}
		if msg.Embedded2.Uint64 != m.Embedded2.Uint64 {
			t.Fatalf("unmarshal embedded: expected %#v, got %#v", msg.Embedded2.Uint64, m.Embedded2.Uint64)
		}
	}
}
