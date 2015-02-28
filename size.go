package protobuf

import (
	"reflect"
	"time"
)

// Size traverses the value v recursively and returns the encoded
// protocol buffer size. The struct underlying v must be a pointer.
func Size(v interface{}) (n int) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("v must be a pointer to a struct")
	}
	return sizeStruct(val.Elem())
}

func sizeStruct(val reflect.Value) (n int) {
	num := val.NumField()
	var ftype int

	for i := 0; i < num; i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		ftype, _ = parseTag(val.Type().Field(i).Tag.Get("protobuf"))
		if ftype > ftypeStart && ftype < ftypeEnd {
			n += sizeTag(ftype, field)
			continue
		}

		if v, ok := field.Interface().(time.Time); ok {
			n += 1 + uvarintSize(uint64(v.UnixNano()))
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			m := sizeStruct(field)
			n += 1 + m + uvarintSize(uint64(m))
		case reflect.Ptr:
			v := field.Elem()
			switch v.Kind() {
			case reflect.Struct:
				m := sizeStruct(v)
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

func sizeTag(ftype int, val reflect.Value) (n int) {
	switch ftype {
	case sfixed32, fixed32:
		switch val.Kind() {
		case reflect.Int32, reflect.Uint32:
			n += 5
		case reflect.Ptr:
			n += sizeTag(ftype, val.Elem())
		case reflect.Slice:
			vlen := val.Len()
			for i := 0; i < vlen; i++ {
				n += 5
			}
		}
	case sfixed64, fixed64:
		switch val.Kind() {
		case reflect.Int64, reflect.Uint64:
			n += 9
		case reflect.Ptr:
			n += sizeTag(ftype, val.Elem())
		case reflect.Slice:
			vlen := val.Len()
			for i := 0; i < vlen; i++ {
				n += 9
			}
		}
	case sint32:
		// TODO
	case sint64:
		// TODO
	}
	return n
}

func sizeSlice(val reflect.Value) (n int) {
	vlen := val.Len()
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
			m := sizeStruct(val.Index(i))
			n += 1 + m + uvarintSize(uint64(m))
		}
	case reflect.Ptr:
		for i := 0; i < vlen; i++ {
			v := val.Index(i).Elem()
			switch v.Kind() {
			case reflect.Struct:
				m := sizeStruct(v)
				n += 1 + m + uvarintSize(uint64(m))
			default:
				// TODO
			}
		}
	}
	return n
}

func sizeType(val reflect.Value) (n int) {
	switch val.Kind() {
	case reflect.Int32, reflect.Int64:
		n += 1 + uvarintSize(uint64(val.Int()))
	case reflect.Uint32, reflect.Uint64:
		n += 1 + uvarintSize(val.Uint())
	case reflect.Float32:
		n += 5
	case reflect.Float64:
		n += 9
	case reflect.Bool:
		n += 2
	case reflect.String:
		m := len(val.String())
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
