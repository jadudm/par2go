package packets

import (
	"fmt"
	"os"
)

// HeaderMagic is an 8-byte constant for a valid packet header.
var HeaderMagic = [...]byte{'P', 'A', 'R', '2', '\x00', 'P', 'K', 'T'}

// MainMagic FIXME
var MainMagic = [...]byte{'P', 'A', 'R', ' ', '2', '.', '0', '\x00', 'M', 'a', 'i', 'n', '\x00', '\x00', '\x00', '\x00'}

// PktType provides an enumeration type for PAR2 packets.
type PktType int

// PAR2 packets are enumated for clarity.
const (
	PAR2PKT PktType = iota
	MAINPKT
)

// Packet FIXME
type Packet struct {
	header PacketHeader
	body   PacketBody
}

// PacketBody FIXME
type PacketBody interface {
	getType() PktType
	getLength() uint64
	toBytes() []byte
}

/*
http://parchive.sourceforge.net/docs/specifications/parity-volume-spec/article-spec.html
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

// HashMD5 is a type for 16-byte MD5 hashes
// found throughout the PAR2 format.
type HashMD5 [16]byte

// ReadPacket FIXME
func ReadPacket(f os.File, at int64) (PacketHeader, PacketBody) {
	ph := readPacketHeader(f, at)
	// I read the header at location 'at'.
	body := readPacketBody(f, ph)

	return ph, body
}

func readPacketBody(f os.File, ph PacketHeader) PacketBody {
	// Headers are 64 bytes.
	start := ph._offset + ph.getLength()

	// Read in the packet body from the given offset
	// and based on the packet type encoded in the header.
	switch ph.packetType {
	case MainMagic:
		return readMainPacket(f, int64(start))
	default:
		fmt.Println("UKNOWN PKT")
	}
	return nil
}
