package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

type handshake struct {
	length byte
	pstr   string
	resv   [8]byte
	info   []byte
	peerId []byte
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

		trackerRequest := makeTrackerRequest(torrent)

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
	case "handshake":
		torrentFile := args[0]
		peerData := strings.Split(args[1], ":")
		peerIp, peerPort := peerData[0], peerData[1]
		torrent, err := readTorrentFile(torrentFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		handshakeMsg := makeHandshakeMsg(handshake{
			length: byte(19),
			pstr:   "BitTorrent protocol",
			resv:   [8]byte{},
			info:   torrent.Info.hash(),
			peerId: []byte("00112233445566778899"),
		})
		respHandshake, err := connectWithPeer(peerIp, peerPort, handshakeMsg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Peer ID:", hex.EncodeToString(respHandshake.peerId))
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
