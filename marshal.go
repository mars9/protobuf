// Package protobuf converts data structures to and from the wire format
// of protocol buffers.
package protobuf

import (
	"encoding/binary"
	"errors"
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

// Marshal traverses the value v recursively and returns the protocol
// buffer encoding of v. The struct underlying v must be a pointer.
//
// Marshal currently encodes all visible field, which does not allow
// distinction between 'required' and 'optional' fields. Marshal ignores
// unsupported struct field types. If the buffer is too small, Marshal
// will panic.
func Marshal(data []byte, v interface{}) (n int, err error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("v must be a pointer to a struct")
	}
	return marshalStruct(data, val.Elem())
}

func marshalStruct(data []byte, val reflect.Value) (n int, err error) {
	num := val.NumField()
	var ftype, fopt int

	for i := 0; i < num; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		ftype, fopt = parseTag(val.Type().Field(i).Tag.Get("protobuf"))
		if ftype > ftypeStart && ftype < ftypeEnd {
			m, err := marshalTag(ftype, fopt, data[n:], i+1, field)
			if err != nil {
				return n + m, err
			}
			n += m
			continue
		}

		if v, ok := field.Interface().(time.Time); ok {
			n += putUint(data[n:], i+1, uint64(v.UnixNano()))
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			m := sizeStruct(field)
			data[n], n = byte(i+1)<<3|wireBytes, n+1
			n += binary.PutUvarint(data[n:], uint64(m))
			m, err = marshalStruct(data[n:], field)
			if err != nil {
				return n + m, err
			}
			n += m
		case reflect.Ptr:
			if fopt == required && field.IsNil() {
				return n, errors.New("required field not set")
			}

			v := field.Elem()
			switch v.Kind() {
			case reflect.Struct:
				m := sizeStruct(field.Elem())
				data[n], n = byte(i+1)<<3|wireBytes, n+1
				n += binary.PutUvarint(data[n:], uint64(m))
				m, err = marshalStruct(data[n:], field.Elem())
				if err != nil {
					return n + m, err
				}
				n += m
			case reflect.Ptr:
				// nothing
			case reflect.Slice:
				// nothing
			default:
				n += marshalType(data[n:], i+1, v)
			}
		case reflect.Slice:
			n += marshalSlice(data[n:], i+1, field)
		default:
			n += marshalType(data[n:], i+1, field)
		}
	}
	return n, err
}

func marshalTag(ftype, fopt int, data []byte, key int, val reflect.Value) (n int, err error) {
	if val.Kind() == reflect.Ptr {
		if fopt == required && val.IsNil() {
			return n, errors.New("required field not set")
		}
		m, err := marshalTag(ftype, fopt, data, key, val.Elem())
		if err != nil {
			return n + m, err
		}
		n += m
	}

	switch ftype {
	case sfixed32:
		switch val.Kind() {
		case reflect.Int32:
			n += putFixed32(data, key, uint32(val.Int()))
		case reflect.Slice:
			vlen := val.Len()
			for i := 0; i < vlen; i++ {
				n += putFixed32(data[n:], key, uint32(val.Index(i).Int()))
			}
		}
	case sfixed64:
		switch val.Kind() {
		case reflect.Int64:
			n += putFixed64(data, key, uint64(val.Int()))
		case reflect.Slice:
			vlen := val.Len()
			for i := 0; i < vlen; i++ {
				n += putFixed64(data[n:], key, uint64(val.Index(i).Int()))
			}
		}
	case fixed32:
		switch val.Kind() {
		case reflect.Uint32:
			n += putFixed32(data, key, uint32(val.Uint()))
		case reflect.Slice:
			vlen := val.Len()
			for i := 0; i < vlen; i++ {
				n += putFixed32(data[n:], key, uint32(val.Index(i).Uint()))
			}
		}
	case fixed64:
		switch val.Kind() {
		case reflect.Uint64:
			n += putFixed64(data, key, val.Uint())
		case reflect.Slice:
			vlen := val.Len()
			for i := 0; i < vlen; i++ {
				n += putFixed64(data[n:], key, val.Index(i).Uint())
			}
		}
	case sint32:
		// TODO
	case sint64:
		// TODO
	}
	return n, err
}

