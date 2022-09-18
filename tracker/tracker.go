package tracker

import (
	"bittorrent-client/config"
	"bittorrent-client/peer"
	"bittorrent-client/torrentfile"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/IncSW/go-bencode"
)

func GetPeers(file torrentfile.TorrentFile) ([]peer.Peer, error) {
	peers := []peer.Peer{}
	c := &http.Client{Timeout: 15 * time.Second}

	url, err := trackerUrl(file)
	if err != nil {
		return peers, err
	}
	respRaw, err := c.Get(url)
	if err != nil {
		panic(err)
	}

	resp, _ := io.ReadAll(respRaw.Body)

	raw, err := bencode.Unmarshal(resp)
	if err != nil {
		panic(err)
	}

	rawMap := raw.(map[string]interface{})

	peersRaw := rawMap["peers"].([]interface{})
	for _, e := range peersRaw {
		eMap := e.(map[string]interface{})

		p := peer.Peer{}
		p.ID = eMap["peer id"].([]byte)
		p.Port = eMap["port"].(int64)
		p.IP = net.ParseIP(string(eMap["ip"].([]byte)))

		peers = append(peers, p)
	}

	return peers, nil
}

func trackerUrl(file torrentfile.TorrentFile) (string, error) {
	base, err := url.Parse(file.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(file.Info.Hash[:])},
		"peer_id":    []string{config.ClientID},
		"port":       []string{strconv.Itoa(int(config.Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(0)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
