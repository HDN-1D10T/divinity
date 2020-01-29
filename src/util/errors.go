package util

import (
	"fmt"
	"log"
)

// PanicErr logs the line and sends an Panic if err != nil
func PanicErr(e error) {
	if e != nil {
		log.Panicln(e)
	}
}

// LogErr simply logs the error to stderr
func LogErr(e error) {
	if e != nil {
		log.Println(e)
	}
}

// PrintErr prints only the error to stdout
func PrintErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
