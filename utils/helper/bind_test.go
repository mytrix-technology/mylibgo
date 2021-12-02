package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/tj/assert"

	utype "github.com/mytrix-technology/mylibgo/utils/types"
	"github.com/mytrix-technology/mylibgo/utils/types/slicer"
)

type TestDTO struct {
	MapField     slicer.SliceOfString `query:"group"`
	CustomerCode string               `query:"customer_code"`
	CityName     slicer.SliceOfString `query:"city_name"`
	ZipcodeID    slicer.SliceOfString `query:"zipcode_id"`
	Zipcode      slicer.SliceOfInt    `query:"zipcode"`
	Detailed     utype.Boolean        `query:"detailed"`
	SortBy       []*Sort              `query:"sort"`
}

func TestBindURLQuery(t *testing.T) {
	input, dto := createTest1()
	testDB := []struct {
		input  url.Values
		expect TestDTO
	}{
		{input, dto},
	}

	for _, test := range testDB {
		var dest TestDTO
		if err := BindURLQuery(&dest, test.input); err != nil {
			t.Errorf(err.Error())
		}

		assert.Equal(t, test.expect, dest)
	}
}
func TestEncodeToURLQuery(t *testing.T) {
	input := TestDTO{
		MapField:     []string{"zipcode", "city_name"},
		CustomerCode: "IDP20120009",
		CityName:     []string{"Jakarta"},
		SortBy:       []*Sort{{"zipcode", SORT_DESC}, {"city_name", SORT_ASC}},
	}

	expected := url.Values{
		"group":         []string{"zipcode", "city_name"},
		"customer_code": []string{"IDP20120009"},
		"city_name":     []string{"Jakarta"},
		"sort":          []string{"-zipcode", "+city_name"},
	}

	res, err := EncodeToURLQuery(input, "query")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expected, res)
}

func createTest1() (url.Values, TestDTO) {
	//q := url.Values{}
	//q.Add("group", "zip_code")
	//q.Add("zipcode", "12860")
	//q.Add("zipcode", "123")
	//q.Add("sort", "field1")
	//q.Add("sort", "-field2")
	//q.Add("detailed", "")

	query := "zipcode=12860&zipcode=123&group=ziska&group=zip_code&sort=field1&sort=-field2&detailed"
	q, err := url.ParseQuery(query)
	if err != nil {
		panic(err)
	}
	dto := TestDTO{
		MapField: []string{"ziska", "zip_code"},
		Zipcode:  []int{12860, 123},
		SortBy:   []*Sort{{"field1", SORT_ASC}, {"field2", SORT_DESC}},
		Detailed: utype.BooleanFrom(true),
	}

	return q, dto
}

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
