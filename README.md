# Protobuf V3 [![GoDoc](https://godoc.org/github.com/mars9/protobuf?status.svg)](https://godoc.org/github.com/mars9/protobuf)

Package protobuf converts data structures to and from the wire format
of protocol buffers. See
https://developers.google.com/protocol-buffers/docs/encoding for
protocol buffers documentation.

This package does not require users to write or compile .proto files,
but does not support all of the Protocol Buffer types.

The following table summarizes the correspondence between .proto
definition types and Go field types:

Go          | Protocol Buffer
----------- | -------------
bool        | optional bool
uint64      | optional uint64
uint32      | optional uint64
int64       | optional int64
int32       | optional int32
float       | optional float32
double      | optional float64
time.Time   | optional int64
[]byte      | optional bytes
string      | optional string
struct      | optional message
[]bool      | repeated bool
[]uint64    | repeated uint64
[]uint32    | repeated uint64
[]int64     | repeated int64
[]int32     | repeated int32
[]float     | repeated float32
[]double    | repeated float64
[][]byte    | repeated bytes
[]string    | repeated string
[]struct    | repeated message

