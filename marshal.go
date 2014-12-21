// Package protobuf converts data structures to and from the wire format
// of protocol buffers.
package protobuf

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

const (
	wireVarint  = 0
	wireFixed64 = 1
	wireBytes   = 2
	wireFixed32 = 5
)

// Marshal returns the protocl buffer encoding of v. The struct underlying
// v must be a pointer.
//
// Marshal currently encodes all visible field, which does not allow
// distinction between 'required' and 'optional' fields. Marshal returns
// an error if a struct field type is not supported. If the buffer is too
// small, Marshal will panic.
func Marshal(data []byte, v interface{}) (n int, err error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return n, errors.New("v must be a pointer to a struct")
	}
	return marshal(data, val.Elem())
}

func marshal(data []byte, val reflect.Value) (n int, err error) {
	num := val.NumField()
	for i := 0; i < num; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		// TODO
		//fmt.Printf("%v\n", val.Type().Field(i).Tag)
		switch field.Kind() {
		case reflect.Int32, reflect.Int64:
			n += marshalUint(data[n:], i+1, uint64(field.Int()))
		case reflect.Uint32, reflect.Uint64:
			n += marshalUint(data[n:], i+1, field.Uint())
		case reflect.Float32:
			x := math.Float32bits(float32(field.Float()))
			n += marshalFloat32(data[n:], i+1, x)
		case reflect.Float64:
			x := math.Float64bits(field.Float())
			n += marshalFloat64(data[n:], i+1, x)
		case reflect.Bool:
			n += marshalBool(data[n:], i+1, field.Bool())
		case reflect.String:
			n += marshalString(data[n:], i+1, field.String())
		case reflect.Slice:
			m, err := marshalSlice(data[n:], i+1, field)
			n += m
			if err != nil {
				return n, err
			}
		default:
			return n, errors.New("invalid type: " + field.Kind().String())
		}
	}
	return n, err
}

// TODO: [][]byte
func marshalSlice(data []byte, key int, val reflect.Value) (n int, err error) {
	vlen := val.Len()
	switch val.Type().Elem().Kind() {
	case reflect.Int32, reflect.Int64:
		for i := 0; i < vlen; i++ {
			n += marshalUint(data[n:], key, uint64(val.Index(i).Int()))
		}
	case reflect.Uint32, reflect.Uint64:
		for i := 0; i < vlen; i++ {
			n += marshalUint(data[n:], key, val.Index(i).Uint())
		}
	case reflect.Float32:
		var x uint32
		for i := 0; i < vlen; i++ {
			x = math.Float32bits(float32(val.Index(i).Float()))
			n += marshalFloat32(data[n:], key, x)
		}
	case reflect.Float64:
		var x uint64
		for i := 0; i < vlen; i++ {
			x = math.Float64bits(val.Index(i).Float())
			n += marshalFloat64(data[n:], key, x)
		}
	case reflect.Bool:
		for i := 0; i < vlen; i++ {
			n += marshalBool(data[n:], key, val.Index(i).Bool())
		}
	case reflect.String:
		for i := 0; i < vlen; i++ {
			n += marshalString(data[n:], key, val.Index(i).String())
		}
	case reflect.Uint8: // byte slice
		n += marshalBytes(data[n:], key, val.Bytes())
	default:
		return n, errors.New("invalid type: " + val.Kind().String())
	}
	return n, err
}

func marshalUint(data []byte, key int, v uint64) (n int) {
	data[n], n = byte(key)<<3|wireVarint, n+1
	n += binary.PutUvarint(data[n:], v)
	return n
}

func marshalFloat32(data []byte, key int, v uint32) (n int) {
	data[n], n = byte(key)<<3|wireFixed32, n+1
	binary.LittleEndian.PutUint32(data[n:], v)
	return n + 4
}

func marshalFloat64(data []byte, key int, v uint64) (n int) {
	data[n], n = byte(key)<<3|wireFixed64, n+1
	binary.LittleEndian.PutUint64(data[n:], v)
	return n + 8
}

func marshalBool(data []byte, key int, v bool) (n int) {
	data[n], n = byte(key)<<3|wireVarint, n+1
	if v {
		data[n], n = 1, n+1
	} else {
		data[n], n = 0, n+1
	}
	return n
}

func marshalString(data []byte, key int, v string) (n int) {
	data[n], n = byte(key)<<3|wireBytes, n+1
	n += binary.PutUvarint(data[n:], uint64(len(v)))
	n += copy(data[n:], v)
	return n
}

func marshalBytes(data []byte, key int, v []byte) (n int) {
	data[n], n = byte(key)<<3|wireBytes, n+1
	n += binary.PutUvarint(data[n:], uint64(len(v)))
	n += copy(data[n:], v)
	return n
}

// Size returns the encoded protocol buffer size. The struct underlying v
// must be a pointer.
func Size(v interface{}) (n int, err error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return n, errors.New("v must be a pointer to a struct")
	}
	return size(val.Elem())
}

func size(val reflect.Value) (n int, err error) {
	num := val.NumField()
	for i := 0; i < num; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.Int32, reflect.Int64:
			n += 1 + uvarintSize(uint64(field.Int()))
		case reflect.Uint32, reflect.Uint64:
			n += 1 + uvarintSize(field.Uint())
		case reflect.Float32:
			n += 5
		case reflect.Float64:
			n += 9
		case reflect.Bool:
			n += 2
		case reflect.String:
			m := len(field.String())
			n += 1 + m + uvarintSize(uint64(m))
		case reflect.Slice:
			m, err := sliceSize(field, field.Len())
			n += m
			if err != nil {
				return n, err
			}
		default:
			return n, errors.New("invalid type: " + field.Kind().String())
		}
	}
	return n, err
}

func sliceSize(val reflect.Value, vlen int) (n int, err error) {
	switch val.Type().Elem().Kind() {
	case reflect.Int32, reflect.Int64:
		for i := 0; i < vlen; i++ {
			n += 1 + uvarintSize(uint64(val.Index(i).Int()))
		}
	case reflect.Uint32, reflect.Uint64:
		for i := 0; i < vlen; i++ {
			n += 1 + uvarintSize(val.Index(i).Uint())
		}
	case reflect.Float32:
		for i := 0; i < vlen; i++ {
			n += 5
		}
	case reflect.Float64:
		for i := 0; i < vlen; i++ {
			n += 9
		}
	case reflect.Bool:
		for i := 0; i < vlen; i++ {
			n += 2
		}
	case reflect.String:
		for i := 0; i < vlen; i++ {
			m := len(val.Index(i).String())
			n += 1 + m + uvarintSize(uint64(m))
		}
	case reflect.Uint8: // byte slice
		m := len(val.Bytes())
		n += 1 + m + uvarintSize(uint64(m))
	default:
		return n, errors.New("invalid type: " + val.Kind().String())
	}
	return n, err
}

func uvarintSize(v uint64) (n int) {
	for {
		n++
		v >>= 7
		if v == 0 {
			break
		}
	}
	return n
}
