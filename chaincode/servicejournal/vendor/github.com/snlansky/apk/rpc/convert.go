package rpc

import (
	"strconv"
	"reflect"
	"unicode/utf8"
	"unicode"
	"strings"
)

//----------------------------------------------------------------------------------------
// JSON standard : all number are Number type, that is float64 in golang.
func convertParam(v interface{}, targetType reflect.Type) (newV reflect.Value, ok bool) {
	defer func() {
		if re := recover(); re != nil {
			ok = false
			logger.Errorf("[RPC]: convertParam, recover: %s", re)
		}
	}()

	ok = true

	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)
	tKind := targetType.Kind()

	if targetType.Kind() == reflect.Interface {
		newV = vValue
		return
	}

	switch vType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		f := convertObj2Float64(v)
		return convertFloat642Value(f, targetType)
	case reflect.Bool:
		return convertBool2Value(v.(bool), targetType)
	case reflect.String:
		return convertString2Value(v.(string), targetType)
	case reflect.Slice:
		return convertSlice2Value(v, targetType)
	case reflect.Map:
		if (tKind == reflect.Ptr && targetType.Elem().Kind() == reflect.Struct) || tKind == reflect.Struct {
			return convertObj2Value(v.(map[string]interface{}), targetType)
		} else if tKind == reflect.Map {
			return convertMap2Value(v, targetType)
		} else {
			ok = false
		}
	default:
		ok = false
	}
	return
}

func convertObj2Float64(v interface{}) (f float64) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Int:
		f = float64(v.(int))
	case reflect.Int8:
		f = float64(v.(int8))
	case reflect.Int16:
		f = float64(v.(int16))
	case reflect.Int32:
		f = float64(v.(int32))
	case reflect.Int64:
		f = float64(v.(int64))
	case reflect.Uint:
		f = float64(v.(uint))
	case reflect.Uint8:
		f = float64(v.(uint8))
	case reflect.Uint16:
		f = float64(v.(uint16))
	case reflect.Uint32:
		f = float64(v.(uint32))
	case reflect.Uint64:
		f = float64(v.(uint64))
	case reflect.Float32:
		f = float64(v.(float32))
	case reflect.Float64:
		f = v.(float64)
	}
	return
}
func convertFloat642Value(v float64, t reflect.Type) (newV reflect.Value, ok bool) {
	ok = true
	switch t.Kind() {
	case reflect.Int:
		newV = reflect.ValueOf(int(v))
	case reflect.Int8:
		newV = reflect.ValueOf(int8(v))
	case reflect.Int16:
		newV = reflect.ValueOf(int16(v))
	case reflect.Int32:
		newV = reflect.ValueOf(int32(v))
	case reflect.Int64:
		newV = reflect.ValueOf(int64(v))
	case reflect.Uint:
		newV = reflect.ValueOf(uint(v))
	case reflect.Uint8:
		newV = reflect.ValueOf(uint8(v))
	case reflect.Uint16:
		newV = reflect.ValueOf(uint16(v))
	case reflect.Uint32:
		newV = reflect.ValueOf(uint32(v))
	case reflect.Uint64:
		newV = reflect.ValueOf(uint64(v))
	case reflect.Float32:
		newV = reflect.ValueOf(float32(v))
	case reflect.Float64:
		newV = reflect.ValueOf(float64(v))
	default:
		ok = false
	}
	return
}

func convertBool2Value(v bool, t reflect.Type) (newV reflect.Value, ok bool) {
	ok = true
	if t.Kind() == reflect.Bool {
		newV = reflect.ValueOf(v)
	} else {
		ok = false
	}
	return
}

