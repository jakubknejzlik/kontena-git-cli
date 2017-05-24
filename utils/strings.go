package utils

import "strings"

// SplitString ...
func SplitString(s, sep string) []string {
	if len(s) == 0 {
		return []string{}
	}
	return strings.Split(strings.Trim(s, " "), sep)
}
