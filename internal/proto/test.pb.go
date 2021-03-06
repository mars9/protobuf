// Code generated by protoc-gen-go.
// source: internal/proto/test.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	internal/proto/test.proto

It has these top-level messages:
	TypesMessage
	SliceMessage
	StructMessage
	NestedStruct
	Struct
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto1.ProtoPackageIsVersion1

type TypesMessage struct {
	Uint32  uint32  `protobuf:"varint,1,opt,name=Uint32,json=uint32" json:"Uint32,omitempty"`
	Uint64  uint64  `protobuf:"varint,2,opt,name=Uint64,json=uint64" json:"Uint64,omitempty"`
	Int32   int32   `protobuf:"varint,3,opt,name=Int32,json=int32" json:"Int32,omitempty"`
	Int64   int64   `protobuf:"varint,4,opt,name=Int64,json=int64" json:"Int64,omitempty"`
	Float32 float32 `protobuf:"fixed32,5,opt,name=Float32,json=float32" json:"Float32,omitempty"`
	Float64 float64 `protobuf:"fixed64,6,opt,name=Float64,json=float64" json:"Float64,omitempty"`
	Bool    bool    `protobuf:"varint,7,opt,name=Bool,json=bool" json:"Bool,omitempty"`
	String_ string  `protobuf:"bytes,8,opt,name=String,json=string" json:"String,omitempty"`
	Bytes   []byte  `protobuf:"bytes,9,opt,name=Bytes,json=bytes,proto3" json:"Bytes,omitempty"`
}

