package protobuf

import (
	"bytes"
	"math"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mars9/protobuf/internal/testproto"
)

var uintMessages = []struct {
	Uint32 uint32
	Uint64 uint64
}{
	{Uint32: 0, Uint64: 0},
	{Uint32: 1, Uint64: 1},
	{Uint32: math.MaxUint32, Uint64: math.MaxUint64},
}

var intMessages = []struct {
	Int32 int32
	Int64 int64
}{
	{Int32: 0, Int64: 0},
	{Int32: 1, Int64: 1},
	{Int32: -1, Int64: -1},
	{Int32: math.MaxInt32, Int64: math.MaxInt64},
	{Int32: math.MinInt32, Int64: math.MinInt64},
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

var sliceMessages = []struct {
	Uint32Slice []uint32
	Uint64Slice []uint64
	Int32Slice  []int32
	Int64Slice  []int64
	BoolSlice   []bool
	StringSlice []string
}{
	{Uint32Slice: []uint32{0, 0, 0}},
	{Uint32Slice: []uint32{1, 1, 1}},
	{Uint32Slice: []uint32{math.MaxUint32, math.MaxUint32, math.MaxUint32}},

	{Uint64Slice: []uint64{0, 0, 0}},
	{Uint64Slice: []uint64{1, 1, 1}},
	{Uint64Slice: []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64}},

	{Int32Slice: []int32{0, 0, 0}},
	{Int32Slice: []int32{1, 1, 1}},
	{Int32Slice: []int32{-1, -1, -1}},
	{Int32Slice: []int32{math.MaxInt32, math.MaxInt32, math.MaxInt32}},
	{Int32Slice: []int32{math.MinInt32, math.MinInt32, math.MinInt32}},

	{Int64Slice: []int64{0, 0, 0}},
	{Int64Slice: []int64{1, 1, 1}},
	{Int64Slice: []int64{-1, -1, -1}},
	{Int64Slice: []int64{math.MaxInt64, math.MaxInt64, math.MaxInt64}},
	{Int64Slice: []int64{math.MinInt64, math.MinInt64, math.MinInt64}},

	{BoolSlice: []bool{true, false}},

	{StringSlice: []string{"", "", ""}},
	{StringSlice: []string{"string", "string", "string"}},
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

func TestSliceMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range sliceMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size slice: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal slice: %v", err)
		}
		if n != size {
			t.Fatalf("marshal slice: expected size %d, got %d", size, n)
		}

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
		if size != len(pb) {
			t.Fatalf("marshal slice: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal bytes: expected bytes %q, got %q", pb, b)
		}
	}
}
