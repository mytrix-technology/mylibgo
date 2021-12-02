package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"reflect"
)

// SerializeToBase64 will encode arbitrary value into gob bytes, and then encode the bytes into base64 string.
// If the value is a struct with an interface field, and you put another struct as the interface,
// you need to put the default value of the struct into the registers field
func SerializeToBase64(value interface{}, registers ...interface{}) (string, error) {
	var buf = new(bytes.Buffer)
	if len(registers) > 0 {
		for _, obj := range registers {
			gob.Register(obj)
		}
	}

	if err := gob.NewEncoder(buf).Encode(value); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DeserializeFromBase64(value string, dst interface{}, registers ...interface{}) error {
	dstTyp := reflect.TypeOf(dst)
	if dstTyp.Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}

	if len(registers) > 0 {
		for _, obj := range registers {
			gob.Register(obj)
		}
	}

	valBytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(valBytes)
	return gob.NewDecoder(buf).Decode(dst)
}
