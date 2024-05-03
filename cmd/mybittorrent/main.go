package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"info"`
}

type Info struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

func (info Info) hash() *bytes.Buffer {
	var b bytes.Buffer
	bencode.Marshal(&b, info)
	sha := sha1.Sum(b.Bytes())
	dst := make([]byte, hex.EncodedLen(len(sha)))
	hex.Encode(dst, sha[:])
	return bytes.NewBuffer(dst)
}

func main() {
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "decode":
		bencodedValue := args[0]

		decoded, _, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	case "info":
		torrentFile := args[0]
		torrent, err := readTorrentFile(torrentFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Tracker URL: ", torrent.Announce)
		fmt.Println("Length: ", torrent.Info.Length)
		fmt.Println("Info Hash: ", torrent.Info.hash())
		fmt.Println("Piece Length: ", torrent.Info.PieceLength)
		peiceHashes := hex.EncodeToString([]byte(torrent.Info.Pieces))
		fmt.Println("Piece Hashes: ")
		for i := 0; i < len(peiceHashes); i += 40 {
			fmt.Println(peiceHashes[i : i+40])
		}
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
