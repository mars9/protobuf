package protobuf

import (
	"encoding/binary"
	"errors"
	"reflect"
)

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
			data = data[n:]
			if err = putNum(field, v); err != nil {
				return err
			}
		case wireBytes:
			v, n := binary.Uvarint(data)
			if n <= 0 {
				return errors.New("invalid varint value")
			}
			data = data[n:]
			if err = putSlice(field, data[:v]); err != nil {
				return err
			}
			data = data[v:]
		}
	}
	return err
}

// Unmarshal parses the protocol buffer representation in data and places
// the decoded result in v. If the struct underlying v does not match the
// data, the results can be unpredictable.
func Unmarshal(data []byte, v interface{}) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("malformed packet")
		}
	}()

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer to a struct")
	}
	return unmarshal(data, val.Elem())
}

func putNum(field reflect.Value, v uint64) error {
	switch field.Kind() {
	case reflect.Int32, reflect.Int64:
		x := int64(v)
		if field.OverflowInt(x) {
			return errors.New("int overflow")
		}
		field.SetInt(x)
	case reflect.Uint32, reflect.Uint64:
		if field.OverflowUint(v) {
			return errors.New("uint overflow")
		}
		field.SetUint(v)
	case reflect.Bool:
		if v != 0 && v != 1 {
			return errors.New("invalid bool value")
		}
		if v == 0 {
			field.SetBool(false)
		} else {
			field.SetBool(true)
		}
	}
	return nil
}

func putSlice(field reflect.Value, b []byte) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(string(b))

	case reflect.Slice:
		switch field.Type().Elem().Kind() {
		case reflect.Uint32, reflect.Uint64:

		case reflect.Int32, reflect.Int64:

		case reflect.Bool:

		case reflect.String:

		case reflect.Uint8: // byte slice
			field.SetBytes(b)
		}
	}
	return nil
}
