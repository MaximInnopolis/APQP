package utils

import "time"

// CustomTime is a custom time structure
type CustomTime struct {
	time.Time
}

// MarshalJSON - implementation of the json.Marshaler interface
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte(`"Time is not defined yet"`), nil
	}
	return []byte(`"` + ct.Format("2006-01-02T15:04:05.9999999-07:00") + `"`), nil
}
