package protobuf

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
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
		case wireFixed32:
			v := binary.LittleEndian.Uint32(data)
			if err = unmarshalFixed32(field, v); err != nil {
				return err
			}
			data = data[4:]
		case wireFixed64:
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

func unmarshalUvarint(val reflect.Value, v uint64) error {
	switch val.Kind() {
	case reflect.Int32, reflect.Int64:
		return setInt(val, int64(v))
	case reflect.Uint32, reflect.Uint64:
		return setUint(val, v)
	case reflect.Bool:
		return setBool(val, v)
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

func unmarshalFixed32(val reflect.Value, v uint32) error {
	switch val.Kind() {
	case reflect.Float32:
		x := float64(math.Float32frombits(v))
		if val.OverflowFloat(float64(x)) {
			return errors.New("float32 overflow")
		}
		val.SetFloat(float64(x))
	case reflect.Slice:
		return unmarshalFixed32Slice(val, v)
	}
	return nil
}

func unmarshalFixed32Slice(val reflect.Value, v uint32) error {
	vtype := val.Type().Elem()
	elem := reflect.New(vtype).Elem()
	switch vtype.Kind() {
	case reflect.Float32:
		if err := unmarshalFixed32(elem, v); err != nil {
			return err
		}
	}
	val.Set(reflect.Append(val, elem))
	return nil
}

func unmarshalFixed64(val reflect.Value, v uint64) error {
	switch val.Kind() {
	case reflect.Float64:
		x := math.Float64frombits(v)
		if val.OverflowFloat(x) {
			return errors.New("float64 overflow")
		}
		val.SetFloat(x)
	case reflect.Slice:
		return unmarshalFixed64Slice(val, v)
	}
	return nil
}

func unmarshalFixed64Slice(val reflect.Value, v uint64) error {
	vtype := val.Type().Elem()
	elem := reflect.New(vtype).Elem()
	switch vtype.Kind() {
	case reflect.Float64:
		if err := unmarshalFixed64(elem, v); err != nil {
			return err
		}
	}
	val.Set(reflect.Append(val, elem))
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
		case reflect.Uint8: // []byte
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
		return unmarshal(b, val.Elem())
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
		return errors.New("uint overflow")
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
