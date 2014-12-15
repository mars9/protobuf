package protobuf

import (
	"bytes"
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
