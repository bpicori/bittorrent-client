package main

import (
	"bittorrent-client/client"
	"bittorrent-client/config"
	"bittorrent-client/torrentfile"
	"bittorrent-client/tracker"
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	// init config
	config.Init()

	// parse file
	torrentFile, err := torrentfile.ParseFile("./ubuntu.torrent")
	if err != nil {
		panic(err)
	}

	peers, err := tracker.GetPeers(torrentFile)
	if err != nil {
		panic(err)
	}

	// first peer
	peer := peers[0]
	conn, err := net.DialTimeout("tcp", peer.HostPortURL(), 3*time.Second)
	if err != nil {
		panic(err)
	}
	handshake := client.NewHandShake(torrentFile.Info.Hash)
	fmt.Println(handshake.Serialize())

	var buf bytes.Buffer
	io.Copy(&buf, conn)
	fmt.Println("total size:", buf.Len())

}
