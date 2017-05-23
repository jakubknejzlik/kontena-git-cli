package utils

import "fmt"

// Log ...
func Log(m ...interface{}) {
	fmt.Println(m...)
}

// LogSection ...
func LogSection(title string, m ...interface{}) {
	fmt.Printf("==== %s ==== \n", title)
	fmt.Println("================================")
	fmt.Println(m...)
	fmt.Println("================================")
}
