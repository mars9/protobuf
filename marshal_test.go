package protobuf

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mars9/protobuf/internal/testproto"
)

var uintMessages = []struct {
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
}{
	{Uint8: 0, Uint16: 0, Uint32: 0, Uint64: 0},
	{Uint8: 1, Uint16: 1, Uint32: 1, Uint64: 1},
	{Uint8: maxUint8, Uint16: maxUint16, Uint32: maxUint32, Uint64: maxUint64},
}

var intMessages = []struct {
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64
}{
	{Int8: 0, Int16: 0, Int32: 0, Int64: 0},
	{Int8: 1, Int16: 1, Int32: 1, Int64: 1},
	{Int8: -1, Int16: -1, Int32: -1, Int64: -1},
	{Int8: maxInt8, Int16: maxInt16, Int32: maxInt32, Int64: maxInt64},
	{Int8: minInt8, Int16: minInt16, Int32: minInt32, Int64: minInt64},
}

var boolMessages = []struct {
	Bool bool
}{
	{Bool: true},
	{Bool: false},
}

var bytesMessages = []struct {
	String string
	Bytes  []byte
}{
	{String: "", Bytes: []byte("")},
	{String: "string", Bytes: []byte("bytes")},
}

func TestUintMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range uintMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size uint: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal uint: %v", err)
		}
		if n != size {
			t.Fatalf("marshal uint: expected size %d, got %d", size, n)
		}

		pbMsg := testproto.TestUint{
			Uint8:  proto.Uint32(uint32(msg.Uint8)),
			Uint16: proto.Uint32(uint32(msg.Uint16)),
			Uint32: proto.Uint32(msg.Uint32),
			Uint64: proto.Uint64(msg.Uint64),
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		if size != len(pb) {
			t.Fatalf("marshal uint: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal uint: expected bytes %q, got %q", pb, b)
		}
	}
}

func TestIntMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range intMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size int: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal int: %v", err)
		}
		if n != size {
			t.Fatalf("marshal int: expected size %d, got %d", size, n)
		}

		pbMsg := testproto.TestInt{
			Int8:  proto.Int32(int32(msg.Int8)),
			Int16: proto.Int32(int32(msg.Int16)),
			Int32: proto.Int32(msg.Int32),
			Int64: proto.Int64(msg.Int64),
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		if size != len(pb) {
			t.Fatalf("marshal int: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal int: expected bytes %q, got %q", pb, b)
		}
	}
}

func TestBoolMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range boolMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size bool: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal bool: %v", err)
		}
		if n != size {
			t.Fatalf("marshal bool: expected size %d, got %d", size, n)
		}

		pbMsg := testproto.TestBool{
			Bool: proto.Bool(msg.Bool),
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		if size != len(pb) {
			t.Fatalf("marshal bool: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal bool: expected bytes %q, got %q", pb, b)
		}
	}
}

func TestBytesMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range bytesMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size bytes: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal bytes: %v", err)
		}
		if n != size {
			t.Fatalf("marshal bytes: expected size %d, got %d", size, n)
		}
		pbMsg := testproto.TestBytes{
			String_: proto.String(msg.String),
			Bytes:   msg.Bytes,
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		if size != len(pb) {
			t.Fatalf("marshal bytes: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal bytes: expected bytes %q, got %q", pb, b)
		}
	}
}
