package protobuf

import (
	"encoding/hex"
	"fmt"
	"log"
)

func ExampleMarshal() {
	type Embedded struct {
		Field1 string
		Field2 uint64
	}
	type MyStruct struct {
		Field1 uint32
		Field2 *Embedded
	}

	var v = &MyStruct{
		Field1: uint32(42),
		Field2: &Embedded{
			Field1: "field1",
			Field2: uint64(43),
		},
	}

	data := make([]byte, Size(v))
	if _, err := Marshal(data, v); err != nil {
		log.Fatalf("marshal: %v", err)
	}

	fmt.Println(hex.EncodeToString(data))

	// Output:
	// 082a120a0a066669656c6431102b
}

func ExampleUnmarshal() {
	type Embedded struct {
		Field1 string
		Field2 uint64
	}
	type MyStruct struct {
		Field1 uint32
		Field2 *Embedded
	}

	data, err := hex.DecodeString("082a120a0a066669656c6431102b")
	if err != nil {
		log.Fatal(err)
	}

	var v = &MyStruct{}
	if err := Unmarshal(data, v); err != nil {
		log.Fatalf("unmarshal %v", err)
	}

	fmt.Printf("%d %s %d\n", v.Field1, v.Field2.Field1, v.Field2.Field2)

	// Output:
	// 42 field1 43
}

func ExampleMarshal_tag() {
	type Embedded struct {
		Field1 string `protobuf:"bytes,required"`
		Field2 uint64 `protobuf:"varint,optional"`
	}
	type MyStruct struct {
		Field1 uint32    `protobuf:"varint,required"`
		Field2 uint32    `protobuf:"fixed32,required"`
		Field3 *Embedded `protobuf:"bytes:required"`
	}

	var v = &MyStruct{
		Field1: uint32(42),
		Field2: uint32(43),
		Field3: &Embedded{
			Field1: "field1",
			Field2: uint64(44),
		},
	}

	data := make([]byte, Size(v))
	if _, err := Marshal(data, v); err != nil {
		log.Fatalf("marshal: %v", err)
	}

	fmt.Println(hex.EncodeToString(data))

	// Output:
	// 082a152b0000001a0a0a066669656c6431102c
}
