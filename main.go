package main

import (
	"bittorrent-client/torrentfile"
	"bittorrent-client/tracker"
	"fmt"
)

func main() {
	torrentFile, err := torrentfile.ParseFile("./ubuntu.torrent")
	if err != nil {
		panic(err)
	}

	peers, err := tracker.GetPeers(torrentFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(peers)

}
