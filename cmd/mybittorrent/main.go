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

		decoded, _, err := decodeBencode(bencodedValue)
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

func decodeBencode(bencodedString string) (interface{}, int, error) {
	switch bencodedString[0] {
	case 'i':
		return decodeInt(bencodedString, 0)
	case 'l':
		return decodeList(bencodedString, 0)
	default:
		if unicode.IsDigit(rune(bencodedString[0])) {
			return decodeString(bencodedString, 0)
		} else {
			return "", -1, fmt.Errorf("invalid bencoded string")
		}
	}
}

func decodeString(bencodedString string, pos int) (string, int, error) {
	firstColonIndex := strings.Index(bencodedString, ":")
	lengthStr := bencodedString[pos:firstColonIndex]
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 0, err
	}
	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], firstColonIndex + length, nil
}

func decodeInt(bencodedString string, pos int) (int, int, error) {
	for i := pos; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			decodedInt, err := strconv.Atoi(bencodedString[pos+1 : i])
			if err != nil {
				return 0, 0, err
			}
			return decodedInt, i, nil
		}
	}
	return 0, 0, fmt.Errorf("invalid bencoded string")
}

func decodeList(bencodedString string, pos int) ([]interface{}, int, error) {
	list := []interface{}{}
	end := 0
	for i := pos + 1; i < len(bencodedString); i++ {
		ch := bencodedString[i]
		// fmt.Printf("i: %d, ch: %c\n", i, ch)
		if ch == 'e' {
			end = i
			break
		} else if ch == 'i' {
			decodedInt, endPos, err := decodeInt(bencodedString, i)
			if err != nil {
				return nil, -1, err
			}
			list = append(list, decodedInt)
			i = endPos
		} else if ch == 'l' {
			decodedList, endPos, err := decodeList(bencodedString, i)
			if err != nil {
				return nil, -1, err
			}
			list = append(list, decodedList)
			i = endPos
		} else if unicode.IsDigit(rune(ch)) {
			decodedString, endPos, err := decodeString(bencodedString, i)
			if err != nil {
				return nil, -1, err
			}
			list = append(list, decodedString)
			i = endPos
		}
	}
	return list, end, nil
}
