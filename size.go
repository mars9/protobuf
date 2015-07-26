package protobuf

import (
	"math"
	"reflect"
	"time"
)

func sizeStruct(val reflect.Value, fields int) (n int) {
	var custom bool
	var m int
	for i := 0; i < fields; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		if custom, m = sizeCustom(field); custom {
			n += m
			continue
		}

		switch field.Kind() {
		case reflect.Interface:
			n += sizeSlice(field.Elem())
		case reflect.Struct:
			m = sizeStruct(field, field.NumField())
			n += 1 + m + uvarintSize(uint64(m))
		case reflect.Ptr:
			v := field.Elem()
			switch v.Kind() {
			case reflect.Struct:
				m = sizeStruct(v, v.NumField())
				n += 1 + m + uvarintSize(uint64(m))
			case reflect.Ptr:
				// nothing
			case reflect.Slice:
				// nothing
			default:
				n += sizeType(v)
			}
		case reflect.Slice:
			n += sizeSlice(field)
		default:
			n += sizeType(field)
		}
	}
	return n
}

func sizeCustom(val reflect.Value) (bool, int) {
	itype := val.Interface()
	if t, ok := itype.(error); ok || val.Type() == errorType {
		if val.IsNil() {
			return true, 0
		}
		v := len(t.Error())
		return true, 1 + v + uvarintSize(uint64(v))
	}
	if t, ok := itype.(time.Time); ok {
		if t.IsZero() {
			return true, 0
		}
		return true, 1 + uvarintSize(uint64(t.UnixNano()))
	}
	return false, 0
}

func sizeSlice(val reflect.Value) (n int) {
	vlen := val.Len()
	if vlen == 0 {
		return 0
	}

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
	case reflect.Uint8:
		m := len(val.Bytes())
		n += 1 + m + uvarintSize(uint64(m))
	case reflect.Slice:
		var v reflect.Value
		for i := 0; i < vlen; i++ {
			v = val.Index(i)
			if v.Type().Elem().Kind() == reflect.Uint8 {
				n += sizeSlice(v)
			}
		}
	case reflect.Struct:
		for i := 0; i < vlen; i++ {
			m := sizeStruct(val.Index(i), val.Index(i).NumField())
			n += 1 + m + uvarintSize(uint64(m))
		}
	case reflect.Ptr:
		for i := 0; i < vlen; i++ {
			v := val.Index(i).Elem()
			if v.Kind() == reflect.Struct {
				m := sizeStruct(v, v.NumField())
				n += 1 + m + uvarintSize(uint64(m))
			}
		}
	}
	return n
}

func sizeType(val reflect.Value) (n int) {
	switch val.Kind() {
	case reflect.Int32, reflect.Int64:
		v := uint64(val.Int())
		if v == 0 {
			return 0
		}
		n += 1 + uvarintSize(v)
	case reflect.Uint32, reflect.Uint64:
		v := val.Uint()
		if v == 0 {
			return 0
		}
		n += 1 + uvarintSize(v)
	case reflect.Float32:
		if math.Float32bits(float32(val.Float())) == 0 {
			return 0
		}
		n += 5
	case reflect.Float64:
		if math.Float64bits(val.Float()) == 0 {
			return 0
		}
		n += 9
	case reflect.Bool:
		if !val.Bool() {
			return 0
		}
		n += 2
	case reflect.String:
		m := len(val.String())
		if m == 0 {
			return 0
		}
		n += 1 + m + uvarintSize(uint64(m))
	}
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
