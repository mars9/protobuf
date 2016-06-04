package protobuf

import (
	"errors"
	"io"
	"reflect"
)

// Marshal traverses the value v recursively and returns the protocol
// buffer encoding of v. The struct underlying v must be a pointer.
//
// Marshal currently encodes all visible field, which does not allow
// distinction between 'required' and 'optional' fields. Marshal ignores
// unsupported struct field types.
//
// The returned slice may be a sub- slice of data if data was large
// enough to hold the entire encoded block. Otherwise, a newly allocated
// slice will be returned.
func Marshal(data []byte, v interface{}) ([]byte, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return data, errors.New("v must be a pointer to a struct")
	}

	val = val.Elem()
	fields := val.NumField()
	b := buffer(data)
	enc := NewEncoder(&b, 0)
	if err := enc.encodeStruct(val, fields); err != nil {
		return nil, err
	}
	return b, nil
}

// Unmarshal parses the protocol buffer representation in data and places
// the decoded result in v. If the struct underlying v does not match the
// data, the results can be unpredictable.
//
// Unmarshal uses the inverse of the encodings that Marshal uses,
// allocating slices and pointers as necessary.
func Unmarshal(data []byte, v interface{}) error {
	val := reflect.ValueOf(v)
	if !val.IsValid() || val.IsNil() {
		return nil
	}
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}

	return decodeStruct(val.Elem(), data, false)
}

// UnmarshalUnsafe parses the protocol buffer representation in data and
// places the decoded result in v. If the struct underlying v does not
// match the data, the results can be unpredictable.
//
// UnmarshalUnsafe uses the inverse of the encodings that Marshal uses,
// allocating slices and pointers as necessary.
//
// UnmarshalUnsafe does not copy raw byte slices. Most code should use
// Unmarshal instead.
func UnmarshalUnsafe(data []byte, v interface{}) error {
	val := reflect.ValueOf(v)
	if !val.IsValid() || val.IsNil() {
		return nil
	}
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}

	return decodeStruct(val.Elem(), data, true)
}

type buffer []byte

func (b *buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *buffer) WriteByte(p byte) error {
	*b = append(*b, p)
	return nil
}

func (b *buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(*b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, *b)
	*b = (*b)[n:]
	return n, nil
}

func (b *buffer) ReadByte() (byte, error) {
	if len(*b) == 0 {
		return 0, io.EOF
	}

	p := (*b)[0]
	*b = (*b)[1:]
	return p, nil
}
