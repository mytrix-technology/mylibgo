package helper

import (
	"encoding"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// BindURLQuery will unmarshal http request query into a struct or map, pointed by dest.
// dest must be a pointer to struct or map
func BindURLQuery(dest interface{}, query url.Values) error {
	return bindData(dest, query, "query")
}

func bindData(ptr interface{}, data map[string][]string, tag string) error {
	if ptr == nil || len(data) == 0 {
		return nil
	}
	typ := reflect.TypeOf(ptr)
	if typ.Kind() != reflect.Ptr {
		return errors.New("destination is not a pointer to struct")
	}
	typ = typ.Elem()
	val := reflect.ValueOf(ptr).Elem()

	// Map
	if typ.Kind() == reflect.Map {
		for k, v := range data {
			val.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
		}
		return nil
	}

	// !struct
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("binding element must be a struct. got %s", typ.Kind().String())
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)

		if inputFieldName == "" {
			inputFieldName = typeField.Name
			// If tag is nil, we inspect if the field is a struct.
			if structFieldKind == reflect.Struct {
				if err := bindData(structField.Addr().Interface(), data, tag); err != nil {
					return err
				}
				continue
			}
			//if _, ok := structField.Addr().Interface().(BindUnmarshaler); !ok && structFieldKind == reflect.Struct {
			//	if err := bindData(structField.Addr().Interface(), data, tag); err != nil {
			//		return err
			//	}
			//	continue
			//}
		}

		rawInputValue, exists := data[inputFieldName]
		if !exists {
			// check again with case insensitive method
			for k, v := range data {
				if strings.EqualFold(k, inputFieldName) {
					rawInputValue = v
					exists = true
					break
				}
			}
		}

		if !exists {
			continue
		}

		//this part is to handle comma separated value
		var inputValue []string
		for _, val := range rawInputValue {
			strSlice := strings.Split(val, ",")
			inputValue = append(inputValue, strSlice...)
		}

		if inputValue == nil {
			continue
		}

		// Call this first, in case we're dealing with an alias to an array type
		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok {
			if err != nil {
				return err
			}
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ {
				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
			return err
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	// But also call it here, in case we're dealing with an array alias
	if ok, err := unmarshalField(valueKind, val, structField); ok {
		return err
	}

	switch valueKind {
	case reflect.Ptr:
		return setWithProperType(structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func unmarshalField(valueKind reflect.Kind, val string, field reflect.Value) (bool, error) {
	switch valueKind {
	case reflect.Ptr:
		return unmarshalFieldPtr(val, field)
	default:
		return unmarshalFieldNonPtr(val, field)
	}
}

func unmarshalFieldNonPtr(value string, field reflect.Value) (bool, error) {
	fieldIValue := field.Addr().Interface()
	if unmarshaler, ok := fieldIValue.(encoding.TextUnmarshaler); ok {
		return true, unmarshaler.UnmarshalText([]byte(value))
	}

	return false, nil
}

func unmarshalFieldPtr(value string, field reflect.Value) (bool, error) {
	if field.IsNil() {
		// Initialize the pointer to a nil value
		field.Set(reflect.New(field.Type().Elem()))
	}
	return unmarshalFieldNonPtr(value, field.Elem())
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func EncodeToURLQuery(ptr interface{}, tag string) (url.Values, error) {
	q := url.Values{}
	return q, encodeData(q, ptr, tag)
}

func encodeData(q url.Values, ptr interface{}, tag string) error {
	if ptr == nil {
		return nil
	}
	typ := reflect.TypeOf(ptr)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	val := reflect.ValueOf(ptr)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Map
	if val.Kind() == reflect.Map {
		if mapKeyIsAllString(val.MapKeys()) {
			return fmt.Errorf("map object should have string keys to be encoded")
		}

		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			v := iter.Value()
			if str, ok := marshalField(v.Type().Kind(), v); ok {
				q.Add(k, str)
				continue
			}
			if v.Kind() == reflect.Slice && v.Len() > 0 {
				for i := 0; i < v.Len(); i++ {
					sval := v.Index(i)
					if str := setToString(v.Elem().Kind(), sval); str != "" {
						q.Add(k, str)
					}
				}
			}
		}

		return nil
	}

	// !struct
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("encoded element must be a struct. got %s", typ.Kind().String())
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)

		if inputFieldName == "" {
			// inputFieldName = typeField.Name
			// If tag is nil, we inspect if the field is a struct.
			if structFieldKind == reflect.Struct {
				if err := encodeData(q, structField.Interface(), tag); err != nil {
					return err
				}
			}
			continue
		}

		// Call this first, in case we're dealing with an alias to an array type
		if str, ok := marshalField(typeField.Type.Kind(), structField); ok {
			if str != "" {
				q.Add(inputFieldName, str)
			}
			continue
		}

		if structFieldKind == reflect.Slice {
			numElems := structField.Len()
			if numElems == 0 {
				continue
			}

			sliceOf := structField.Type().Elem().Kind()
			for j := 0; j < numElems; j++ {
				if str := setToString(sliceOf, structField.Index(j)); str != "" {
					q.Add(inputFieldName, str)
				}
			}
		} else if str := setToString(structField.Kind(), structField); str != "" {
			q.Add(inputFieldName, str)
		}
	}
	return nil
}

func marshalField(valueKind reflect.Kind, value reflect.Value) (string, bool) {
	switch valueKind {
	case reflect.Ptr:
		return marshalFieldPtr(value)
	default:
		return marshalFieldNonPtr(value)
	}
}

func marshalFieldNonPtr(value reflect.Value) (string, bool) {
	fieldIValue := value.Interface()
	if marshaler, ok := fieldIValue.(encoding.TextMarshaler); ok {
		if val, err := marshaler.MarshalText(); err == nil {
			return string(val), true
		}
	}
	return "", false
}

func marshalFieldPtr(value reflect.Value) (string, bool) {
	if value.IsNil() {
		// Initialize the pointer to a nil value
		value.Set(reflect.New(value.Type().Elem()))
	}
	return marshalFieldNonPtr(value.Elem())
}

func setToString(valueKind reflect.Kind, value reflect.Value) string {
	if str, ok := marshalField(valueKind, value); ok {
		return str
	}

	switch value.Kind() {
	case reflect.Ptr:
		return setToString(value.Elem().Kind(), value.Elem())
	case reflect.Int:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Int8:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Int16:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Int32:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint8:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint16:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint32:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.String:
		return value.String()
	}
	return ""
}

func mapKeyIsAllString(keys []reflect.Value) bool {
	allString := true
	for _, v := range keys {
		if v.Kind() != reflect.String {
			allString = false
			break
		}
	}

	return allString
}