func convertObj2Value(v map[string]interface{}, targetType reflect.Type) (newV reflect.Value, ok bool) {
	ok = true
	tempType := targetType
	if targetType.Kind() == reflect.Ptr {
		tempType = targetType.Elem()
	}
	newV = reflect.New(tempType)

	for i := 0; i < tempType.NumField(); i++ {
		field := tempType.Field(i)
		char, _ := utf8.DecodeRuneInString(field.Name)
		if !unicode.IsUpper(char) {
			continue
		}
		fieldValue := newV.Elem().FieldByName(field.Name)

		jsonName := field.Name
		fieldTag := field.Tag.Get("json")
		if fieldTag != "" {
			info := strings.Split(fieldTag, ",")
			jsonName = info[0]
		}

		if jsonFieldValue, find := v[jsonName]; find {
			if jsonFieldValue == nil {
				continue
			}
			value, ok := convertParam(jsonFieldValue, field.Type)
			if !ok {
				return newV, false
			}
			fieldValue.Set(value)
		}
	}
	if targetType.Kind() == reflect.Struct {
		newV = newV.Elem()
	}
	return
}

func convertMap2Value(v interface{}, targetType reflect.Type) (newV reflect.Value, ok bool) {
	// key must string
	ok = true
	// key type must equal
	if targetType.Key().Kind() != targetType.Key().Kind() {
		return newV, false
	}

	vValue := reflect.ValueOf(v)
	newV = reflect.MakeMap(targetType)
	for _, key := range vValue.MapKeys() {
		value, success := convertParam(vValue.MapIndex(key).Interface(), targetType.Elem())
		if !success {
			return newV, false
		}
		newV.SetMapIndex(key, value)
	}
	return

}

func convertString2Value(v string, targetType reflect.Type) (newV reflect.Value, ok bool) {
	ok = true

	switch targetType.Kind() {
	case reflect.Int:
		val, err := strconv.Atoi(v)
		check(err)
		newV = reflect.ValueOf(val)
	case reflect.Uint8:
		val, err := strconv.ParseUint(v, 10, 8)
		check(err)
		newV = reflect.ValueOf(uint8(val))
	case reflect.Uint16:
		val, err := strconv.ParseUint(v, 10, 16)
		check(err)
		newV = reflect.ValueOf(uint16(val))
	case reflect.Uint32:
		val, err := strconv.ParseUint(v, 10, 32)
		check(err)
		newV = reflect.ValueOf(uint32(val))
	case reflect.Uint64:
		val, err := strconv.ParseUint(v, 10, 64)
		check(err)
		newV = reflect.ValueOf(val)
	case reflect.Int8:
		val, err := strconv.ParseInt(v, 10, 8)
		check(err)
		newV = reflect.ValueOf(int8(val))
	case reflect.Int16:
		val, err := strconv.ParseInt(v, 10, 16)
		check(err)
		newV = reflect.ValueOf(int16(val))
	case reflect.Int32:
		val, err := strconv.ParseInt(v, 10, 32)
		check(err)
		newV = reflect.ValueOf(int32(val))
	case reflect.Int64:
		val, err := strconv.ParseInt(v, 10, 64)
		check(err)
		newV = reflect.ValueOf(val)
	case reflect.Float32:
		val, err := strconv.ParseFloat(v, 32)
		check(err)
		newV = reflect.ValueOf(float32(val))
	case reflect.Float64:
		val, err := strconv.ParseFloat(v, 64)
		check(err)
		newV = reflect.ValueOf(val)
	case reflect.String:
		newV = reflect.ValueOf(v)
	default:
		ok = false
	}
	return
}

func convertSlice2Value(v interface{}, targetType reflect.Type) (newV reflect.Value, ok bool) {
	ok = true

	vValue := reflect.ValueOf(v)
	newV = reflect.MakeSlice(targetType, vValue.Len(), vValue.Len())
	for i := 0; i < vValue.Len(); i++ {
		index := vValue.Index(i)
		value, success := convertParam(index.Interface(), targetType.Elem())
		if !success {
			ok = false
			return
		}
		newV.Index(i).Set(value)
	}
	return
}
