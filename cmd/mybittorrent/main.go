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
	} else if rune(bencodedString[0]) == rune('i') && rune(bencodedString[len(bencodedString)-1]) == rune('e') {
		return bencodedString[1 : len(bencodedString)-1], nil, "int"
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment"), "error"
	}
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
		if types == "int" {
			fmt.Println(decoded)
		}

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
