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
	} else if bencodedString[0] == 'l' && bencodedString[len(bencodedString)-1] == 'e' {
		return decodeList(bencodedString)
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

func decodeList(bencodedString string) ([]interface{}, error) {
	list := []interface{}{}
	bencodedString = bencodedString[1 : len(bencodedString)-1]

	for len(bencodedString) > 0 {
		if bencodedString[0] == 'i' {
			endIndex := strings.Index(bencodedString, "e")
			intValue, err := decodeInt(bencodedString[:endIndex+1])
			if err != nil {
				return nil, err
			}
			list = append(list, intValue)
			bencodedString = bencodedString[endIndex+1:]
		} else if unicode.IsDigit(rune(bencodedString[0])) {
			endIndex := strings.Index(bencodedString, ":")
			length, _ := strconv.Atoi(bencodedString[:endIndex])
			strValue, err := decodeString(bencodedString[:endIndex+1+length])
			if err != nil {
				return nil, err
			}
			list = append(list, strValue)
			bencodedString = bencodedString[endIndex+1+length:]
		} else {
			fmt.Println("Invalid bencoded string", bencodedString)
			return nil, fmt.Errorf("invalid bencoded string")
		}
	}
	return list, nil
}
