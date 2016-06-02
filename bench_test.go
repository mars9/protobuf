package protobuf

import (
	"bytes"
	"crypto/rand"
	"io"
	"math"
	"testing"

	"github.com/golang/protobuf/proto"
	testproto "github.com/mars9/protobuf/internal/proto"
)

var (
	benchBytes [][]byte
	benchSlice = &testproto.SliceMessage{
		Int32:   []int32{math.MinInt32, math.MaxInt32, math.MinInt32, math.MaxInt32, math.MinInt32, math.MaxInt32},
		Int64:   []int64{math.MinInt64, math.MaxInt64, math.MinInt64, math.MaxInt64, math.MinInt64, math.MaxInt64},
		Uint32:  []uint32{0, math.MaxUint32, 0, math.MaxUint32, 0, math.MaxUint32},
		Uint64:  []uint64{0, math.MaxUint64, 0, math.MaxUint64, 0, math.MaxUint64},
		Float32: []float32{math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32},
		Float64: []float64{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64},
		Bool:    []bool{false, true, false, true, false, true, true, false, false, true, true, false, true, true},
	}
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
}

func BenchmarkProtobuf(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, 0)

	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(benchSlice); err != nil {
			b.Fatalf("protobuf encode: %v", err)
		}
	}
}

func BenchmarkGoogleProtobuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := proto.Marshal(benchSlice); err != nil {
			b.Fatalf("google protobuf marshal: %v", err)
		}
	}
}
