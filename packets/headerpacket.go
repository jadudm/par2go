package packets

import (
	"fmt"
	"os"

	"jadud.com/par2go/verify"
)

/*
	Bytes		Type				Description
	8				byte[8]			Magic sequence. Used to quickly identify location of packets. Value = {'P', 'A', 'R', '2', '\0', 'P', 'K', 'T'} (ASCII)
	8				8-byte uint	Length of the entire packet. Must be multiple of 4. (NB: Includes length of header.)
	16			MD5 Hash		MD5 Hash of packet. Used as a checksum for the packet. Calculation starts at first byte of Recovery Set ID and ends at last byte of body. Does not include the magic sequence, length field or this field. NB: The MD5 Hash, by its definition, includes the length as if it were appended to the packet.
	16			MD5 Hash		Recovery Set ID. All packets that belong together have the same recovery set ID. (See "main packet" for how it is calculated.)
	16			byte[16]		Type. Can be anything. All beginning "PAR " (ASCII) are reserved for specification-defined packets. Application-specific packets are recommended to begin with the ASCII name of the client.
	?*4			?						Body of Packet. Must be a multiple of 4 bytes.
*/

// PacketHeader is at the front of every PAR2 packet.
// PAR2 files are made up of sequences of packets; therefore,
// this header is used everywhere.
type PacketHeader struct {
	_offset uint64  // Offset into file where this header was found.
	_type   PktType // 8-byte magic sequence. MCJ: Converted to enum.
	// MD5 Hash of packet.
	// Used as a checksum for the packet.
	// Calculation starts at first byte of Recovery Set ID and
	// ends at last byte of body. Does not include the magic sequence,
	// length field or this field.
	// NB: The MD5 Hash, by its definition, includes the length
	// as if it were appended to the packet.
	packetLength  uint64
	checksum      HashMD5
	recoverySetID HashMD5
	packetType    [16]byte
}

// -------------- INTERFACE PacketBody -------------- //
func (ph PacketHeader) getType() PktType {
	return ph._type
}

func (ph PacketHeader) toBytes() []byte {
	// FIXME: Not a proper implementation
	b := make([]byte, 1)
	return b
}

func (ph PacketHeader) getLength() uint64 {
	return 64
}

func checkForMagicSequence(f os.File, fileIndex int64) {
	tag := make([]byte, 8)
	_, e1 := f.ReadAt(tag, fileIndex)
	verify.Check(e1)

	if !(string(tag[0:8]) == string(HeaderMagic[0:8])) {
		panic(fmt.Sprintf(
			"invalid packet tag at %d: %s is not %s",
			fileIndex, tag, HeaderMagic,
		))
	}
}

// ReadPacketHeader ...
func readPacketHeader(f os.File, at int64) PacketHeader {
	// Make sure this is a PAR2PKT
	checkForMagicSequence(f, at)

	// Read the length of the packet, including this header.
	packetLength := readIntOfSize(f, at+8, 8)
	// fmt.Println("packet length", packetLength)

	// Read the full packet
	fullPacket := make([]byte, packetLength)
	// Read from the start index; the length includes this header.
	_, e1 := f.ReadAt(fullPacket, at)
	verify.Check(e1)

	// Read in the two MD5 hashes. Each 16 bytes.
	// I'm already 16 bytes in.
	var checksum [16]byte
	var recoverySetID [16]byte
	copy(checksum[:], fullPacket[16:32])
	copy(recoverySetID[:], fullPacket[32:48])

	// Get the packet type.
	var packetType [16]byte
	copy(packetType[:], fullPacket[48:64])

	return PacketHeader{
		_type:         PAR2PKT,
		_offset:       uint64(at),
		packetLength:  packetLength,
		checksum:      checksum,
		recoverySetID: recoverySetID,
		packetType:    packetType,
	}
}
