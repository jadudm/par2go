package packets

import (
	"os"

	"jadud.com/par2go/verify"
)

// MainPacket FIXME
type MainPacket struct {
	_offset               uint64
	_type                 PktType
	sliceSize             uint64 // 8-byte uint
	numFilesInRecoverySet uint32 // 4-byte uint
	recoverySetFileIDs    map[int][]byte
	nonRecoverySetFileIDs map[int][]byte
}

// -------------- INTERFACE PacketBody -------------- //
func (mp MainPacket) getType() PktType {
	return mp._type
}

func (mp MainPacket) toBytes() []byte {
	// FIXME: Not a proper implementation
	b := make([]byte, 1)
	return b
}

func (mp MainPacket) getLength() uint64 {
	return 64 //FIXME
}

func readMainPacket(f os.File, at int64) MainPacket {
	// Read the slice size.
	sliceSize := readUint64(f, at)  // 8 Bytes
	numFiles := readUint32(f, at+8) // 4 bytes
	// I am now 12 bytes in.
	fileIndex := at + 12

	// Set up maps of IDs
	recoverySetFileIDs := make(map[int][]byte)
	nonRecoverySetFileIDs := make(map[int][]byte)

	// Read the file IDs of everything in the recovery set.
	for ndx := 0; ndx < int(numFiles); ndx++ {
		id := make([]byte, 16)
		_, e := f.ReadAt(id, fileIndex)
		verify.Check(e)
		recoverySetFileIDs[ndx] = id
		fileIndex += 16
	}

	// 20200412 MCJ WARNING
	// The docs do not make clear how to know how many
	// files are *not* in the recovery set. Therefore,
	// I have no idea how to read in the MD5 hash array.

	return MainPacket{
		_type:                 MAINPKT,
		_offset:               uint64(at),
		sliceSize:             sliceSize,
		numFilesInRecoverySet: numFiles,
		recoverySetFileIDs:    recoverySetFileIDs,
		nonRecoverySetFileIDs: nonRecoverySetFileIDs,
	}
}
