package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// PktType provides an enumeration type for PAR2 packets.
type PktType int

// PAR2 packets are enumated for clarity.
const (
	PAR2PKT PktType = iota
)

/*
	A PAR 2.0 file consists of a sequence of "packets".
	A packet has a fixed sized header and a variable length body.
	The packet header contains a checksum for the packet -
	if the packet is damaged, the packet is ignored.
	The packet header also contains a packet-type.
	If the client does not understand the packet type, the packet is ignored.
	To be compliant with this specification, a client must understand the
	"core" set of packets. Client may process the optional packets or
	create their own application-specific packets.
*/

/*
	Bytes		Type				Description
	8				byte[8]			Magic sequence. Used to quickly identify location of packets. Value = {'P', 'A', 'R', '2', '\0', 'P', 'K', 'T'} (ASCII)
	8				8-byte uint	Length of the entire packet. Must be multiple of 4. (NB: Includes length of header.)
	16			MD5 Hash		MD5 Hash of packet. Used as a checksum for the packet. Calculation starts at first byte of Recovery Set ID and ends at last byte of body. Does not include the magic sequence, length field or this field. NB: The MD5 Hash, by its definition, includes the length as if it were appended to the packet.
	16			MD5 Hash		Recovery Set ID. All packets that belong together have the same recovery set ID. (See "main packet" for how it is calculated.)
	16			byte[16]		Type. Can be anything. All beginning "PAR " (ASCII) are reserved for specification-defined packets. Application-specific packets are recommended to begin with the ASCII name of the client.
	?*4			?						Body of Packet. Must be a multiple of 4 bytes.
*/

// HashMD5 is a wrapper struct for 16-byte MD5 hashes
// found throughout the PAR2 format.
type HashMD5 struct {
	value [16]byte
}

// PacketHeader is at the front of every PAR2 packet.
// PAR2 files are made up of sequences of packets; therefore,
// this header is used everywhere.
type PacketHeader struct {
	_type  PktType // 8-byte magic sequence. MCJ: Converted to enum.
	length uint64  // 8-byte uint. MCJ: Converted to uint64.
	// MD5 Hash of packet.
	// Used as a checksum for the packet.
	// Calculation starts at first byte of Recovery Set ID and
	// ends at last byte of body. Does not include the magic sequence,
	// length field or this field.
	// NB: The MD5 Hash, by its definition, includes the length
	// as if it were appended to the packet.
	checksum      HashMD5
	recoverySetID HashMD5
	packetType    [16]byte
	// NOTE: The 'body' of the packet must be length of mod 4.
	// FIXME: There should be a constructor for this struct
	// that provides guarantees. Or, some other check.
	body string
}

func arraysAreSame(a []byte, b []byte) bool {
	return string(a) == string(b)
}

func readAndCheckHeader(f os.File, fileIndex int64) []byte {
	tag := make([]byte, 8)
	_, e1 := f.ReadAt(tag, fileIndex)
	Check(e1)

	if !arraysAreSame(tag[0:8], ValidHeader[0:8]) {
		panic(fmt.Sprintf(
			"invalid packet tag at %d: %s is not %s",
			fileIndex, tag, ValidHeader,
		))
	}

	return tag
}

func readPacketLength(f os.File, fileIndex int64) uint64 {
	// Read the length.
	lengthBytes := make([]byte, 8)
	_, e2 := f.ReadAt(lengthBytes, fileIndex)
	Check(e2)
	packetLength := binary.LittleEndian.Uint64(lengthBytes)
	return packetLength
}

// ReadPacketHeader ...
func ReadPacketHeader(f os.File, fileIndex int64) PacketHeader {
	readAndCheckHeader(f, fileIndex)
	packetLength := readPacketLength(f, fileIndex+8)
	fmt.Println("packet length", packetLength)
	// Read the full packet, but remember that the 16-byte
	// header is part of the length.
	fullPacket := make([]byte, packetLength)
	_, e1 := f.ReadAt(fullPacket, fileIndex)
	Check(e1)

	// Read in the two MD5 hashes. Each 16 bytes.
	// We're already 16 bytes in.
	checksum := HashMD5{}
	recoverySetID := HashMD5{}
	copy(checksum.value[:], fullPacket[16:32])
	copy(recoverySetID.value[:], fullPacket[32:48])

	return PacketHeader{
		_type:         PAR2PKT,
		length:        packetLength,
		checksum:      checksum,
		recoverySetID: recoverySetID,
		body:          string(fullPacket[48:]),
	}
}
