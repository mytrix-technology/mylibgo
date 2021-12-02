package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Sequence int

func (s *Sequence) UnmarshalJSON(data []byte) error {
	var str string
	var num int

	strErr := json.Unmarshal(data, &str)
	numErr := json.Unmarshal(data, &num)

	if strErr != nil && numErr != nil {
		return numErr
	}

	if strErr == nil {
		if len(str) > 2 {
			suffix := str[len(str)-2:]
			switch {
			case strings.EqualFold(suffix, "st"):
				if str[len(str)-3:len(str)-2] != "1" {
					return fmt.Errorf("invalid sequence string: %s", str)
				}
				str = str[:len(str)-2]
			case strings.EqualFold(suffix, "nd"):
				if str[len(str)-3:len(str)-2] != "2" {
					return fmt.Errorf("invalid sequence string: %s", str)
				}
				str = str[:len(str)-2]
			case strings.EqualFold(suffix, "rd"):
				if str[len(str)-3:len(str)-2] != "3" {
					return fmt.Errorf("invalid sequence string: %s", str)
				}
				str = str[:len(str)-2]
			case strings.EqualFold(suffix, "th"):
				str = str[:len(str)-2]
			}
		}

		num, numErr = strconv.Atoi(str)
		if numErr != nil {
			return fmt.Errorf("cannot convert sequence string %s into sequence number. %s", data, numErr)
		}
	}

	if num < 1 {
		return fmt.Errorf("sequence start from 1 or 1st upward. got %d", num)
	}

	*s = Sequence(num)

	return nil
}

func (s *Sequence) String() string {
	suffix := "th"
	str := strconv.Itoa(int(*s))
	switch str[len(str)-1:] {
	case "1":
		suffix = "st"
	case "2":
		suffix = "nd"
	case "3":
		suffix = "rd"
	}

	return str + suffix
}

func (s *Sequence) Valid() bool {
	return *s > 0
}
