package protobuf

import (
	"errors"
	"io"
	"reflect"
)

// DecodeBuffer is a buffer for decoding Protocol Buffers. It may be reused
// between invocations to save alloctions.
type DecodeBuffer struct {
	dec *Decoder
	b   *buffer
}

// NewDecodeBuffer creates and initializes a new DecodeBuffer.
func NewDecodeBuffer() *DecodeBuffer {
	b := &buffer{}
	return &DecodeBuffer{
		dec: NewDecoder(b, 0),
		b:   b,
	}
}

// Write sets the contents of p to the buffer. The return value is the
// length of p and err is always nil.
func (b *DecodeBuffer) Write(p []byte) (int, error) {
	*b.b = p
	return len(p), nil
}

// Decode parses the protocol buffer representation in data and places the
// decoded result in v. If the struct underlying v does not match the data,
// the results can be unpredictable.
//
// Decode uses the inverse of the encodings that Marshal uses, allocating
// slices and pointers as necessary.
func (b *DecodeBuffer) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if !val.IsValid() || val.IsNil() {
		return nil
	}
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}

	val = val.Elem()
	return b.dec.decodeStruct(val, val.NumField(), len(*b.b))
}

// EncodeBuffer is a buffer for encoding Protocol Buffers. It may be reused
// between invocations to save alloctions.
type EncodeBuffer struct {
	enc *Encoder
	b   *buffer
}

// NewEncodeBuffer creates and initializes a new EncodeBuffer.
func NewEncodeBuffer() *EncodeBuffer {
	b := &buffer{}
	return &EncodeBuffer{
		enc: NewEncoder(b, 0),
		b:   b,
	}
}

// Encode traverses the value v recursively and returns the protocol buffer
// encoding of v. The struct underlying v must be a pointer.
//
// Encode currently encodes all visible field, which does not allow
// distinction between 'required' and 'optional' fields. Encode ignores
// unsupported struct field types.
func (b *EncodeBuffer) Encode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}

	val = val.Elem()
	return b.enc.encodeStruct(val, val.NumField())
}

// Bytes returns a slice of the contents of the buffer.
func (b *EncodeBuffer) Bytes() []byte { return *b.b }

// Len returns the number of bytes of the slice.
func (b *EncodeBuffer) Len() int { return len(*b.b) }

// Reset resets the buffer so it has no content.
func (b *EncodeBuffer) Reset() { *b.b = (*b.b)[0:0] }

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
