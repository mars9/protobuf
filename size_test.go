package protobuf

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestTypeSize(t *testing.T) {
	t.Parallel()

	for _, v := range typesMessages {
		n := proto.Size(v)

		val := reflect.ValueOf(v)
		m := sizeStruct(val.Elem(), val.Elem().NumField())

		if n != m {
			t.Fatalf("type size: expected size %d, got %d", n, m)
		}
	}
}

func TestSliceSize(t *testing.T) {
	t.Parallel()

	for _, v := range sliceMessages {
		n := proto.Size(v)

		val := reflect.ValueOf(v)
		m := sizeStruct(val.Elem(), val.Elem().NumField())

		if n != m {
			t.Fatalf("slice size: expected size %d, got %d", n, m)
		}
	}
}

func TestStructSize(t *testing.T) {
	t.Parallel()

	for _, v := range structMessages {
		n := proto.Size(v)

		val := reflect.ValueOf(v)
		m := sizeStruct(val.Elem(), val.Elem().NumField())

		if n != m {
			t.Fatalf("struct size: expected size %d, got %d", n, m)
		}
	}
}
