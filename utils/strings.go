package utils

import "strings"

// SplitString ...
func SplitString(s, sep string) []string {
	if len(s) == 0 {
		return []string{}
	}
	return strings.Split(strings.Trim(strings.Trim(s, " "), "\n"), sep)
}

// ArrayOfStringsContains ...
func ArrayOfStringsContains(items []string, value string) bool {
	for _, x := range items {
		if x == value {
			return true
		}
	}
	return false
}
