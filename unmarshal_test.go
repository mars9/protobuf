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
