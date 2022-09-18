package peer

import (
	"net"
	"strconv"
)

type Peer struct {
	ID   []byte
	IP   net.IP
	Port int64
}

func (p *Peer) HostPortURL() string {
	return net.JoinHostPort(p.IP.String(), strconv.FormatInt(p.Port, 10))
}
