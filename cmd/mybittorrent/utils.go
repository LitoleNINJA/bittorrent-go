package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/jackpal/bencode-go"
)

func decodeBencode(bencodedString string) (interface{}, int, error) {
	switch bencodedString[0] {
	case 'i':
		return decodeInt(bencodedString, 0)
	case 'l':
		return decodeList(bencodedString, 0)
	case 'd':
		return decodeDict(bencodedString, 0)
	default:
		if unicode.IsDigit(rune(bencodedString[0])) {
			return decodeString(bencodedString, 0)
		} else {
			return "", -1, fmt.Errorf("invalid bencoded string")
		}
	}
}

func decodeString(bencodedString string, pos int) (string, int, error) {
	firstColonIndex := strings.Index(bencodedString[pos:], ":") + pos
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

func decodeDict(bencodedString string, pos int) (map[string]interface{}, int, error) {
	dict := make(map[string]interface{})
	key := ""
	end := 0
	for i := pos + 1; i < len(bencodedString); i++ {
		ch := bencodedString[i]
		if ch == 'e' {
			end = i
			break
		} else if ch == 'i' {
			decodedInt, endPos, err := decodeInt(bencodedString, i)
			if err != nil {
				return nil, -1, err
			}
			i = endPos
			if key == "" {
				key = strconv.Itoa(decodedInt)
			} else {
				dict[key] = decodedInt
				key = ""
			}
		} else if unicode.IsDigit(rune(ch)) {
			decodedString, endPos, err := decodeString(bencodedString, i)
			if err != nil {
				return nil, -1, err
			}
			i = endPos
			if key == "" {
				key = decodedString
			} else {
				dict[key] = decodedString
				key = ""
			}
		} else if ch == 'l' {
			decodedList, endPos, err := decodeList(bencodedString, i)
			if err != nil {
				return nil, -1, err
			}
			i = endPos
			if key == "" {
				key = "list"
			} else {
				dict[key] = decodedList
				key = ""
			}
		} else if ch == 'd' {
			decodedDict, endPos, err := decodeDict(bencodedString, i)
			if err != nil {
				return nil, -1, err
			}
			i = endPos
			if key == "" {
				key = "dict"
			} else {
				dict[key] = decodedDict
				key = ""
			}
		} else {
			return nil, -1, fmt.Errorf("invalid bencoded string")
		}
	}
	return dict, end, nil
}

func encodeBencode(data interface{}) (string, error) {
	switch v := data.(type) {
	case string:
		return EncodeString(v)
	case int:
		return EncodeNumber(v)
	case []interface{}:
		return EncodeList(v)
	case map[string]interface{}:
		return EncodeDictionary(v)
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

func EncodeString(data string) (string, error) {
	return strconv.Itoa(len(data)) + ":" + data, nil
}

func EncodeNumber(data int) (string, error) {
	return "i" + strconv.Itoa(data) + "e", nil
}

func EncodeList(data []interface{}) (string, error) {
	encodedList := "l"
	for _, v := range data {
		encoded, err := encodeBencode(v)
		if err != nil {
			return "", err
		}
		encodedList += encoded
	}
	encodedList += "e"
	return encodedList, nil
}

func EncodeDictionary(data map[string]interface{}) (string, error) {
	encodedDict := "d"
	for k, v := range data {
		encodedKey, err := EncodeString(k)
		if err != nil {
			return "", err
		}
		encodedDict += encodedKey
		encodedValue, err := encodeBencode(v)
		if err != nil {
			return "", err
		}
		encodedDict += encodedValue
	}
	encodedDict += "e"
	return encodedDict, nil
}

func readTorrentFile(torrentFile string) (Torrent, error) {
	file, err := os.Open(torrentFile)
	if err != nil {
		return Torrent{}, err
	}
	defer file.Close()

	var torrent Torrent
	err = bencode.Unmarshal(file, &torrent)
	if err != nil {
		return Torrent{}, err
	}
	return torrent, nil
}
