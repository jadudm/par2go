package main

import (
	"fmt"
	"os"
)

// checkInputFile makes sure the string passed on the
// command line is a file, and a valid PAR2 file to boot.
func checkInputFile(file string) {
	_, e1 := CheckForFile(file)
	Check(e1)
	_, e2 := CheckForPAR2Header(file)
	Check(e2)
}

func main() {
	fmt.Println("par2go")

	// FIXME: Check length of arguments list.
	// Grab the arguments without the application name.
	args := os.Args[1:]
	// TODO: check how things come in. For now, assume one
	// file on the command line.
	file := args[0]
	fmt.Printf("%s OK\n", file)
	f, err := os.Open(file)
	Check(err)
	pkt := ReadPacketHeader(*f, 0)
	fmt.Println(pkt)
}
