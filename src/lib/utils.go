package lib

import (
	"strconv"
	"strings"
)

// IntToString converts an integer to a string
func IntToString(n int) string {
	return strconv.Itoa(n)
}

// JoinStrings joins a slice of strings with a separator
func JoinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}
