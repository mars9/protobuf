package protobuf

import (
	"reflect"
	"testing"
)

type Message struct {
	Type  int32
	Tag   uint32
	ID    uint32
	Path  []string
	Name  string
	Mode  uint32
	Error string
}

func TestMarshalUnmarshalMessage(t *testing.T) {
	m := &Message{
		Type: 101,
		Tag:  32000,
		ID:   100,
		Path: []string{"a", "b", "c", "d"},
		Name: "file",
		Mode: 2,
	}
	size := Size(m)
	data := make([]byte, size)
	n, err := Marshal(data, m)
	if err != nil {
		t.Fatalf("integration test marshal: %v", err)
	}
	if size != n {
		t.Fatalf("integration test: expected marshal size %d, got %d", size, n)
	}

	r := &Message{}
	if err = Unmarshal(data, r); err != nil {
		t.Fatalf("integration test unmarshal: %v", err)
	}
	if !reflect.DeepEqual(&m, &r) {
		t.Fatalf("integration test: marshal/unmarshal differ: expected %#v, got %#v", m, r)
	}
}

func TestMarshalUnmarshalMessage1(t *testing.T) {
	const WalkMessage int32 = 101

	m := &Message{
		Type: WalkMessage,
		Tag:  32000,
		ID:   100,
		Path: []string{"a", "b", "c", "d"},
	}
	size := Size(m)
	data := make([]byte, size)
	n, err := Marshal(data, m)
	if err != nil {
		t.Fatalf("integration test marshal: %v", err)
	}
	if size != n {
		t.Fatalf("integration test: expected marshal size %d, got %d", size, n)
	}

	r := &Message{}
	if err = Unmarshal(data, r); err != nil {
		t.Fatalf("integration test unmarshal: %v", err)
	}
	if !reflect.DeepEqual(&m, &r) {
		t.Fatalf("integration test: marshal/unmarshal differ: expected %#v, got %#v", m, r)
	}
}
