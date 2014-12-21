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

var fixedMessages = []struct {
	Float32 float32
	Float64 float64
}{
	{Float32: 0, Float64: 0},
	{Float32: 1, Float64: 1},
	{Float32: -1, Float64: -1},
	{Float32: math.MaxFloat32, Float64: math.MaxFloat64},
	{Float32: math.SmallestNonzeroFloat32, Float64: math.SmallestNonzeroFloat64},
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

var byteSliceMessages = []struct {
	ByteSlice [][]byte
}{
	{ByteSlice: [][]byte{[]byte(""), []byte(""), []byte("")}},
	{ByteSlice: [][]byte{[]byte("bytes"), []byte("bytes"), []byte("bytes")}},
}

var fixedSliceMessages = []struct {
	Float32Slice []float32
	Float64Slice []float64
}{
	{Float32Slice: []float32{0, 0, 0}},
	{Float32Slice: []float32{1, 1, 1}},
	{Float32Slice: []float32{-1, -1, -1}},
	{Float32Slice: []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}},
	{Float32Slice: []float32{math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32}},

	{Float64Slice: []float64{0, 0, 0}},
	{Float64Slice: []float64{1, 1, 1}},
	{Float64Slice: []float64{-1, -1, -1}},
	{Float64Slice: []float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}},
	{Float64Slice: []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64}},
}

type embeddedMessage struct {
	Uint32 uint32
	Uint64 uint64
}

var embeddedMessages = []struct {
	Embedded1 embeddedMessage
	Embedded2 *embeddedMessage
}{
	{
		Embedded1: embeddedMessage{Uint32: math.MaxUint32, Uint64: math.MaxUint64},
		Embedded2: &embeddedMessage{Uint32: math.MaxUint32, Uint64: math.MaxUint64},
	},
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

func TestFixedMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range fixedMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size fixed: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal fixed: %v", err)
		}
		if n != size {
			t.Fatalf("marshal fixed: expected size %d, got %d", size, n)
		}

		pbMsg := testproto.TestFixed{
			Float32: proto.Float32(msg.Float32),
			Float64: proto.Float64(msg.Float64),
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		if size != len(pb) {
			t.Fatalf("marshal fixed: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal fixed: expected bytes %q, got %q", pb, b)
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

func TestByteSliceMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range byteSliceMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size bytes slice: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal bytes slice: %v", err)
		}
		if n != size {
			t.Fatalf("marshal bytes slice: expected size %d, got %d", size, n)
		}

		pbMsg := testproto.TestByteSlice{
			ByteSlice: msg.ByteSlice,
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}
		if size != len(pb) {
			t.Fatalf("marshal bytes slice: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal bytes slice: expected bytes %q, got %q", pb, b)
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
			t.Fatalf("marshal slice: expected bytes %q, got %q", pb, b)
		}
	}
}

func TestFixedSliceMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range fixedSliceMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size fixed slice: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal fixed slice: %v", err)
		}
		if n != size {
			t.Fatalf("marshal fixed slice: expected size %d, got %d", size, n)
		}

		pbMsg := testproto.TestFixedSlice{
			Float32Slice: msg.Float32Slice,
			Float64Slice: msg.Float64Slice,
		}
		pb, err := proto.Marshal(&pbMsg)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}
		if size != len(pb) {
			t.Fatalf("marshal fixed slice: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal fixed slice: expected bytes %q, got %q", pb, b)
		}
	}
}

func TestEmbeddedMarshal(t *testing.T) {
	t.Parallel()

	for _, msg := range embeddedMessages {
		size, err := Size(&msg)
		if err != nil {
			t.Fatalf("size embedded: %v", err)
		}
		b := make([]byte, size)
		n, err := Marshal(b, &msg)
		if err != nil {
			t.Fatalf("marshal embedded: %v", err)
		}
		if n != size {
			t.Fatalf("marshal embedded: expected size %d, got %d", size, n)
		}

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
		if size != len(pb) {
			t.Fatalf("marshal embedded: expected size %d, got %d", len(pb), size)
		}
		if bytes.Compare(b, pb) != 0 {
			t.Fatalf("marshal embedded: expected bytes %q, got %q", pb, b)
		}
	}
}
