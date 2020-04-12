package main

import (
	"fmt"
	"os"

	"jadud.com/par2go/packets"
	"jadud.com/par2go/verify"
)

// checkInputFile makes sure the string passed on the
// command line is a file, and a valid PAR2 file to boot.
func checkInputFile(file string) {
	_, e1 := verify.CheckForFile(file)
	verify.Check(e1)
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
	verify.Check(err)
	header, body := packets.ReadPacket(*f, 0)
	fmt.Println(header)
	fmt.Println(body)
}
