package main

import (
	// Uncomment this line to pass the first stage
	"encoding/json"
	"fmt"
	"os"
	// "reflect" // to reflect variables to anthor type
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error, int) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		return decodeString(bencodedString)
	} else if rune(bencodedString[0]) == rune('i') {
		return decodeInt(bencodedString)
	} else if rune(bencodedString[0]) == rune('l') {
		return decodeLists(bencodedString)
	} else if rune(bencodedString[0]) == rune('d') {
		return decodeDictionary(bencodedString)
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment"), 0
	}
}

func decodeString(bencodedString string) (interface{}, error, int) {
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)

	if err != nil {
		return "", err, 0
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil, firstColonIndex + 1 + length
}

func decodeInt(bencodedString string) (int, error, int) {
	var charindex int
	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			charindex = i
			break
		}
	}
	n , err :=strconv.Atoi(bencodedString[1:charindex])
	if err != nil {
		return 0,  fmt.Errorf("failed to decode integer %s: %v", bencodedString , err),0
	}

	return  n , nil, charindex + 1
}

func decodeLists(bencodedString string) ([]interface{}, error, int) {
	array := make([]interface{}, 0)
	i := 1
	for {
		if i >= len(bencodedString) {
			return nil, fmt.Errorf("Not found"), i
		}

		if bencodedString[i] == 'e' {
			break
		}

		value, err, n := decodeBencode(bencodedString[i:])
		if err != nil {
			return nil, err , i
		}

		i += n
		array = append(array, value)
	}
	return array, nil, i+1
}

func decodeDictionary(bencodedString string) (interface{}, error, int) {
	result := make(map[string]interface{})

	for len(bencodedString) > 0 {
		key, err, _ := decodeBencode(bencodedString)
		if err != nil {
			return nil, err, 0
		}

		bencodedString = bencodedString[len(key.(string))+2:]

		value, err, _ := decodeBencode(bencodedString)
		if err != nil {
			return nil, err, 0
		}

		bencodedString = bencodedString[len(value.(string)):]

		result[key.(string)] = value
	}

	return result, nil, 0
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage

		bencodedValue := os.Args[2]

		decoded, err, _ := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
