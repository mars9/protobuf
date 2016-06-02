package protobuf

import (
	"reflect"
	"testing"

	testproto "github.com/mars9/protobuf/internal/proto"
)

func TestEncodeDecodeBuffer(t *testing.T) {
	t.Parallel()

	for _, v := range structMessages {
		encData, err := Marshal(nil, v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}


		m := &testproto.StructMessage{}
		if err := Unmarshal(encData, m); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		if !reflect.DeepEqual(v, m) {
			t.Fatalf("buffer enc/dec: expected %#v, got %#v", v, m)
		}
	}
}
