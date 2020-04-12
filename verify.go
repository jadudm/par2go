package main

import (
	"fmt"
	"os"
)

// ValidHeader is an 8-byte constant for a valid header.
var ValidHeader = [...]byte{'P', 'A', 'R', '2', '\x00', 'P', 'K', 'T'}

// Check is a cheat.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// CheckForPAR2Header makes sure a file at the given path
// has a valid PAR2 header.
// FIXME: This reads the whole file, and looks at
// the first 8 bytes... should be improved.
func CheckForPAR2Header(file string) (bool, error) {

	// Does the file exist?
	f, err := os.Open(file)
	// If not, send back the failure to open the file.
	if err != nil {
		return false, err
	}

	// Does it have 8 bytes?
	header := make([]byte, 8)
	n, err := f.Read(header)
	// If not, send back the failure to read.
	if err != nil {
		return false, err
	}

	// Check if the bytes read equal a valid header.
	// If so,
	var foundHeader = string(header[:n])
	var result bool = foundHeader == string(ValidHeader[:8])
	if result {
		return result, nil
	}
	var str string = "file header is not PAR2PKT, found %s"
	return result, fmt.Errorf(str, foundHeader)
}

// CheckForFile does stuff.
func CheckForFile(file string) (bool, error) {
	return true, nil
}
