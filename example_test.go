package protobuf

import (
	"encoding/hex"
	"fmt"
	"log"
)

func Example_MarshalUnmarshal() {
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
		log.Fatal("marshal: %v", err)
	}

	fmt.Println(hex.EncodeToString(data))

	var x = &MyStruct{}
	if err := Unmarshal(data, x); err != nil {
		log.Fatalf("unmarshal %v", err)
	}

	fmt.Printf("%d %s %d\n", x.Field1, x.Field2.Field1, x.Field2.Field2)

	// Output:
	// 082a120a0a066669656c6431102b
	// 42 field1 43
}

func Example_TagMarshalUnmarshal() {
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
		log.Fatal("marshal: %v", err)
	}

	fmt.Println(hex.EncodeToString(data))

	var x = &MyStruct{}
	if err := Unmarshal(data, x); err != nil {
		log.Fatalf("unmarshal %v", err)
	}

	fmt.Printf("%d %d %s %d\n", x.Field1, x.Field2, x.Field3.Field1, x.Field3.Field2)

	// Output:
	// 082a152b0000001a0a0a066669656c6431102c
	// 42 43 field1 44
}
