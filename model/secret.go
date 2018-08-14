package model

import (
	"regexp"
	"time"
)

// Secret ...
type Secret struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SecretParseList ...
func SecretParseList(rows []string) ([]Secret, error) {
	valueRE := regexp.MustCompile("[^\\s]+\\s+([^\\s]+)")
	result := []Secret{}
	rows = rows[1:]
	for len(rows) > 2 {
		name := valueRE.FindStringSubmatch(rows[0])[1]
		createdAt := valueRE.FindStringSubmatch(rows[1])[1]
		updatedAt := valueRE.FindStringSubmatch(rows[2])[1]
		createdAtTime, _ := time.Parse(time.RFC3339Nano, createdAt)
		updatedAtTime, _ := time.Parse(time.RFC3339Nano, updatedAt)
		result = append(result, Secret{name, createdAtTime, updatedAtTime})
		rows = rows[3:]
	}
	return result, nil
}
