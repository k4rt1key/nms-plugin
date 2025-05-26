package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"nms-plugin/src/discovery"
	"nms-plugin/src/polling"
)

func main() {

	input := strings.Join(os.Args[1:], " ")

	var request map[string]interface{}

	if err := json.Unmarshal([]byte(input), &request); err != nil {

		log.Fatal("Failed to parse JSON:", err)

	}

	var response interface{}

	switch request["type"] {

	case "discovery":

		response = discovery.Execute(request)

	case "polling":

		response = polling.Execute(request)

	default:

		log.Fatal("Unknown request type:", request["type"])

	}

	output, _ := json.MarshalIndent(response, "", "  ")

	fmt.Println(string(output))
}
