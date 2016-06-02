package protobuf

import (
	"errors"
	"io"
	"math"
	"reflect"
	"time"
)

const (
	wireVarint  = 0
	wireFixed64 = 1
	wireBytes   = 2
	wireFixed32 = 5
)

// Writer defines the encode writer. Typically this is a *bufio.Writer.
type Writer interface {
	WriteString(string) (int, error)
	io.ByteWriter
	io.Writer
}

// Encoder manages the transmission of type and data information to the
// other side of a connection.
type Encoder struct {
	w   Writer
	max int
}

// NewEncoder returns a new encoder that will transmit on the io.Writer.
//
// Max defines the maximum size that can be transmitted, if max is 0,
// the maximum message size is not checked.
func NewEncoder(w Writer, max int) *Encoder {
	return &Encoder{w: w, max: max}
}

// Encode transmits the data item represented by the empty interface value,
// guaranteeing that all necessary type information has been transmitted
// first.
//
// Encode first writes the varint encoded message size, traverses the value
// v recursively and writes the Protocol Buffer encoding of v. The struct
// underlying v must be a pointer.
//
// Encode currently encodes all visible field and ignores unsupported
// Protocol Buffer struct field types.
func (e *Encoder) Encode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}

	val = val.Elem()
	fields := val.NumField()
	size := sizeStruct(val, fields)

	if err := writeLength(e.w, size, e.max); err != nil {
		return err
	}
	return e.encodeStruct(val, fields)
}

func (e *Encoder) encodeStruct(val reflect.Value, fields int) error {
	var field reflect.Value
	var custom bool
	var err error
	for i := 0; i < fields && err == nil; i++ {
		field = val.Field(i)
		if !field.CanSet() {
			continue
		}

		if custom, err = e.encodeCustom(field, i+1); err != nil || custom {
			continue
		}

		switch field.Kind() {
		case reflect.Interface:
			panic("interface encoding not implemented")
		case reflect.Struct:
			err = e.writeStruct(i+1, field)
		case reflect.Ptr:
			v := field.Elem()
			switch v.Kind() {
			case reflect.Struct:
				err = e.writeStruct(i+1, v)
			case reflect.Ptr:
				// nothing
			case reflect.Slice:
				// nothing
			default:
				err = e.encodeBasic(v, i+1)
			}
		case reflect.Slice:
			err = e.encodeSlice(field, i+1)
		default:
			err = e.encodeBasic(field, i+1)
		}
	}
	return err
}

func (e *Encoder) encodeCustom(val reflect.Value, key int) (bool, error) {
	itype := val.Interface()
	if t, ok := itype.(error); ok || val.Type() == errorType {
		if val.IsNil() {
			return true, nil
		}
		return true, e.writeString(key, t.Error())
	}
	if t, ok := itype.(time.Time); ok {
		if t.IsZero() {
			return true, nil
		}
		return true, e.writeUvarint(key, uint64(t.UnixNano()))
	}
	return false, nil
}

func (e *Encoder) encodeSlice(val reflect.Value, key int) (err error) {
	vlen := val.Len()
	if vlen == 0 {
		return nil
	}

	switch val.Type().Elem().Kind() {
	case reflect.Int32, reflect.Int64:
		for i := 0; i < vlen && err == nil; i++ {
			err = e.writeUvarint(key, uint64(val.Index(i).Int()))
		}
	case reflect.Uint32, reflect.Uint64:
		for i := 0; i < vlen && err == nil; i++ {
			err = e.writeUvarint(key, val.Index(i).Uint())
		}
	case reflect.Float32:
		for i := 0; i < vlen && err == nil; i++ {
			x := math.Float32bits(float32(val.Index(i).Float()))
			err = e.writeFixed32(key, x)
		}
	case reflect.Float64:
		for i := 0; i < vlen && err == nil; i++ {
			x := math.Float64bits(val.Index(i).Float())
			err = e.writeFixed64(key, x)
		}
	case reflect.Bool:
		for i := 0; i < vlen && err == nil; i++ {
			err = e.writeBool(key, val.Index(i).Bool())
		}
	case reflect.String:
		for i := 0; i < vlen && err == nil; i++ {
			err = e.writeString(key, val.Index(i).String())
		}
	case reflect.Uint8:
		err = e.writeBytes(key, val.Bytes())
	case reflect.Slice:
		for i := 0; i < vlen && err == nil; i++ {
			err = e.encodeSlice(val.Index(i), key)
		}
	case reflect.Struct:
		panic("struct slice not implemented")
	case reflect.Ptr:
		panic("pointer slice not implemented")
	}
	return err
}

func (e *Encoder) encodeBasic(val reflect.Value, key int) error {
	switch val.Kind() {
	case reflect.Int32, reflect.Int64:
		v := uint64(val.Int())
		if v == 0 {
			return nil
		}
		return e.writeUvarint(key, v)
	case reflect.Uint32, reflect.Uint64:
		v := val.Uint()
		if v == 0 {
			return nil
		}
		return e.writeUvarint(key, v)
	case reflect.Float32:
		v := math.Float32bits(float32(val.Float()))
		if v == 0 {
			return nil
		}
		return e.writeFixed32(key, v)
	case reflect.Float64:
		v := math.Float64bits(val.Float())
		if v == 0 {
			return nil
		}
		return e.writeFixed64(key, v)
	case reflect.Bool:
		v := val.Bool()
		if !v {
			return nil
		}
		return e.writeBool(key, v)
	case reflect.String:
		v := val.String()
		if v == "" {
			return nil
		}
		return e.writeString(key, v)
	}
	return nil
}

func (e *Encoder) writeUvarint(key int, v uint64) (err error) {
	if err = e.w.WriteByte(byte(key)<<3 | wireVarint); err != nil {
		return err
	}
	return writeUvarint(e.w, v)
}

func (e *Encoder) writeBool(key int, v bool) (err error) {
	if err = e.w.WriteByte(byte(key)<<3 | wireVarint); err != nil {
		return err
	}
	if v {
		return e.w.WriteByte(1)
	}
	return e.w.WriteByte(0)
}

func (e *Encoder) writeFixed32(key int, v uint32) (err error) {
	if err = e.w.WriteByte(byte(key)<<3 | wireFixed32); err != nil {
		return err
	}
	return writeFixed32(e.w, v)
}

func (e *Encoder) writeFixed64(key int, v uint64) (err error) {
	if err = e.w.WriteByte(byte(key)<<3 | wireFixed64); err != nil {
		return err
	}
	return writeFixed64(e.w, v)
}

func (e *Encoder) writeBytes(key int, v []byte) (err error) {
	if err = e.w.WriteByte(byte(key)<<3 | wireBytes); err != nil {
		return err
	}
	if err = writeUvarint(e.w, uint64(len(v))); err != nil {
		return err
	}
	_, err = e.w.Write(v)
	return err
}

func (e *Encoder) writeString(key int, v string) (err error) {
	if err = e.w.WriteByte(byte(key)<<3 | wireBytes); err != nil {
		return err
	}
	if err = writeUvarint(e.w, uint64(len(v))); err != nil {
		return err
	}
	_, err = e.w.WriteString(v)
	return err
}

func (e *Encoder) writeStruct(key int, v reflect.Value) (err error) {
	if err = e.w.WriteByte(byte(key)<<3 | wireBytes); err != nil {
		return err
	}

	n := uint64(sizeStruct(v, v.NumField()))
	if err = writeUvarint(e.w, n); err != nil {
		return err
	}
	return e.encodeStruct(v, v.NumField())
}
