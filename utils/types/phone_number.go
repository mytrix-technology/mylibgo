package types

import (
	"github.com/dongri/phonenumber"
	"fmt"
	"database/sql/driver"
	"encoding/json"
)

type UTMobileNumber string

func (p *UTMobileNumber) IsValid() bool {
	str := string(*p)
	number := phonenumber.Parse(str, "ID")
	if number == "" {
		return false
	}

	return true
}

func (p *UTMobileNumber) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	number := phonenumber.Parse(str, "ID")
	if number == "" {
		number = str
	} else {
		number = "0" + number[2:]
	}

	*p = UTMobileNumber(number)
	return nil
}

func (p *UTMobileNumber) Scan(value interface{}) error {
	// if value is nil, false
	if value == nil {
		*p = ""
		return nil
	}

	if bs, err := driver.String.ConvertValue(value); err == nil {
		if err != nil {
			*p = ""
			return nil
		}

		var str string
		switch bs.(type) {
		case string:
			str = bs.(string)
		case []byte:
			str = string(bs.([]byte))
		}

		number := phonenumber.Parse(str, "ID")
		if number == "" {
			number = str
		}else{
			number = "0" + number[2:]
		}
		*p = UTMobileNumber(number)
		return nil
	}

	// otherwise, return an error
	return fmt.Errorf("failed to scan into UTMobileNumber")
}

func (t UTMobileNumber) Value() (driver.Value, error) {
	return string(t), nil
}