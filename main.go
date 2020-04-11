package main

import (
	"fmt"
	"os"
)

// checkInputFile makes sure the string passed on the
// command line is a file, and a valid PAR2 file to boot.
func checkInputFile(file string) {
	CheckForFile(file)
	CheckForPAR2Header(file)
}

func main() {
	fmt.Println("par2go")

	// Grab the arguments without the application name.
	args := os.Args[1:]
	// TODO: check how things come in. For now, assume one
	// file on the command line.
	file := args[0]
	result, err := CheckForPAR2Header(file)
	Check(err)
	fmt.Println(result)
}
