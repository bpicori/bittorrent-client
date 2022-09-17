package torrentfile

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/IncSW/go-bencode"
)

type Files struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"`
}

type Info struct {
	Name        string `bencode:"name"`
	Length      int64  `bencode:"length"`
	Files       []Files
	PieceLength int64    `bencode:"piece length"`
	Pieces      [][]byte `bencode:"pieces"`
	Hash        [20]byte
}

type TorrentFile struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"info"`
}

func ParseFile(path string) (TorrentFile, error) {
	torrentFile := TorrentFile{}
	data, err := os.ReadFile(path)
	if err != nil {
		return torrentFile, err
	}

	raw, err := bencode.Unmarshal(data)
	if err != nil {
		return torrentFile, err
	}
	rawMap := raw.(map[string]interface{})

	announce, err := parseAnnounce(rawMap)
	if err != nil {
		return torrentFile, err
	}

	info, err := parseInfo(rawMap)
	if err != nil {
		return torrentFile, err
	}

	torrentFile.Announce = announce
	torrentFile.Info = info

	return torrentFile, nil
}

func parseInfo(raw map[string]interface{}) (Info, error) {
	info := Info{}

	infoMap := raw["info"].(map[string]interface{})

	info.Name = string(infoMap["name"].([]byte))
	info.Length = infoMap["length"].(int64)
	info.PieceLength = infoMap["piece length"].(int64)

	piecesBuffer := infoMap["pieces"].([]byte)

	hashLength := 20
	nrOfPieces := len(piecesBuffer) / hashLength

	var res [][]byte = [][]byte{}
	for i := 0; i < nrOfPieces; i++ {
		from := i * 20
		to := (i + 1) * 20
		tmp := piecesBuffer[from:to]
		res = append(res, tmp)
	}

	info.Pieces = res
	t, err := bencode.Marshal(infoMap)
	if err != nil {
		return info, err
	}
	info.Hash = sha1.Sum(t)

	return info, nil

}

func parseAnnounce(raw map[string]interface{}) (string, error) {
	announce, ok := raw["announce"].([]byte)
	if !ok {
		return "", fmt.Errorf("couldn't parse announce to string, value %v", raw["announce"])
	}
	return string(announce), nil
}

// func parseError(key string) error {
// 	return fmt.Errorf(msg)
// }
