package main

import (
	"fmt"
	"os"
)

var major, minor, patch string

func showVersion() {
	fmt.Printf("%s.%s.%s", major, minor, patch)
	os.Exit(0)
}
