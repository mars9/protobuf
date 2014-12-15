// Package protobuf converts data structures to and from the wire format
// of protocol buffers.
package protobuf

import (
	"encoding/binary"
	"errors"
	"reflect"
)

const (
	wireVarint  = 0
	wireFixed64 = 1
	wireBytes   = 2
	wireFixed32 = 5
)

// Marshal returns the protocl buffer encoding of v. The struct
// underlying v must be a pointer.
func Marshal(data []byte, v interface{}) (n int, err error) {
	val := reflect.ValueOf(v).Elem()
	num := val.NumField()
	for i := 0; i < num; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.Uint32, reflect.Uint64:
			data[n], n = byte(i+1)<<3|wireVarint, n+1
			n += binary.PutUvarint(data[n:], field.Uint())

		case reflect.Int32, reflect.Int64:
			data[n], n = byte(i+1)<<3|wireVarint, n+1
			n += binary.PutUvarint(data[n:], uint64(field.Int()))

		case reflect.Bool:
			data[n], n = byte(i+1)<<3|wireVarint, n+1
			if field.Bool() {
				data[n], n = 1, n+1
			} else {
				data[n], n = 0, n+1
			}

		case reflect.String:
			data[n], n = byte(i+1)<<3|wireBytes, n+1
			s := field.String()
			n += binary.PutUvarint(data[n:], uint64(len(s)))
			n += copy(data[n:], s)

		case reflect.Slice:
			n += marshalSlice(data[n:], i+1, field)

		default:
			errors.New("unsupported type: " + field.Kind().String())
		}
	}
	return n, err
}

// TODO
func marshalSlice(data []byte, key int, val reflect.Value) (n int) {
	switch val.Type().Elem().Kind() {
	case reflect.Uint32, reflect.Uint64:

	case reflect.Int32, reflect.Int64:

	case reflect.Bool:

	case reflect.String:

	case reflect.Uint8: // byte slice
		data[n], n = byte(key)<<3|wireBytes, n+1
		b := val.Bytes()
		n += binary.PutUvarint(data[n:], uint64(len(b)))
		n += copy(data[n:], b)
	}
	return n
}

// Size returns the encoded protocol buffer size.
func Size(v interface{}) (n int, err error) {
	val := reflect.ValueOf(v).Elem()
	num := val.NumField()
	for i := 0; i < num; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.Uint32, reflect.Uint64:
			n += 1 + uvarintSize(field.Uint())
		case reflect.Int32, reflect.Int64:
			n += 1 + uvarintSize(uint64(field.Int()))
		case reflect.Bool:
			n += 2
		case reflect.String:
			m := len(field.String())
			n += 1 + m + uvarintSize(uint64(m))
		case reflect.Slice:
			n += sliceSize(field, field.Len())
		default:
			errors.New("unsupported type: " + field.Kind().String())
		}
	}
	return n, err
}

func sliceSize(val reflect.Value, vlen int) (n int) {
	elem := val.Type().Elem().Kind()
	if elem == reflect.Uint8 { // byte slice
		m := len(val.Bytes())
		n += 1 + m + uvarintSize(uint64(m))
		return n
	}
	if vlen == 0 {
		return 0
	}

	switch val.Type().Elem().Kind() {
	case reflect.Uint32, reflect.Uint64:
		for i := 0; i < vlen; i++ {
			n += uvarintSize(val.Index(i).Uint())
		}
	case reflect.Int32, reflect.Int64:
		for i := 0; i < vlen; i++ {
			n += uvarintSize(uint64(val.Index(i).Int()))
		}
	case reflect.Bool:
		for i := 0; i < vlen; i++ {
			n++
		}
		// TODO
		//case reflect.String:
		//	for i := 0; i < vlen; i++ {
		//		n += uvarintSize(uint64(len(val.Index(i).String())))
		//	}
	}
	n += 1 + uvarintSize(uint64(n))
	return n
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
