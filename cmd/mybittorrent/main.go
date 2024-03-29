package main

import (
	// Uncomment this line to pass the first stage
	"bufio"
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
	} else if bencodedString[0] == 'd' {
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
	n, err := strconv.Atoi(bencodedString[1:charindex])
	if err != nil {
		return 0, fmt.Errorf("failed to decode integer %s: %v", bencodedString, err), 0
	}

	return n, nil, charindex + 1
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
			return nil, err, i
		}

		i += n
		array = append(array, value)
	}
	return array, nil, i + 1
}

func decodeDictionary(bencodedString string) (interface{}, error, int) {
	result := make(map[string]interface{})
	i := 1
	for {
		if i >= len(bencodedString) {
			return nil, fmt.Errorf("Not found"), i
		}

		if bencodedString[i] == 'e' {
			break
		}

		key, err, n := decodeBencode(bencodedString[i:])
		if err != nil {
			return nil, err, 0
		}
		i += n

		value, err2, n2 := decodeBencode(bencodedString[i:])
		if err2 != nil {
			return nil, err2, 0
		}

		i += n2
		result[key.(string)] = value
	}

	return result, nil, i + 1
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

	} else if command == "info" {
		file, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			decoded, err, _ := decodeBencode(scanner.Text())
			if err != nil {
				fmt.Println(err)
				return
			}

			jsonOutput, _ := json.Marshal(decoded)
			var data map[string]interface{}
			err = json.Unmarshal([]byte(jsonOutput), &data)
			if err != nil {
				fmt.Printf("could not unmarshal json: %s\n", err)
				return
			}

			jsonOutput2, _ := json.Marshal(data["info"])
			var data2 map[string]interface{}
			err = json.Unmarshal([]byte(jsonOutput2), &data2)
			if err != nil {
				fmt.Printf("could not unmarshal json: %s\n", err)
				return
			}

			fmt.Printf(
				"Tracker URL: %v\nLength: %v\n",
				data["announce"],
				data2["length"])
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
