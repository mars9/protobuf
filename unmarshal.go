package protobuf

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"time"
)

// Unmarshal parses the protocol buffer representation in data and places
// the decoded result in v. If the struct underlying v does not match the
// data, the results can be unpredictable.
//
// Unmarshal uses the inverse of the encodings that Marshal uses,
// allocating slices and pointers as necessary.
func Unmarshal(data []byte, v interface{}) (err error) {
	defer func() {
		if v := recover(); v != nil {
			if msg, ok := v.(string); ok {
				err = errors.New(msg)
				return
			}
			err = errors.New("malformed packet")
		}
	}()

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to a struct")
	}
	return unmarshal(data, val.Elem())
}

func unmarshal(data []byte, val reflect.Value) (err error) {
	num := val.NumField()
	var field reflect.Value

	for len(data) > 0 {
		key, n := binary.Uvarint(data)
		if n <= 0 {
			return errors.New("invalid field key")
		}
		data = data[n:]

		fnum := int(key >> 3)
		if fnum > 0 && fnum <= num {
			field = val.Field(fnum - 1)
		} else {
			break
		}

		switch key & 7 {
		case wireVarint:
			v, n := binary.Uvarint(data)
			if n <= 0 {
				return errors.New("invalid varint value")
			}
			if err = unmarshalUvarint(field, v); err != nil {
				return err
			}
			data = data[n:]
		case wireFixed32:
			if len(data) < 4 {
				return errors.New("bad 32-bit value")
			}
			v := binary.LittleEndian.Uint32(data)
			if err = unmarshalFixed32(field, v); err != nil {
				return err
			}
			data = data[4:]
		case wireFixed64:
			if len(data) < 8 {
				return errors.New("bad 64-bit value")
			}
			v := binary.LittleEndian.Uint64(data)
			if err = unmarshalFixed64(field, v); err != nil {
				return err
			}
			data = data[8:]
		case wireBytes:
			v, n := binary.Uvarint(data)
			if n <= 0 {
				return errors.New("invalid varint value")
			}
			data = data[n:]
			if err = unmarshalBytes(field, data[:v]); err != nil {
				return err
			}
			data = data[v:]
		}
	}
	return err
}

func unmarshalBytes(val reflect.Value, b []byte) error {
	switch val.Kind() {
	case reflect.String:
		val.SetString(string(b))
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.String:
			vtype := val.Type().Elem()
			elem := reflect.New(vtype).Elem()
			elem.SetString(string(b))
			val.Set(reflect.Append(val, elem))
		case reflect.Uint8: // []byte
			// copy bytes
			//val.SetBytes(append([]byte(nil), b...))
			val.SetBytes(b)
		case reflect.Slice: // [][]byte
			vtype := val.Type().Elem()
			elem := reflect.New(vtype).Elem()
			if err := unmarshalBytes(elem, b); err != nil {
				return err
			}
			val.Set(reflect.Append(val, elem))
		}
	case reflect.Struct:
		return unmarshal(b, val)
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return unmarshalBytes(val.Elem(), b)
	}
	return nil
}

func unmarshalFixed32(val reflect.Value, v uint32) error {
	switch val.Kind() {
	case reflect.Int32:
		return setInt(val, int64(int32(v)))
	case reflect.Uint32:
		return setUint(val, uint64(v))
	case reflect.Float32:
		return setFloat(val, float64(math.Float32frombits(v)))
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return unmarshalFixed32(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := unmarshalFixed32(elem, v); err != nil {
			return err
		}
		val.Set(reflect.Append(val, elem))
	}
	return nil
}

func unmarshalFixed64(val reflect.Value, v uint64) error {
	switch val.Kind() {
	case reflect.Int64:
		return setInt(val, int64(v))
	case reflect.Uint64:
		return setUint(val, v)
	case reflect.Float64:
		return setFloat(val, math.Float64frombits(v))
	case reflect.Ptr:
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		return unmarshalFixed64(val.Elem(), v)
	case reflect.Slice:
		elem := reflect.New(val.Type().Elem()).Elem()
		if err := unmarshalFixed64(elem, v); err != nil {
			return err
		}
		val.Set(reflect.Append(val, elem))
	}
	return nil
}

func unmarshalUvarint(val reflect.Value, v uint64) error {
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
		return unmarshalUvarint(val.Elem(), v)
	case reflect.Slice:
		return unmarshalUintSlice(val, v)
	}
	return nil
}

func unmarshalUintSlice(val reflect.Value, v uint64) error {
	vtype := val.Type().Elem()
	elem := reflect.New(vtype).Elem()
	switch vtype.Kind() {
	case reflect.Int32, reflect.Int64:
		if err := setInt(elem, int64(v)); err != nil {
			return err
		}
	case reflect.Uint32, reflect.Uint64:
		if err := setUint(elem, v); err != nil {
			return err
		}
	case reflect.Bool:
		if err := setBool(elem, v); err != nil {
			return err
		}
	}
	val.Set(reflect.Append(val, elem))
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
