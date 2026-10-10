package validator

import (
	"fmt"
	"regexp"
	"unicode/utf8"
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

// MaxLen reports a failure when the value is a string longer than n
// characters.
func (v *FieldValidator) MaxLen(n int) *FieldValidator {
	v.registerRule(func(key string, data interface{}) (string, bool) {
		if s, ok := data.(string); ok && utf8.RuneCountInString(s) > n {
			return fmt.Sprintf("%s must be at most %d characters", key, n), false
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
