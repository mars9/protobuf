package protobuf

import (
	"crypto/rand"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	testproto "github.com/mars9/protobuf/internal/proto"
)

var (
	testBytes  [8192]byte
	testString string
)

func init() {
	if _, err := io.ReadFull(rand.Reader, testBytes[:]); err != nil {
		panic("not enough entropy")
	}
	testString = string(testBytes[:])
}

var (
	typesMessages = []*testproto.TypesMessage{
		&testproto.TypesMessage{},
		&testproto.TypesMessage{
			Uint32:  0,
			Uint64:  0,
			Int32:   math.MinInt32,
			Int64:   math.MinInt64,
			Float32: math.SmallestNonzeroFloat32,
			Float64: math.SmallestNonzeroFloat64,
			Bool:    false,
			String_: "abc",
			Bytes:   []byte("abc"),
		},
		&testproto.TypesMessage{
			Uint32:  math.MaxUint32,
			Uint64:  math.MaxUint64,
			Int32:   math.MaxInt32,
			Int64:   math.MaxInt64,
			Float32: math.MaxFloat32,
			Float64: math.MaxFloat64,
			Bool:    true,
			String_: testString,
			Bytes:   testBytes[:],
		},
	}

	sliceMessages = []*testproto.SliceMessage{
		&testproto.SliceMessage{},
		&testproto.SliceMessage{
			Uint32:  []uint32{0, 0, 0},
			Uint64:  []uint64{0, 0, 0},
			Int32:   []int32{math.MinInt32, math.MinInt32, math.MinInt32},
			Int64:   []int64{math.MinInt64, math.MinInt64, math.MinInt64},
			Float32: []float32{math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32},
			Float64: []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
			Bool:    []bool{false, true, false},
			String_: []string{"abc", "def", "ghi"},
			Bytes:   [][]byte{[]byte("abc"), []byte("def"), []byte("ghi")},
		},
		&testproto.SliceMessage{
			Uint32:  []uint32{math.MaxUint32, math.MaxUint32, math.MaxUint32},
			Uint64:  []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64},
			Int32:   []int32{math.MaxInt32, math.MaxInt32, math.MaxInt32},
			Int64:   []int64{math.MaxInt64, math.MaxInt64, math.MaxInt64},
			Float32: []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32},
			Float64: []float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64},
			Bool:    []bool{false, true, false},
			String_: []string{testString, testString, testString},
			Bytes:   [][]byte{testBytes[:], testBytes[:], testBytes[:]},
		},
	}

	structMessages = []*testproto.StructMessage{
		&testproto.StructMessage{},
		&testproto.StructMessage{Types: typesMessages[1]},
		&testproto.StructMessage{Types: typesMessages[2]},
		&testproto.StructMessage{Slices: sliceMessages[0]},
		&testproto.StructMessage{Slices: sliceMessages[1]},
		&testproto.StructMessage{Slices: sliceMessages[2]},
		&testproto.StructMessage{Types: typesMessages[0], Slices: sliceMessages[0]},
		&testproto.StructMessage{Types: typesMessages[1], Slices: sliceMessages[1]},
		&testproto.StructMessage{Types: typesMessages[2], Slices: sliceMessages[2]},
	}
)

func TestTypeDecode(t *testing.T) {
	t.Parallel()

	for _, v := range typesMessages {
		data, err := proto.Marshal(v)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		m := &testproto.TypesMessage{}
		val := reflect.ValueOf(m)
		if err = decodeStruct(val.Elem(), data, true); err != nil {
			t.Fatalf("decode type: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("decode type: expected %#v, got %#v", v, m)
		}
	}
}

func TestSliceDecode(t *testing.T) {
	t.Parallel()

	for _, v := range sliceMessages {
		data, err := proto.Marshal(v)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		m := &testproto.SliceMessage{}
		val := reflect.ValueOf(m)
		if err = decodeStruct(val.Elem(), data, true); err != nil {
			t.Fatalf("decode slice: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("decode slice: expected %#v, got %#v", v, m)
		}
	}
}

func TestStrucDecode(t *testing.T) {
	t.Parallel()

	for _, v := range structMessages {
		data, err := proto.Marshal(v)
		if err != nil {
			t.Fatalf("marshal protobuf: %v", err)
		}

		m := &testproto.StructMessage{}
		val := reflect.ValueOf(m)
		if err = decodeStruct(val.Elem(), data, true); err != nil {
			t.Fatalf("decode struct: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("decode struct: expected %#v, got %#v", v, m)
		}
	}
}
