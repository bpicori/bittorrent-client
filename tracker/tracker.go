package tracker

import (
	"bittorrent-client/torrentfile"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/IncSW/go-bencode"
)

const PORT = 6881

type Peer struct {
	ID   []byte
	IP   net.IP
	Port int64
}

func GetPeers(file torrentfile.TorrentFile) ([]Peer, error) {
	peers := []Peer{}
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
		p := Peer{}
		p.ID = eMap["peer id"].([]byte)
		p.Port = eMap["port"].(int64)
		p.IP = net.IP(eMap["ip"].([]byte))

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
		"peer_id":    []string{generateClientId()},
		"port":       []string{strconv.Itoa(int(PORT))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(0)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func generateClientId() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 20)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
