package datetime

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gopkg.in/guregu/null.v4/zero"
	"reflect"
	"time"
)

var UTJSONDateTimeFormat = "2006-01-02 15:04:05"

// UTJSONDateTime extending zero.Time type with custom JSON Marshalling
// into format 'YYYY-MM-DD HH:MM:SS'
type UTJSONDateTime struct {
	zero.Time
}

func NewUTJSONDateTime(t time.Time) UTJSONDateTime {
	return UTJSONDateTime{zero.NewTime(t, true)}
}

func (t *UTJSONDateTime) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		var date time.Time
		if x != "" {
			date, err = time.Parse(UTJSONDateTimeFormat, x)
			if err != nil {
				return fmt.Errorf("failed to parse time %s", x)
			}
		}

		*t = UTJSONDateTime{zero.TimeFrom(date)}
		return nil
	case nil:
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("json: cannot unmarshal %v into Go value of type UTJSONDateTime", reflect.TypeOf(v).Name())
	}
}

func (t UTJSONDateTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		//return (time.Time{}).MarshalJSON()
		return json.Marshal("")
	}
	//return t.Time.MarshalJSON()
	s := t.Time.Time
	return json.Marshal(s.Format(UTJSONDateTimeFormat))
}

func (t UTJSONDateTime) TimeNow() time.Time {
	return time.Now()
}

func (t UTJSONDateTime) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

func (t UTJSONDateTime) ToString() string {
	return t.Time.Time.Format(UTJSONDateTimeFormat)
}

func (t UTJSONDateTime) Value() (driver.Value, error) {
	if t.Valid {
		return t.ToString(), nil
	} else {
		return nil, nil
	}
}
