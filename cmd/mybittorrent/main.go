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
func decodeBencode(bencodedString string) (interface{}, error, string) {
	if unicode.IsDigit(rune(bencodedString[0])) {
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
			return "", err, "error"
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil, "string"
	} else if rune(bencodedString[0]) == rune('i') {
		result, err, types := decodeInt(bencodedString)
		if err != nil {
			return "", err, "error"
		}
		return result, err, types
	} else if rune(bencodedString[0]) == rune('l') {
		result, err, types := decodeLists(bencodedString[1 : len(bencodedString)-1])
		if err != nil {
			return "", err, "error"
		}
		return result, err, types
	} else if rune(bencodedString[0]) == rune('d') {
		return decodeDictionary(bencodedString[1 : len(bencodedString)-1])
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment"), "error"
	}
}

func decodeInt(bencodedString string) (interface{}, error, string) {
	var charindex int
	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			charindex = i
			break
		}
	}
	return bencodedString[1:charindex], nil, "int"
}

func decodeLists(bencodedString string) (interface{}, error, string) {
	// var index int = 1
	var length int
	var lengthOfInteger int

	if len(bencodedString) == 0 {
		return "[]", nil, "int"
	}

	decoded, err, types := decodeBencode(bencodedString[0:])

	if err != nil {
		return "", err, "error"
	}

	length = len(decoded.(string))

	if types == "string" {
		lengthOfInteger = len(strconv.Itoa(length))

		if (lengthOfInteger + length + 2) < len(bencodedString) {
			decoded2, err, _ := decodeBencode(bencodedString[lengthOfInteger+length+1:])
			if err != nil {
				return "", err, "error"
			}

			return fmt.Sprint("[" + "\"" + decoded.(string) + "\"" + "," + decoded2.(string) + "]"), err, "array"
		}
	}

	if types == "int" {
		decoded2, err, types := decodeBencode(bencodedString[length+2:])
		if err != nil {
			return "", err, "error"
		}

		if types == "string" {
			decoded2 = "\"" + decoded2.(string) + "\""
		}

		return fmt.Sprint("[" + decoded.(string) + "," + decoded2.(string) + "]"), err, "array"
	}

	return fmt.Sprint("[" + decoded.(string) + "]"), err, "array"
}

func decodeDictionary(bencodedString string) (interface{}, error, string) {
	result := make(map[string]interface{})

	for len(bencodedString) > 0 {
		key, err, _ := decodeBencode(bencodedString)
		if err != nil {
			return nil, err, "error"
		}

		bencodedString = bencodedString[len(key.(string))+2:]

		value, err, _ := decodeBencode(bencodedString)
		if err != nil {
			return nil, err, "error"
		}

		bencodedString = bencodedString[len(value.(string)):]

		result[key.(string)] = value
	}

	return result, nil, "dictionary"
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage

		bencodedValue := os.Args[2]

		decoded, err, types := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}
		if types == "string" {
			jsonOutput, _ := json.Marshal(decoded)
			fmt.Println(string(jsonOutput))
		}
		if types == "int" || types == "array" {
			fmt.Println(decoded)
		}

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
