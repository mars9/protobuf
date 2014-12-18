package protobuf

import (
	"encoding/binary"
	"errors"
	"reflect"
)

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
			if err = unmarshalUvarint(field, v); err != nil {
				return err
			}
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

func unmarshalUvarint(val reflect.Value, v uint64) error {
	switch val.Kind() {
	case reflect.Int32, reflect.Int64:
		return unmarshalInt(val, int64(v))
	case reflect.Uint32, reflect.Uint64:
		return unmarshalUint(val, v)
	case reflect.Bool:
		return unmarshalBool(val, v)
	case reflect.Slice:
		return unmarshalUintSlice(val, v)
	}
	return nil
}

func unmarshalUintSlice(val reflect.Value, v uint64) (err error) {
	vtype := val.Type().Elem()
	elem := reflect.New(vtype).Elem()
	switch vtype.Kind() {
	case reflect.Int32, reflect.Int64:
		if err = unmarshalInt(elem, int64(v)); err != nil {
			return err
		}
	case reflect.Uint32, reflect.Uint64:
		if err = unmarshalUint(elem, v); err != nil {
			return err
		}
	case reflect.Bool:
		if err = unmarshalBool(elem, v); err != nil {
			return err
		}
	}
	val.Set(reflect.Append(val, elem))
	return nil
}

func unmarshalUint(val reflect.Value, v uint64) error {
	if val.OverflowUint(v) {
		return errors.New("uint overflow")
	}
	val.SetUint(v)
	return nil
}

func unmarshalInt(val reflect.Value, v int64) error {
	if val.OverflowInt(v) {
		return errors.New("uint overflow")
	}
	val.SetInt(v)
	return nil
}

func unmarshalBool(val reflect.Value, v uint64) error {
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
		case reflect.Uint8: // byte slice
			val.SetBytes(b)
		}
	}
	return nil
}
