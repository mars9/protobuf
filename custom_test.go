package protobuf

import (
	"reflect"
	"testing"
	"time"
)

type myUint32 uint32

type customStruct struct {
	Duration time.Duration
	MyUint32 myUint32
}

type timeStruct struct {
	Time1 time.Time
	Time2 time.Time
}

type emptyStruct struct{}

func TestCustomTypes(t *testing.T) {
	t.Parallel()

	s := customStruct{
		Duration: time.Duration(42),
		MyUint32: myUint32(42),
	}

	data := make([]byte, Size(&s))
	if _, err := Marshal(data, &s); err != nil {
		t.Fatalf("custom marshal: %v", err)
	}

	v := customStruct{}
	if err := Unmarshal(data, &v); err != nil {
		t.Fatalf("custom unmarshal: %v", err)
	}

	if !reflect.DeepEqual(s, v) {
		t.Fatalf("custom types: expected %#v, got %#v", s, v)
	}
}

func TestTimeType(t *testing.T) {
	t.Parallel()

	s := timeStruct{
		Time1: time.Now(),
		Time2: time.Unix(1, 1),
	}

	size := Size(&s)
	data := make([]byte, size)
	n, err := Marshal(data, &s)
	if err != nil {
		t.Fatalf("time marshal: %v", err)
	}
	if n != size {
		t.Fatalf("time marshal: expected size %d, got %d", size, n)
	}

	v := timeStruct{}
	if err := Unmarshal(data, &v); err != nil {
		t.Fatalf("time unmarshal: %v", err)
	}

	if !reflect.DeepEqual(s, v) {
		t.Fatalf("time types: expected %#v, got %#v", s, v)
	}
}

func TestEmptyStruct(t *testing.T) {
	t.Parallel()

	s := emptyStruct{}

	size := Size(&s)
	data := make([]byte, size)
	n, err := Marshal(data, &s)
	if err != nil {
		t.Fatalf("empty marshal: %v", err)
	}
	if n != size {
		t.Fatalf("empty marshal: expected size %d, got %d", size, n)
	}

	v := emptyStruct{}
	if err := Unmarshal(data, &v); err != nil {
		t.Fatalf("empty unmarshal: %v", err)
	}

	if !reflect.DeepEqual(s, v) {
		t.Fatalf("empty types: expected %#v, got %#v", s, v)
	}
}
