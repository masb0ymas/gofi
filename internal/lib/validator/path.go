package validator

import "strings"

type path []string

func (p path) key() string {
	return strings.Join(p, ".")
}

func (p path) last() string {
	if len(p) > 0 {
		return p[len(p)-1]
	}
	return ""
}
