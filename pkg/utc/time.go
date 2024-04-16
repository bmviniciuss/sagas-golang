package utc

import (
	"strconv"
	"time"
)

var ISO8601Layout = "2006-01-02T15:04:05.000Z"

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Time(t).UTC().Format(ISO8601Layout))), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	parsed, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	tt, err := time.Parse(ISO8601Layout, parsed)
	if err != nil {
		return err
	}
	*t = Time(tt)
	return nil
}

func (t *Time) String() string {
	return time.Time(*t).UTC().Format(ISO8601Layout)
}

func NewFromTime(t time.Time) Time {
	return Time(t.UTC())
}
