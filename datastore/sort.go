package datastore

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type SortType string

const (
	SORT_DESC SortType = "DESC"
	SORT_ASC  SortType = "ASC"
)

var nullBytes = []byte("null")

type Sorts []Sort
type Sort struct {
	Field string
	Type  SortType
}

func (s *Sort) Encode() string {
	if s.Field == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", s.Field, s.Type)
}

func (ss *Sorts) Add(field string, stype SortType) {
	s := Sort{field, stype}
	if field != "" {
		*ss = append(*ss, s)
	}
}

func (ss *Sorts) Encode(fieldPrefix ...string) string {
	sorts := ""
	prefix := ""
	if len(fieldPrefix) > 0 {
		prefix = fieldPrefix[0]
	}
	for _, s := range *ss {
		if s.Field == "" {
			continue
		}

		if len(sorts) > 0 {
			sorts += ","
		}

		sorts += prefix + s.Encode()
	}
	return sorts
}

func (ss *Sorts) ValidEncode(tb *Table, fieldPrefix ...string) string {
	sorts := ""
	prefix := ""
	if len(fieldPrefix) > 0 {
		prefix = fieldPrefix[0]
	}

	for _, s := range *ss {
		if s.Field == "" {
			continue
		}

		if _, ok := tb.GetColumn(s.Field); !ok {
			continue
		}

		if len(sorts) > 0 {
			sorts += ","
		}

		sorts += prefix + s.Encode()
	}
	return sorts
}

func (ss *Sorts) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*ss = nil
		return nil
	}

	var value Sort
	var slice []Sort

	valErr := json.Unmarshal(data, &value)
	slErr := json.Unmarshal(data, &slice)

	if valErr != nil && slErr != nil {
		return slErr
	}

	if valErr == nil {
		*ss = []Sort{value}
		return nil
	}

	*ss = slice
	return nil
}

func (s *Sort) UnmarshalText(text []byte) error {
	stype := SORT_ASC
	sfield := ""
	if len(text) == 0 {
		*s = Sort{
			Field: "",
			Type:  stype,
		}
		return nil
	}

	switch string(text[0]) {
	case "-":
		stype = SORT_DESC
		sfield = string(text[1:])
	case "+":
		stype = SORT_ASC
		sfield = string(text[1:])
	default:
		stype = SORT_ASC
		sfield = string(text)
	}

	*s = Sort{sfield, stype}
	return nil
}

func (s Sort) MarshalText() (text []byte, err error) {
	sprefix := ""
	sfield := s.Field
	if len(s.Field) == 0 {
		return []byte{}, nil
	}

	switch s.Type {
	case SORT_DESC:
		sprefix = "-"
	case SORT_ASC:
		sprefix = "+"
	}

	return []byte(sprefix + sfield), nil
}

func (ss *Sorts) ValueOrZero() []Sort {
	if *ss == nil {
		return []Sort{}
	}
	return *ss
}