func marshalSlice(data []byte, key int, val reflect.Value) (n int) {
	vlen := val.Len()
	switch val.Type().Elem().Kind() {
	case reflect.Int32, reflect.Int64:
		for i := 0; i < vlen; i++ {
			n += putUint(data[n:], key, uint64(val.Index(i).Int()))
		}
	case reflect.Uint32, reflect.Uint64:
		for i := 0; i < vlen; i++ {
			n += putUint(data[n:], key, val.Index(i).Uint())
		}
	case reflect.Float32:
		var x uint32
		for i := 0; i < vlen; i++ {
			x = math.Float32bits(float32(val.Index(i).Float()))
			n += putFixed32(data[n:], key, x)
		}
	case reflect.Float64:
		var x uint64
		for i := 0; i < vlen; i++ {
			x = math.Float64bits(val.Index(i).Float())
			n += putFixed64(data[n:], key, x)
		}
	case reflect.Bool:
		for i := 0; i < vlen; i++ {
			n += putBool(data[n:], key, val.Index(i).Bool())
		}
	case reflect.String:
		for i := 0; i < vlen; i++ {
			n += putString(data[n:], key, val.Index(i).String())
		}
	case reflect.Uint8:
		n += putBytes(data[n:], key, val.Bytes())
	case reflect.Slice:
		for i := 0; i < vlen; i++ {
			v := val.Index(i)
			if v.Type().Elem().Kind() == reflect.Uint8 {
				n += marshalSlice(data[n:], key, v)
			}
		}
	}
	return n
}

func marshalType(data []byte, key int, val reflect.Value) (n int) {
	switch val.Kind() {
	case reflect.Int32, reflect.Int64:
		n += putUint(data, key, uint64(val.Int()))
	case reflect.Uint32, reflect.Uint64:
		n += putUint(data, key, val.Uint())
	case reflect.Float32:
		x := math.Float32bits(float32(val.Float()))
		n += putFixed32(data, key, x)
	case reflect.Float64:
		x := math.Float64bits(val.Float())
		n += putFixed64(data, key, x)
	case reflect.Bool:
		n += putBool(data, key, val.Bool())
	case reflect.String:
		n += putString(data, key, val.String())
	}
	return n
}

func putUint(data []byte, key int, v uint64) (n int) {
	data[n], n = byte(key)<<3|wireVarint, n+1
	n += binary.PutUvarint(data[n:], v)
	return n
}

func putFixed32(data []byte, key int, v uint32) (n int) {
	data[n], n = byte(key)<<3|wireFixed32, n+1
	binary.LittleEndian.PutUint32(data[n:], v)
	return n + 4
}

func putFixed64(data []byte, key int, v uint64) (n int) {
	data[n], n = byte(key)<<3|wireFixed64, n+1
	binary.LittleEndian.PutUint64(data[n:], v)
	return n + 8
}

func putBool(data []byte, key int, v bool) (n int) {
	data[n], n = byte(key)<<3|wireVarint, n+1
	if v {
		data[n], n = 1, n+1
	} else {
		data[n], n = 0, n+1
	}
	return n
}

func putString(data []byte, key int, v string) (n int) {
	data[n], n = byte(key)<<3|wireBytes, n+1
	n += binary.PutUvarint(data[n:], uint64(len(v)))
	n += copy(data[n:], v)
	return n
}

func putBytes(data []byte, key int, v []byte) (n int) {
	data[n], n = byte(key)<<3|wireBytes, n+1
	n += binary.PutUvarint(data[n:], uint64(len(v)))
	n += copy(data[n:], v)
	return n
}
