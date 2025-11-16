package lib

import (
	"errors"
	"strconv"
	"time"
)

type StringTo struct{}

func (s *StringTo) ToInt32(str string, defaultValue ...int32) int32 {
	def := int32(0)
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

	v, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return def
	}
	return int32(v)
}

func (s *StringTo) ToInt64(str string, defaultValue ...int64) int64 {
	def := int64(0)
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return def
	}
	return v
}

func (s *StringTo) ToTime(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, str); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid time format")
}
