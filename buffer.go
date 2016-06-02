package protobuf

import "io"

// Marshal traverses the value v recursively returns the protocol buffer
// encoding of v. The returned slice may be a sub- slice of data if data
// was large enough to hold the entire encoded block. Otherwise, a newly
// allocated slice will be returned.
func Marshal(data []byte, v interface{}) ([]byte, error) {
	b := buffer(data)
	if err := NewEncoder(&b, 0).Encode(v); err != nil {
		return nil, err
	}
	return b, nil
}

// Unmarshal parses the protocol buffer encoded data and stores the
// result in the value pointed to by v. Unmarshal uses the inverse of
// the encodings that Marshal uses
func Unmarshal(data []byte, v interface{}) error {
	b := buffer(data)
	return NewDecoder(&b, 0).Decode(v)
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
