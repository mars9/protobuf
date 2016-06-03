package protobuf

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	testproto "github.com/mars9/protobuf/internal/proto"
)

func TestEncodeDecodeBuffer(t *testing.T) {
	t.Parallel()

	for _, v := range structMessages {
		xdata, err := Marshal(nil, v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		pdata, err := proto.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if bytes.Compare(xdata, pdata) != 0 {
			t.Fatalf("marshal buffer:\nexpected bytes %q\ngot bytes     %q", pdata, xdata)
		}

		xm, pm := &testproto.StructMessage{}, &testproto.StructMessage{}
		if err := Unmarshal(pdata, xm); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if err := proto.Unmarshal(xdata, pm); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		if !reflect.DeepEqual(xm, pm) {
			t.Fatalf("buffer enc/dec:\nexpected %#v\ngot      %#v", pm, xm)
		}
	}
}

type NestedStruct struct {
	Arg int32
}

type xstruct struct {
	NestedStruct
	Arg int32
}

func TestNestedStruct(t *testing.T) {
	t.Parallel()

	pv := &testproto.Struct{
		NestedStruct: &testproto.NestedStruct{Arg: 42},
		Arg:          42,
	}
	xv := xstruct{
		NestedStruct: NestedStruct{Arg: 42},
		Arg:          42,
	}

	pdata, perr := proto.Marshal(pv)
	xdata, xerr := Marshal(nil, &xv)
	if perr != nil || xerr != nil {
		t.Fatalf("nestest struct test: %v, %v", perr, xerr)
	}
	if bytes.Compare(xdata, pdata) != 0 {
		t.Fatalf("marshal buffer:\nexpected bytes %q\ngot bytes     %q", pdata, xdata)
	}

	yv := xstruct{}
	if err := Unmarshal(pdata, &yv); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(xv, yv) {
		t.Fatalf("buffer enc/dec:\nexpected %#v\ngot      %#v", xv, yv)
	}
}
