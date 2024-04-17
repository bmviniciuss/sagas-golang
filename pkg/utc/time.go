package utc

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

var (
	ISO8601Layout = "2006-01-02T15:04:05.000Z"
	Layouts       = []string{
		ISO8601Layout,
		time.RFC3339,
		time.RFC3339Nano,
	}
)

type Time struct {
	t time.Time
}

func NewFromTime(t time.Time) Time {
	return Time{t: t.UTC()}
}

func NewFromString(s string) (Time, error) {
	for _, layout := range Layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return NewFromTime(t), nil
		}
	}
	return Time{}, fmt.Errorf("failed to parse time: %s", s)
}

func Now() Time {
	return NewFromTime(time.Now())
}

func (t *Time) Time() time.Time {
	return time.Time(t.t).UTC()
}

// formaters

func (t *Time) String() string {
	return time.Time(t.t).UTC().Format(ISO8601Layout)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	parsed, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	tt, err := NewFromString(parsed)
	if err != nil {
		return err
	}
	*t = tt
	return nil
}

func (t Time) Value() (driver.Value, error) {
	return t.Time(), nil
}

func (t *Time) Scan(value interface{}) error {
	switch src := value.(type) {
	case time.Time:
		*t = NewFromTime(src)
		return nil
	case string:
		tt, err := NewFromString(src)
		if err != nil {
			return err
		}
		*t = tt
		return nil
	case []byte:
		tt, err := NewFromString(string(src))
		if err != nil {
			return err
		}
		*t = tt
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}
