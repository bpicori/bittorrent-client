package client

import "bittorrent-client/config"

const ProtocolIdentifier = "BitTorrent protocol"

type HandShake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   []byte
}

func NewHandShake(infoHash [20]byte) *HandShake {
	return &HandShake{
		Pstr:     ProtocolIdentifier,
		InfoHash: infoHash,
		PeerID:   []byte(config.ClientID),
	}
}

/*
Format parts
1- length of the protocol identifier (BitTorrent protocol)  -> 19
2- Protocol Identifier
3- 8 bytes reserved for extension, for this clients are all "00000000"
4- Info hash of the info 20 bytes
5- PeerId that is generated
*/
func (h *HandShake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}
