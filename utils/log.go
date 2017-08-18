package utils

import "fmt"

// Log ...
func Log(m ...interface{}) {
	fmt.Println(m...)
}

// LogSection ...
func LogSection(title string, m ...interface{}) {
	fmt.Println("\n================================")
	fmt.Printf("==== %s ==== \n", title)
	if len(m) > 0 {
		fmt.Println("================================")
		fmt.Println(m...)
	}
	fmt.Println("================================")
}
