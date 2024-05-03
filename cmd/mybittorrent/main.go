package main

import (
	"encoding/json"
	"fmt"
	"os"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

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
		data, err := readTorrentFile(torrentFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		url := data["announce"].(string)
		length := data["info"].(map[string]interface{})["length"].(int)
		fmt.Printf("Tracker URL: %s\nLength: %d\n", url, int(length))
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
