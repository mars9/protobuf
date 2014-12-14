package protobuf

import (
	"encoding/binary"
	"errors"
	"reflect"
)

const (
	maxInt8  = 1<<7 - 1
	minInt8  = -1 << 7
	maxInt16 = 1<<15 - 1
	minInt16 = -1 << 15
	maxInt32 = 1<<31 - 1
	minInt32 = -1 << 31
	maxInt64 = 1<<63 - 1
	minInt64 = -1 << 63

	maxUint8  = 1<<8 - 1
	maxUint16 = 1<<16 - 1
	maxUint32 = 1<<32 - 1
	maxUint64 = 1<<64 - 1
)

// Unmarshal parses the protocol buffer representation in data and places
// the decoded result in v. If the struct underlying v does not match the
// data, the results can be unpredictable.
func Unmarshal(data []byte, v interface{}) (err error) {
	val := reflect.ValueOf(v).Elem()
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
			if err = putUint(field, v); err != nil {
				return err
			}

		case wireBytes:
			v, n := binary.Uvarint(data)
			if n <= 0 {
				return errors.New("invalid varint value")
			}
			data = data[n:]
			if err = putBytes(field, data[:v]); err != nil {
				return err
			}
			data = data[v:]

		case wireFixed64:
			// TODO
		case wireFixed32:
			// TODO
		}
	}
	return err
}

func putUint(field reflect.Value, v uint64) error {
	switch field.Kind() {
	case reflect.Int8:
		x := int64(v)
		if x < minInt8 && x > maxInt8 {
			return errors.New("int8 overflow")
		}
		field.SetInt(x)
	case reflect.Int16:
		x := int64(v)
		if x < minInt16 && x > maxInt16 {
			return errors.New("int8 overflow")
		}
		field.SetInt(x)
	case reflect.Int32:
		x := int64(v)
		if x < minInt32 && x > maxInt32 {
			return errors.New("int32 overflow")
		}
		field.SetInt(x)
	case reflect.Int64:
		x := int64(v)
		if x < minInt64 && x > maxInt64 {
			return errors.New("int64 overflow")
		}
		field.SetInt(x)
	case reflect.Uint8:
		if v > maxUint8 {
			return errors.New("uint8 overflow")
		}
		field.SetUint(v)
	case reflect.Uint16:
		if v > maxUint16 {
			return errors.New("uint16 overflow")
		}
		field.SetUint(v)
	case reflect.Uint32:
		if v > maxUint32 {
			return errors.New("uint32 overflow")
		}
		field.SetUint(v)
	case reflect.Uint64:
		if v > maxUint64 {
			return errors.New("uint64 overflow")
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
	default:
		// TODO: handle int and uint
	}
	return nil
}

func putBytes(field reflect.Value, b []byte) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(string(b))
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			field.SetBytes(b)
		} else {
			// TODO: not a byte slice
		}
	default:
		// TODO: can this happen?
	}
	return nil
}
