package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func main() {
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "decode":
		bencodedValue := args[0]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		return decodeString(bencodedString)
	} else if bencodedString[0] == 'i' && bencodedString[len(bencodedString)-1] == 'e' {
		return decodeInt(bencodedString)
	} else {
		return "", fmt.Errorf("invalid bencoded string")
	}
}

func decodeString(bencodedString string) (string, error) {
	firstColonIndex := strings.Index(bencodedString, ":")
	lengthStr := bencodedString[:firstColonIndex]
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}
	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
}

func decodeInt(bencodedString string) (int, error) {
	return strconv.Atoi(bencodedString[1 : len(bencodedString)-1])
}
