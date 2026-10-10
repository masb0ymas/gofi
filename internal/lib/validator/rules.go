package validator

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// rule reports a failure message; ok=false means validation failed. data
// always comes from JSON decoding, so it is never a pointer.
type rule func(key string, data interface{}) (string, bool)

func (v *FieldValidator) Required() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if data == nil || data == "" {
			return fmt.Sprintf("%s is required", key), false
		}
		return "", true
	})
	return v
}

func (v *FieldValidator) String() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if _, ok := data.(string); data != nil && !ok {
			return fmt.Sprintf("%s must be a string", key), false
		}
		return "", true
	})
	return v
}

// MinLen reports a failure when the value is a string shorter than n
// characters.
func (v *FieldValidator) MinLen(n int) *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if s, ok := data.(string); ok && utf8.RuneCountInString(s) < n {
			return fmt.Sprintf("%s must be at least %d characters", key, n), false
		}
		return "", true
	})
	return v
}

func (v *FieldValidator) Regex(pattern string) *FieldValidator {
	re := regexp.MustCompile(pattern)
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if data == nil {
			return "", true
		}
		if s, ok := data.(string); !ok || !re.MatchString(s) {
			return fmt.Sprintf("%s must match the pattern %s", key, re), false
		}
		return "", true
	})
	return v
}

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (v *FieldValidator) Email() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if data == nil {
			return "", true
		}
		if s, ok := data.(string); !ok || !emailPattern.MatchString(s) {
			return fmt.Sprintf("%s must be a valid email address", key), false
		}
		return "", true
	})
	return v
}

func (v *FieldValidator) WithinS(vals ...string) *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if s, ok := data.(string); ok && !slices.Contains(vals, s) {
			return fmt.Sprintf("%s may only contain %s", key, strings.Join(vals, ", ")), false
		}
		return "", true
	})
	return v
}

// Array reports a failure when the value is absent or an empty JSON array.
func (v *FieldValidator) Array() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if data == nil {
			return fmt.Sprintf("%s is required", key), false
		}
		if arr, ok := data.([]interface{}); !ok || len(arr) == 0 {
			return fmt.Sprintf("%s must be a non-empty array", key), false
		}
		return "", true
	})
	return v
}

// Number reports a failure when the value is present and not a JSON number.
func (v *FieldValidator) Number() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		switch data.(type) {
		case nil, float64, float32, int, int32, int64:
			return "", true
		}
		return fmt.Sprintf("%s must be a number", key), false
	})
	return v
}

// Bool reports a failure when the value is present and not a JSON boolean.
func (v *FieldValidator) Bool() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		switch data.(type) {
		case nil, bool:
			return "", true
		}
		return fmt.Sprintf("%s must be a boolean", key), false
	})
	return v
}

// UUID reports a failure when the value is present and not a UUID string.
func (v *FieldValidator) UUID() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if data == nil {
			return "", true
		}
		s, ok := data.(string)
		if !ok {
			return fmt.Sprintf("%s must be a valid UUID", key), false
		}
		if _, err := uuid.Parse(s); err != nil {
			return fmt.Sprintf("%s must be a valid UUID", key), false
		}
		return "", true
	})
	return v
}

// dateFormat is the accepted date encoding, matching encoding/json's
// time.Time unmarshaling so a value that validates also binds.
const dateFormat = time.RFC3339

// Date reports a failure when the value is present and not a parseable date.
func (v *FieldValidator) Date() *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if data == nil {
			return "", true
		}
		s, ok := data.(string)
		if !ok {
			return fmt.Sprintf("%s must be a valid date", key), false
		}
		if _, err := time.Parse(dateFormat, s); err != nil {
			return fmt.Sprintf("%s must be a valid date", key), false
		}
		return "", true
	})
	return v
}