func (m *TypesMessage) Reset()                    { *m = TypesMessage{} }
func (m *TypesMessage) String() string            { return proto1.CompactTextString(m) }
func (*TypesMessage) ProtoMessage()               {}
func (*TypesMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SliceMessage struct {
	Uint32  []uint32  `protobuf:"varint,1,rep,name=Uint32,json=uint32" json:"Uint32,omitempty"`
	Uint64  []uint64  `protobuf:"varint,2,rep,name=Uint64,json=uint64" json:"Uint64,omitempty"`
	Int32   []int32   `protobuf:"varint,3,rep,name=Int32,json=int32" json:"Int32,omitempty"`
	Int64   []int64   `protobuf:"varint,4,rep,name=Int64,json=int64" json:"Int64,omitempty"`
	Float32 []float32 `protobuf:"fixed32,5,rep,name=Float32,json=float32" json:"Float32,omitempty"`
	Float64 []float64 `protobuf:"fixed64,6,rep,name=Float64,json=float64" json:"Float64,omitempty"`
	Bool    []bool    `protobuf:"varint,7,rep,name=Bool,json=bool" json:"Bool,omitempty"`
	String_ []string  `protobuf:"bytes,8,rep,name=String,json=string" json:"String,omitempty"`
	Bytes   [][]byte  `protobuf:"bytes,9,rep,name=Bytes,json=bytes,proto3" json:"Bytes,omitempty"`
}

func (m *SliceMessage) Reset()                    { *m = SliceMessage{} }
func (m *SliceMessage) String() string            { return proto1.CompactTextString(m) }
func (*SliceMessage) ProtoMessage()               {}
func (*SliceMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type StructMessage struct {
	Types  *TypesMessage `protobuf:"bytes,1,opt,name=types" json:"types,omitempty"`
	Slices *SliceMessage `protobuf:"bytes,2,opt,name=slices" json:"slices,omitempty"`
}

func (m *StructMessage) Reset()                    { *m = StructMessage{} }
func (m *StructMessage) String() string            { return proto1.CompactTextString(m) }
func (*StructMessage) ProtoMessage()               {}
func (*StructMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StructMessage) GetTypes() *TypesMessage {
	if m != nil {
		return m.Types
	}
	return nil
}

func (m *StructMessage) GetSlices() *SliceMessage {
	if m != nil {
		return m.Slices
	}
	return nil
}

type NestedStruct struct {
	Arg int32 `protobuf:"varint,1,opt,name=arg" json:"arg,omitempty"`
}

func (m *NestedStruct) Reset()                    { *m = NestedStruct{} }
func (m *NestedStruct) String() string            { return proto1.CompactTextString(m) }
func (*NestedStruct) ProtoMessage()               {}
func (*NestedStruct) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type Struct struct {
	NestedStruct *NestedStruct `protobuf:"bytes,1,opt,name=NestedStruct,json=nestedStruct" json:"NestedStruct,omitempty"`
	Arg          int32         `protobuf:"varint,2,opt,name=arg" json:"arg,omitempty"`
}

func (m *Struct) Reset()                    { *m = Struct{} }
func (m *Struct) String() string            { return proto1.CompactTextString(m) }
func (*Struct) ProtoMessage()               {}
func (*Struct) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Struct) GetNestedStruct() *NestedStruct {
	if m != nil {
		return m.NestedStruct
	}
	return nil
}

func init() {
	proto1.RegisterType((*TypesMessage)(nil), "proto.TypesMessage")
	proto1.RegisterType((*SliceMessage)(nil), "proto.SliceMessage")
	proto1.RegisterType((*StructMessage)(nil), "proto.StructMessage")
	proto1.RegisterType((*NestedStruct)(nil), "proto.NestedStruct")
	proto1.RegisterType((*Struct)(nil), "proto.Struct")
}

var fileDescriptor0 = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x4e, 0x02, 0x31,
	0x10, 0xc6, 0xb3, 0xcc, 0xb6, 0x40, 0x5d, 0x12, 0x53, 0x8d, 0xa9, 0x37, 0xc2, 0x49, 0x63, 0x02,
	0x09, 0x10, 0xbc, 0x73, 0x30, 0xf1, 0xa0, 0x87, 0xa2, 0x0f, 0xb0, 0x60, 0xdd, 0x90, 0x6c, 0x76,
	0xc9, 0xb6, 0x1c, 0x78, 0x5b, 0x1f, 0xc5, 0xe9, 0x94, 0x7f, 0x9a, 0x72, 0xda, 0xce, 0x37, 0xdf,
	0x74, 0xe6, 0xd7, 0x59, 0x71, 0xbf, 0xae, 0x9c, 0x69, 0xaa, 0xbc, 0x1c, 0x6d, 0x9a, 0xda, 0xd5,
	0x23, 0x67, 0xac, 0x1b, 0xd2, 0x51, 0x32, 0xfa, 0x0c, 0x7e, 0x12, 0x91, 0x7d, 0xec, 0x36, 0xc6,
	0xbe, 0x19, 0x6b, 0xf3, 0xc2, 0xc8, 0x3b, 0xc1, 0x3f, 0xb1, 0x6a, 0x32, 0x56, 0x49, 0x3f, 0x79,
	0xe8, 0x69, 0xbe, 0xa5, 0xe8, 0xa0, 0xcf, 0xa6, 0xaa, 0x85, 0x7a, 0x1a, 0xf4, 0xd9, 0x54, 0xde,
	0x0a, 0xf6, 0x4a, 0x76, 0x40, 0x99, 0x69, 0x16, 0xdc, 0x41, 0x45, 0x73, 0x8a, 0x2a, 0x90, 0x8a,
	0x5e, 0x25, 0xda, 0x2f, 0x65, 0x9d, 0x7b, 0x37, 0x43, 0xbd, 0xa5, 0xdb, 0xdf, 0x21, 0x3c, 0x66,
	0xb0, 0x82, 0x63, 0x26, 0xd9, 0x67, 0xb0, 0x46, 0x8a, 0x74, 0x5e, 0xd7, 0xa5, 0x6a, 0xa3, 0xdc,
	0xd1, 0xe9, 0x12, 0xcf, 0x7e, 0x96, 0x85, 0x6b, 0xd6, 0x55, 0xa1, 0x3a, 0xa8, 0x76, 0x35, 0xb7,
	0x14, 0xf9, 0xae, 0xf3, 0x1d, 0x32, 0xaa, 0x2e, 0xca, 0x99, 0x66, 0x4b, 0x1f, 0x10, 0xe2, 0xa2,
	0x5c, 0xaf, 0x4c, 0x0c, 0x11, 0x2e, 0x20, 0x42, 0x1c, 0x11, 0xa2, 0x88, 0x70, 0x01, 0x11, 0x2e,
	0x22, 0x42, 0x1c, 0x11, 0xa2, 0x88, 0x10, 0x47, 0x84, 0x13, 0x62, 0x21, 0x7a, 0xe8, 0xde, 0xae,
	0xdc, 0x01, 0xf1, 0x51, 0x30, 0xe7, 0xb7, 0x4a, 0x4b, 0xbc, 0x1a, 0xdf, 0x84, 0xa5, 0x0f, 0xcf,
	0x37, 0xad, 0x83, 0x43, 0x3e, 0x09, 0x6e, 0xfd, 0xeb, 0x58, 0x5a, 0xec, 0xc9, 0x7b, 0xfe, 0x64,
	0x7a, 0x6f, 0x19, 0xf4, 0x45, 0xf6, 0x8e, 0xff, 0x90, 0xf9, 0x0a, 0xed, 0xe4, 0xb5, 0x80, 0xbc,
	0x29, 0xa8, 0x0b, 0xd3, 0xfe, 0x38, 0x58, 0xd0, 0xe0, 0x3e, 0xf7, 0xfc, 0xd7, 0xfb, 0x6f, 0x94,
	0xf3, 0x94, 0xce, 0xaa, 0xc8, 0xa5, 0xad, 0xe3, 0xa5, 0x4b, 0x4e, 0x35, 0x93, 0xdf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x5f, 0x91, 0xe9, 0x14, 0xd0, 0x02, 0x00, 0x00,
}
