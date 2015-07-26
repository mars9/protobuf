package protobuf

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestReadUvarint(t *testing.T) {
	buf := make([]byte, binary.MaxVarintLen64)

	for _, x := range []uint64{0, 128, 255, 1024, 16 * 1024, 1<<64 - 1} {
		n := binary.PutUvarint(buf, x)
		y, m, err := readUvarint(bytes.NewReader(buf))
		if err != nil {
			t.Fatalf("readUvarint: %v", err)
		}

		if x != y {
			t.Fatalf("readUvarint: expected value %d, got %d", x, y)
		}
		if n != m {
			t.Fatalf("readUvarint: expected length %d, got %d", n, m)
		}
	}
}

func TestWriteUvarint(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	for _, x := range []uint64{0, 128, 255, 1024, 16 * 1024, 1<<64 - 1} {
		buf.Reset()
		if err := writeUvarint(buf, x); err != nil {
			t.Fatalf("writeUvarint: %v", err)
		}

		y, m := binary.Uvarint(buf.Bytes())

		if x != y {
			t.Fatalf("writeUvarint: expected value %d, got %d", x, y)
		}
		if buf.Len() != m {
			t.Fatalf("writeUvarint: expected length %d, got %d", buf.Len(), m)
		}
	}
}
