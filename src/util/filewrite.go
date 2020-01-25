package util

import (
	"fmt"
	"os"
)

// FileWrite Helper Utility:
// Creates a file if it doesn't already exist
// Appends string to the file
func FileWrite(msg, outputFile string) {
	f, err := os.OpenFile(outputFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := f.WriteString(msg + "\n"); err != nil {
		fmt.Println(err)
	}
	if err := f.Sync(); err != nil {
		fmt.Println(err)
	}
	f.Close()
}

// LogWrite Helper Utility:
// Prints string to stdout and appends to file
// File will be created if it doesn't exist
func LogWrite(msg, outputFile string) {
	fmt.Println(msg)
	f, err := os.OpenFile(outputFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	if _, err := f.WriteString(msg + "\n"); err != nil {
		return
	}
	if err := f.Sync(); err != nil {
		return
	}
}
