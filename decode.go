package protobuf

import (
	"errors"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"time"
)

// Reader defines the decode reader. Typicall this is a *bufio.Reader.
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

	size, err := ReadLength(d.r, d.max)
	if err != nil {
		return err
	}
	return d.decodeStruct(val.Elem(), val.Elem().NumField(), size)
}

func (d *Decoder) decodeNil() error {
	size, err := ReadLength(d.r, d.max)
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

func (d *Decoder) decodeStruct(val reflect.Value, fields, size int) error {
	var field reflect.Value
	var n, fnum int
	var v uint64
	var err error
	for off := 0; off < size && err == nil; {
		if v, n, err = readUvarint(d.r); err != nil {
			return err
		}
		off += n

		fnum = int(v >> 3)
		if fnum > 0 && fnum <= fields {
			field = val.Field(fnum - 1)
		} else {
			break
		}

		switch v & 7 {
		case wireVarint:
			if v, n, err = readUvarint(d.r); err != nil {
				return err
			}
			off += n
			err = d.decodeUvarint(field, v)
		case wireFixed32:
			v1, err := readFixed32(d.r)
			if err != nil {
				return err
			}
			off += 4
			err = d.decodeFixed32(field, v1)
		case wireFixed64:
			if v, err = readFixed64(d.r); err != nil {
				return err
			}
			off += 8
			err = d.decodeFixed64(field, v)
		case wireBytes:
			if v, n, err = readUvarint(d.r); err != nil {
				return err
			}
			off += n
			n, err = d.decodeBytes(field, int(v))
			off += n
		}
	}
	return err
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

/*
func (d *Decoder) decodeError(val reflect.Value, v int) (int, bool, error) {
	if _, ok := val.Interface().(error); ok || val.Type() == errorType {
		b := make([]byte, v)
		n, err := io.ReadFull(d.r, b)
		if err != nil {
			return n, true, err
		}
		val.Set(reflect.ValueOf(errors.New(string(b))))
		return n, true, nil
	}
	return 0, false, nil
}
*/

func (d *Decoder) decodeBytes(val reflect.Value, v int) (int, error) {
	/*
		n, custom, err := d.decodeError(val, v)
		if custom {
			return n, err
		}
	*/

	kind := val.Kind()
	switch kind {
	case reflect.Interface:
		return v, d.decodeStruct(val.Elem(), val.Elem().NumField(), v)
	case reflect.Struct:
		return v, d.decodeStruct(val, val.NumField(), v)
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.Slice, reflect.Struct, reflect.Ptr:
			elem := reflect.New(val.Type().Elem()).Elem()
			n, err := d.decodeBytes(elem, v)
			if err != nil {
				return n, err
			}
			val.Set(reflect.Append(val, elem))
			return n, err
		}
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return d.decodeBytes(val.Elem(), v)
	}

	b := make([]byte, v)
	n, err := io.ReadFull(d.r, b)
	if err != nil {
		return n, err
	}

	switch kind {
	case reflect.String:
		val.SetString(string(b))
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.String:
			elem := reflect.New(val.Type().Elem()).Elem()
			elem.SetString(string(b))
			val.Set(reflect.Append(val, elem))
		case reflect.Uint8: // []byte
			val.SetBytes(b)
		}
	}
	return n, err
}

func (d *Decoder) decodeUvarint(val reflect.Value, v uint64) error {
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
		return d.decodeUvarint(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := d.decodeUvarint(elem, v); err != nil {
			return err
		}
		val.Set(reflect.Append(val, elem))
	}
	return nil
}

func (d *Decoder) decodeFixed64(val reflect.Value, v uint64) error {
	switch val.Kind() {
	//case reflect.Int64:
	//	return setInt(val, int64(v))
	//case reflect.Uint64:
	//	return setUint(val, v)
	case reflect.Float64:
		return setFloat(val, math.Float64frombits(v))
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return d.decodeFixed64(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := d.decodeFixed64(elem, v); err != nil {
			return err
		}
		val.Set(reflect.Append(val, elem))
	}
	return nil
}

func (d *Decoder) decodeFixed32(val reflect.Value, v uint32) error {
	switch val.Kind() {
	//case reflect.Int32:
	//	return setInt(val, int64(int32(v)))
	//case reflect.Uint32:
	//	return setUint(val, uint64(v))
	case reflect.Float32:
		return setFloat(val, float64(math.Float32frombits(v)))
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return d.decodeFixed32(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := d.decodeFixed32(elem, v); err != nil {
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
