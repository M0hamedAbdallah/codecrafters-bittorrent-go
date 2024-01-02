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
	} else if rune(bencodedString[0]) == rune('l') && rune(bencodedString[len(bencodedString)-1]) == rune('e') {
		decoded, err, _ := decodeBencode(bencodedString[1:])
		if err != nil {
			return "", err, "error"
		}

		length := len(decoded.(string))

		lengthOfInteger := len(strconv.Itoa(length))

		decode2, err, _ := decodeBencode(bencodedString[lengthOfInteger+length+2 : len(bencodedString)-1])
		if err != nil {
			return "", err, "error"
		}

		return fmt.Sprint("[" + "\"" + decoded.(string) + "\"" + "," + decode2.(string) + "]"), err, "array"
	} else if rune(bencodedString[0]) == rune('d') && rune(bencodedString[len(bencodedString)-1]) == rune('e') {
		return "", nil, ""
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
		if types == "int" || types == "array" {
			fmt.Println(decoded)
		}

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
