package types

import (
	"encoding/json"
	"time"
)

type RFC3339Time time.Time

func (r RFC3339Time) String() string {
	t := time.Time(r)

	return t.Format(time.RFC3339)
}

func RFCFromTime(t time.Time) RFC3339Time {
	return RFC3339Time(t)
}

func (r RFC3339Time) MarshalJSON() ([]byte, error) {
	t := time.Time(r)

	return json.Marshal(t.Format(time.RFC3339))
}

func (r *RFC3339Time) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	tt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*r = RFCFromTime(tt)

	return nil
}

func (r RFC3339Time) MarshalText() (text []byte, err error) {
	t := time.Time(r)

	return []byte(t.Format(time.RFC3339)), nil
}
