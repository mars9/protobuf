package protobuf

import (
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"time"
)

// Reader defines the decode reader. Typically this is a *bufio.Reader.
type Reader interface {
	io.ByteReader
	io.Reader
}

// Decoder manages the receipt of type and data information read from the
// remote side of a connection.
type Decoder struct {
	r   Reader
	max int
}

// NewDecoder returns a new decoder that reads from the io.Reader.
//
// Max defines the maximum size that can be read, if max is 0, the maximum
// message size is not checked.
func NewDecoder(r Reader, max int) *Decoder {
	return &Decoder{r: r, max: max}
}

// Decode first reads the varint encoded message size and then reads
// the next value from the input stream and stores it in the data
// represented by the empty interface value. If v is nil, the value will
// be discarded. Otherwise, the value underlying v must be a pointer to
// the correct type for the next data item received.
func (d *Decoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if !val.IsValid() || val.IsNil() {
		return d.decodeNil()
	}
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}

	size, err := readLength(d.r, d.max)
	if err != nil {
		return err
	}
	data := make([]byte, size)
	if _, err = io.ReadFull(d.r, data); err != nil {
		return err
	}
	return decodeStruct(val.Elem(), data, true)
}

// Reset discards any buffered data, resets all state, and switches the
// decoder reader to read from r.
func (d *Decoder) Reset(r Reader, max int) {
	d.max = max
	d.r = r
}

func (d *Decoder) decodeNil() error {
	size, err := readLength(d.r, d.max)
	if err != nil {
		return err
	}

	msize := int64(size)
	for done := int64(0); ; {
		n, err := io.CopyN(ioutil.Discard, d.r, msize-done)
		if err != nil {
			return err
		}
		done += n
		if done >= msize {
			break
		}
	}
	return nil
}

func decodeStruct(val reflect.Value, data []byte, unsafe bool) error {
	fields, size := val.NumField(), len(data)
	var field reflect.Value
	var err error
	for off := 0; off < size && err == nil; {
		key, n := binary.Uvarint(data[off:])
		if n <= 0 {
			return errors.New("invalid field key")
		}
		off += n

		fnum := int(key >> 3)
		if fnum > 0 && fnum <= fields {
			field = val.Field(fnum - 1)
		} else {
			break
		}

		switch key & 7 {
		case wireVarint:
			v, n := binary.Uvarint(data[off:])
			if n <= 0 {
				return errors.New("bad varint value")
			}
			if err = decodeUvarint(field, v); err != nil {
				return err
			}
			off += n
		case wireFixed32:
			if off+4 >= size {
				return errors.New("bad 32-bit value")
			}
			v := binary.LittleEndian.Uint32(data[off:])
			if err = decodeFixed32(field, v); err != nil {
				return err
			}
			off += 4
		case wireFixed64:
			if off+8 >= size {
				return errors.New("bad 64-bit value")
			}
			v := binary.LittleEndian.Uint64(data[off:])
			if err = decodeFixed64(field, v); err != nil {
				return err
			}
			off += 8
		case wireBytes:
			v, n := binary.Uvarint(data[off:])
			if n <= 0 {
				return errors.New("bad varint size value")
			}
			m := int(v)
			off += n
			err = decodeBytes(field, data[off:off+m], unsafe)
			if err != nil {
				return err
			}
			off += m
		}
	}
	return err
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func decodeError(val reflect.Value, v []byte) (bool, error) {
	if _, ok := val.Interface().(error); ok || val.Type() == errorType {
		val.Set(reflect.ValueOf(errors.New(string(v))))
		return true, nil
	}
	return false, nil
}

func decodeBytes(val reflect.Value, v []byte, unsafe bool) error {
	custom, err := decodeError(val, v)
	if custom {
		return err
	}

	kind := val.Kind()
	switch kind {
	case reflect.Interface:
		return decodeStruct(val.Elem(), v, unsafe)
	case reflect.Struct:
		return decodeStruct(val, v, unsafe)
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.Slice, reflect.Struct, reflect.Ptr:
			elem := reflect.New(val.Type().Elem()).Elem()
			if err = decodeBytes(elem, v, unsafe); err != nil {
				return err
			}
			val.Set(reflect.Append(val, elem))
			return nil
		}
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return decodeBytes(val.Elem(), v, unsafe)
	}

	switch kind {
	case reflect.String:
		val.SetString(string(v))
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.String:
			elem := reflect.New(val.Type().Elem()).Elem()
			elem.SetString(string(v))
			val.Set(reflect.Append(val, elem))
		case reflect.Uint8: // []byte
			if !unsafe {
				val.SetBytes(append([]byte(nil), v...))
			} else {
				val.SetBytes(v)
			}
		}
	}
	return err
}

func decodeUvarint(val reflect.Value, v uint64) error {
	if _, ok := val.Interface().(time.Time); ok {
		ns := int64(v)
		t := time.Unix(ns/int64(time.Second), ns%int64(time.Second))
		val.Set(reflect.ValueOf(t))
		return nil
	}

	switch val.Kind() {
	case reflect.Int32, reflect.Int64:
		return setInt(val, int64(v))
	case reflect.Uint32, reflect.Uint64:
		return setUint(val, v)
	case reflect.Bool:
		return setBool(val, v)
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return decodeUvarint(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := decodeUvarint(elem, v); err != nil {
			return err
		}
		val.Set(reflect.Append(val, elem))
	}
	return nil
}

func decodeFixed64(val reflect.Value, v uint64) error {
	switch val.Kind() {
	//case reflect.Int64:
	//  return setInt(val, int64(v))
	//case reflect.Uint64:
	//  return setUint(val, v)
	case reflect.Float64:
		return setFloat(val, math.Float64frombits(v))
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return decodeFixed64(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := decodeFixed64(elem, v); err != nil {
			return err
		}
		val.Set(reflect.Append(val, elem))
	}
	return nil
}

func decodeFixed32(val reflect.Value, v uint32) error {
	switch val.Kind() {
	//case reflect.Int32:
	//  return setInt(val, int64(int32(v)))
	//case reflect.Uint32:
	//  return setUint(val, uint64(v))
	case reflect.Float32:
		return setFloat(val, float64(math.Float32frombits(v)))
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return decodeFixed32(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := decodeFixed32(elem, v); err != nil {
			return err
		}
		val.Set(reflect.Append(val, elem))
	}
	return nil
}

func setUint(val reflect.Value, v uint64) error {
	if val.OverflowUint(v) {
		return errors.New("uint overflow")
	}
	val.SetUint(v)
	return nil
}

func setInt(val reflect.Value, v int64) error {
	if val.OverflowInt(v) {
		return errors.New("int overflow")
	}
	val.SetInt(v)
	return nil
}

func setBool(val reflect.Value, v uint64) error {
	if v != 0 && v != 1 {
		return errors.New("invalid bool value")
	}
	if v == 0 {
		val.SetBool(false)
	} else {
		val.SetBool(true)
	}
	return nil
}

func setFloat(val reflect.Value, v float64) error {
	if val.OverflowFloat(v) {
		return errors.New("float64 overflow")
	}
	val.SetFloat(v)
	return nil
}
