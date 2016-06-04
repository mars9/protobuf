package protobuf

import (
	"errors"
	"io"
	"sync"
)

const maxInt = uint64(^uint(0) >> 1)

// readLength reads an varint encoded integer from an io.ByteWriter. Max
// defines the maximum message size that can be read, if max is 0, the
// maximum message size is not checked.
func readLength(r io.ByteReader, max int) (int, error) {
	size, _, err := readUvarint(r)
	if err != nil {
		return 0, err
	}
	if size > maxInt {
		return 0, errors.New("integer overflow")
	}

	msize := int(size)
	if max > 0 && msize > max {
		return 0, errors.New("message too large")
	}
	return msize, nil
}

// writeLength encodes size and writes it to an io.ByteWriter. Max defines
// the maximum message size that can be read, if max is 0, the maximum
// message size is not checked.
func writeLength(w io.ByteWriter, size, max int) error {
	if max > 0 && size > max {
		return errors.New("message too large")
	}
	return writeUvarint(w, uint64(size))
}

func readUvarint(r io.ByteReader) (uint64, int, error) {
	var v uint64
	for shift, n := uint(0), 0; ; shift += 7 {
		b, err := r.ReadByte()
		if err != nil {
			return v, n, err
		}
		n++

		if b < 0x80 {
			if n > 10 || n == 10 && b > 1 {
				return v, n, errors.New("64-bit integer overflow")
			}
			return v | uint64(b)<<shift, n, nil
		}
		v |= uint64(b&0x7f) << shift
	}
}

func writeUvarint(w io.ByteWriter, v uint64) (err error) {
	for v >= 0x80 {
		if err = w.WriteByte(byte(v) | 0x80); err != nil {
			return err
		}
		v >>= 7
	}
	return w.WriteByte(byte(v))
}

var (
	pool64 = sync.Pool{New: func() interface{} { return [8]byte{} }}
	pool32 = sync.Pool{New: func() interface{} { return [4]byte{} }}
)

func readFixed32(r io.Reader) (uint32, error) {
	b := pool32.Get().([4]byte)
	defer pool32.Put(b)

	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}

	return uint32(b[0]) |
		uint32(b[1])<<8 |
		uint32(b[2])<<16 |
		uint32(b[3])<<24, nil
}

func writeFixed32(w io.Writer, v uint32) error {
	b := pool32.Get().([4]byte)
	defer pool32.Put(b)

	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)

	_, err := w.Write(b[:])
	return err
}

func readFixed64(r io.Reader) (uint64, error) {
	b := pool64.Get().([8]byte)
	defer pool64.Put(b)

	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}

	return uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56, nil
}

func writeFixed64(w io.Writer, v uint64) error {
	b := pool64.Get().([8]byte)
	defer pool64.Put(b)

	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)

	_, err := w.Write(b[:])
	return err
}
