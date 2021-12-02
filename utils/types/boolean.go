package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/guregu/null.v4"
)

var nullBytes = []byte("null")

type Boolean struct {
	null.Bool
}

func BooleanFrom(b bool) Boolean {
	return Boolean{null.BoolFrom(b)}
}

func (b *Boolean) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		b.Valid = false
		return nil
	}

	var str string
	var num int
	var raw null.Bool

	strErr := json.Unmarshal(data, &str)
	numErr := json.Unmarshal(data, &num)
	rawErr := json.Unmarshal(data, &raw)

	if strErr != nil && rawErr != nil && numErr != nil {
		fmt.Println("nothing works")
		return rawErr
	}

	if rawErr == nil {
		b.Bool = raw
		return nil
	}

	if numErr == nil {
		str = strconv.Itoa(num)
	}

	if strings.EqualFold(str, "FALSE") || str == "0" {
		*b = BooleanFrom(false)
	} else if strings.EqualFold(str, "TRUE") || str == "1" {
		*b = BooleanFrom(true)
	} else {
		return fmt.Errorf("cannot unmarshal the value %q into boolean type", str)
	}

	return nil
}

func (b *Boolean) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*b = BooleanFrom(true)
		return nil
	}

	str := string(text)
	if strings.EqualFold(str, "FALSE") || str == "0" {
		*b = BooleanFrom(false)
	} else if strings.EqualFold(str, "TRUE") || str == "1" {
		*b = BooleanFrom(true)
	} else {
		return fmt.Errorf("cannot unmarshal the value %q into boolean type", str)
	}

	return nil
}

func (b *Boolean) MarshalText() (text []byte, err error) {
	if !b.Valid {
		return []byte{}, nil
	}

	if b.ValueOrZero() {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}
