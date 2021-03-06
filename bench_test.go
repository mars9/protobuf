package protobuf

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"io"
	"math"
	"testing"

	"github.com/golang/protobuf/proto"
	testproto "github.com/mars9/protobuf/internal/proto"
)

var (
	benchSlice = &testproto.SliceMessage{
		Int32:   []int32{math.MinInt32, math.MaxInt32, math.MinInt32, math.MaxInt32, math.MinInt32, math.MaxInt32},
		Int64:   []int64{math.MinInt64, math.MaxInt64, math.MinInt64, math.MaxInt64, math.MinInt64, math.MaxInt64},
		Uint32:  []uint32{0, math.MaxUint32, 0, math.MaxUint32, 0, math.MaxUint32},
		Uint64:  []uint64{0, math.MaxUint64, 0, math.MaxUint64, 0, math.MaxUint64},
		Float32: []float32{math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32},
		Float64: []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
		Bool:    []bool{false, true, false, true, false, true, true, false, false, true, true, false, true, true},
	}
	benchSimple = &testproto.TypesMessage{
		Uint32:  math.MaxUint32,
		Uint64:  math.MaxUint64,
		Int32:   math.MaxInt32,
		Int64:   math.MaxInt64,
		Float32: math.MaxFloat32,
		Float64: math.MaxFloat64,
		Bool:    true,
		String_: testString,
		Bytes:   testBytes[:],
	}
	//encodedBuffer = proto.NewBuffer(nil)
)

func init() {
	benchSlice.String_ = make([]string, 32)
	benchSlice.Bytes = make([][]byte, 32)
	for i := 0; i < 32; i++ {
		benchSlice.Bytes[i] = make([]byte, 1024)
		if _, err := io.ReadFull(rand.Reader, benchSlice.Bytes[i]); err != nil {
			panic("out of entropy")
		}
		benchSlice.String_[i] = string(benchSlice.Bytes[i])
	}

	//if err := encodedBuffer.Marshal(benchSlice); err != nil {
	//	panic("marshal benchslice: " + err.Error())
	//}
}

func BenchmarkProtobufStreamEncode(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)

	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(benchSlice); err != nil {
			b.Fatalf("stream encode: %v", err)
		}
	}
}

func BenchmarkProtobufBufferEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(nil, benchSlice); err != nil {
			b.Fatalf("protobuf marshal: %v", err)
		}
	}
}

func BenchmarkGoogleProtobufEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := proto.Marshal(benchSlice); err != nil {
			b.Fatalf("google protobuf marshal: %v", err)
		}
	}
}

func BenchmarkGobEncode(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(benchSlice); err != nil {
			b.Fatalf("gob encode: %v", err)
		}
	}
}

func BenchmarkProtobufStream(b *testing.B) {
	v := &testproto.SliceMessage{}
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)
	dec := NewDecoder(buf, 0)

	for i := 0; i < b.N; i++ {
		if err := enc.Encode(benchSlice); err != nil {
			b.Fatalf("stream encode: %v", err)
		}

		v.Reset()
		if err := dec.Decode(v); err != nil {
			b.Fatalf("stream encode: %v", err)
		}
	}
}

func BenchmarkProtobufBuffer(b *testing.B) {
	v := &testproto.SliceMessage{}
	for i := 0; i < b.N; i++ {
		data, err := Marshal(nil, benchSlice)
		if err != nil {
			b.Fatalf("protobuf marshal: %v", err)
		}

		v.Reset()
		if err := Unmarshal(data, v); err != nil {
			b.Fatalf("protobuf unmarshal: %v", err)
		}
	}
}

func BenchmarkProtobufBufferUnsafe(b *testing.B) {
	v := &testproto.SliceMessage{}
	for i := 0; i < b.N; i++ {
		data, err := Marshal(nil, benchSlice)
		if err != nil {
			b.Fatalf("protobuf marshal: %v", err)
		}

		v.Reset()
		if err := UnmarshalUnsafe(data, v); err != nil {
			b.Fatalf("protobuf unmarshal: %v", err)
		}
	}
}

func BenchmarkGoogleProtobuf(b *testing.B) {
	v := &testproto.SliceMessage{}
	for i := 0; i < b.N; i++ {
		data, err := proto.Marshal(benchSlice)
		if err != nil {
			b.Fatalf("google protobuf marshal: %v", err)
		}

		v.Reset()
		if err := proto.Unmarshal(data, v); err != nil {
			b.Fatalf("google protobuf unmarshal: %v", err)
		}
	}
}

func BenchmarkGob(b *testing.B) {
	v := &testproto.SliceMessage{}
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	dec := gob.NewDecoder(buf)

	for i := 0; i < b.N; i++ {
		if err := enc.Encode(benchSlice); err != nil {
			b.Fatalf("stream encode: %v", err)
		}

		v.Reset()
		if err := dec.Decode(v); err != nil {
			b.Fatalf("stream encode: %v", err)
		}
	}
}

func BenchmarkProfile(b *testing.B) {
	v := &testproto.TypesMessage{}
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)
	dec := NewDecoder(buf, 0)

	for i := 0; i < b.N; i++ {
		if err := enc.Encode(benchSimple); err != nil {
			b.Fatalf("stream encode: %v", err)
		}

		v.Reset()
		if err := dec.Decode(v); err != nil {
			b.Fatalf("stream encode: %v", err)
		}
	}
}
