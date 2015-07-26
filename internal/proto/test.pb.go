// Code generated by protoc-gen-go.
// source: test.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	test.proto

It has these top-level messages:
	TypesMessage
	SliceMessage
	StructMessage
*/
package proto

import proto1 "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal

type TypesMessage struct {
	Uint32  uint32  `protobuf:"varint,1,opt" json:"Uint32,omitempty"`
	Uint64  uint64  `protobuf:"varint,2,opt" json:"Uint64,omitempty"`
	Int32   int32   `protobuf:"varint,3,opt" json:"Int32,omitempty"`
	Int64   int64   `protobuf:"varint,4,opt" json:"Int64,omitempty"`
	Float32 float32 `protobuf:"fixed32,5,opt" json:"Float32,omitempty"`
	Float64 float64 `protobuf:"fixed64,6,opt" json:"Float64,omitempty"`
	Bool    bool    `protobuf:"varint,7,opt" json:"Bool,omitempty"`
	String_ string  `protobuf:"bytes,8,opt" json:"String,omitempty"`
	Bytes   []byte  `protobuf:"bytes,9,opt,proto3" json:"Bytes,omitempty"`
}

func (m *TypesMessage) Reset()         { *m = TypesMessage{} }
func (m *TypesMessage) String() string { return proto1.CompactTextString(m) }
func (*TypesMessage) ProtoMessage()    {}

type SliceMessage struct {
	Uint32  []uint32  `protobuf:"varint,1,rep" json:"Uint32,omitempty"`
	Uint64  []uint64  `protobuf:"varint,2,rep" json:"Uint64,omitempty"`
	Int32   []int32   `protobuf:"varint,3,rep" json:"Int32,omitempty"`
	Int64   []int64   `protobuf:"varint,4,rep" json:"Int64,omitempty"`
	Float32 []float32 `protobuf:"fixed32,5,rep" json:"Float32,omitempty"`
	Float64 []float64 `protobuf:"fixed64,6,rep" json:"Float64,omitempty"`
	Bool    []bool    `protobuf:"varint,7,rep" json:"Bool,omitempty"`
	String_ []string  `protobuf:"bytes,8,rep" json:"String,omitempty"`
	Bytes   [][]byte  `protobuf:"bytes,9,rep,proto3" json:"Bytes,omitempty"`
}

func (m *SliceMessage) Reset()         { *m = SliceMessage{} }
func (m *SliceMessage) String() string { return proto1.CompactTextString(m) }
func (*SliceMessage) ProtoMessage()    {}

type StructMessage struct {
	Types  *TypesMessage `protobuf:"bytes,1,opt,name=types" json:"types,omitempty"`
	Slices *SliceMessage `protobuf:"bytes,2,opt,name=slices" json:"slices,omitempty"`
}

func (m *StructMessage) Reset()         { *m = StructMessage{} }
func (m *StructMessage) String() string { return proto1.CompactTextString(m) }
func (*StructMessage) ProtoMessage()    {}

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

func init() {
}
