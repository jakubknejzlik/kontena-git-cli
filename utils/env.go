package utils

import (
	"fmt"
	"os"
)

// GetenvStrict return environment variable for key, panic if value is not specified
func GetenvStrict(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("env %s is not specified", key))
	}
	return value
}

// GetenvStrictWithTip return environment variable for key, panic with tip if value is not specified
func GetenvStrictWithTip(key, tip string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("env %s is not specified. %s", key, tip))
	}
	return value
}

// Getenv return environment variable for key, panic if value is not specified
func Getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
