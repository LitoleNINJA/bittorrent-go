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

type trackerRequest struct {
	URL        string
	InfoHash   string
	PeerID     string
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Compact    int
}

type trackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (info Info) hexHash() *bytes.Buffer {
	var b bytes.Buffer
	bencode.Marshal(&b, info)
	sha := sha1.Sum(b.Bytes())
	dst := make([]byte, hex.EncodedLen(len(sha)))
	hex.Encode(dst, sha[:])
	return bytes.NewBuffer(dst)
}

func (info Info) hash() []byte {
	var b bytes.Buffer
	bencode.Marshal(&b, info)
	sha := sha1.Sum(b.Bytes())
	return sha[:]
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

		fmt.Println("Tracker URL:", torrent.Announce)
		fmt.Println("Length:", torrent.Info.Length)
		fmt.Println("Info Hash:", torrent.Info.hexHash())
		fmt.Println("Piece Length:", torrent.Info.PieceLength)
		peiceHashes := hex.EncodeToString([]byte(torrent.Info.Pieces))
		fmt.Println("Piece Hashes:")
		for i := 0; i < len(peiceHashes); i += 40 {
			fmt.Println(peiceHashes[i : i+40])
		}
	case "peers":
		torrentFile := args[0]
		torrent, err := readTorrentFile(torrentFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		trackerRequest := trackerRequest{
			URL:        torrent.Announce,
			InfoHash:   string(torrent.Info.hash()),
			PeerID:     "00112233445566778899",
			Port:       6881,
			Uploaded:   0,
			Downloaded: 0,
			Left:       torrent.Info.Length,
			Compact:    1,
		}

		peers, err := requestPeers(trackerRequest)
		if err != nil {
			fmt.Println(err)
			return
		}

		peerIps := ""
		for i := 0; i < len(peers.Peers); i += 6 {
			ip := fmt.Sprintf("%d.%d.%d.%d", peers.Peers[i], peers.Peers[i+1], peers.Peers[i+2], peers.Peers[i+3])
			port := int(peers.Peers[i+4])<<8 | int(peers.Peers[i+5])
			peerIps += fmt.Sprintf("%s:%d\n", ip, port)
		}
		fmt.Println(peerIps)
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
