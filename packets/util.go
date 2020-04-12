package packets

import (
	"encoding/binary"
	"os"

	"jadud.com/par2go/verify"
)

func readIntOfSize(f os.File, fileIndex int64, size int32) uint64 {
	// fmt.Printf("readIntOfSize ndx: %d, size: %d\n", fileIndex, size)

	// Read the length.
	valueAsBytes := make([]byte, size)
	_, e2 := f.ReadAt(valueAsBytes, fileIndex)
	verify.Check(e2)

	var value uint64
	if size == 4 {
		value = uint64(binary.LittleEndian.Uint32(valueAsBytes))
	} else if size == 8 {
		value = binary.LittleEndian.Uint64(valueAsBytes)
	}

	// fmt.Printf("\tvalue: %d\n", value)
	return value
}

func readUint32(f os.File, at int64) uint32 {
	return uint32(readIntOfSize(f, at, 4))
}

func readUint64(f os.File, at int64) uint64 {
	return uint64(readIntOfSize(f, at, 8))
}
