package lib

import "slices"

func Contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
