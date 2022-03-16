package dto

import "time"

type DateTime time.Time

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(dt).Format(`"2006-01-02T15:04:05Z"`)), nil
}
