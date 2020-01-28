package util

import (
	"fmt"
	"log"
	"os"

	"github.com/HDN-1D10T/divinity/src/config"
)

// Configuration imported from src/config
type Configuration struct{ config.Options }

var (
	conf       = Configuration{config.ParseConfiguration()}
	outputFile = *conf.OutputFile
)

// FileWrite Helper Utility:
// Creates a file if it doesn't already exist
// Appends string to the file
func FileWrite(msg string) {
	if len(outputFile) > 0 {
		f, err := os.OpenFile(outputFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := f.WriteString(msg); err != nil {
			log.Println(err)
			f.Close()
			return
		}
		if err := f.Sync(); err != nil {
			log.Println(err)
			f.Close()
			return
		}
		f.Close()
	}
}

// LogWrite Helper Utility:
// Prints string to stdout and appends to file
// File will be created if it doesn't exist
func LogWrite(msg string) {
	if len(outputFile) > 0 {
		fmt.Println(msg)
		f, err := os.OpenFile(outputFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := f.WriteString(msg + "\n"); err != nil {
			log.Println(err)
			f.Close()
			return
		}
		if err := f.Sync(); err != nil {
			log.Println(err)
			f.Close()
			return
		}
		f.Close()
		return
	}
	fmt.Println(msg)
	return
}
