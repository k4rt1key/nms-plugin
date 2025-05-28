package main

import (
	"encoding/json"
	"log"
	"nms-plugin/src/discovery"
	"nms-plugin/src/polling"
	"os"
	"strings"
)

func main() {

	input := strings.Join(os.Args[1:], " ")

	var request map[string]interface{}

	if err := json.Unmarshal([]byte(input), &request); err != nil {

		log.Fatal("Failed to parse JSON:", err)

	}

	switch request["type"] {

	case "discovery":

		discovery.Discover(request)

	case "polling":

		polling.Poll(request)

	default:

		log.Fatal("Unknown request type:", request["type"])

	}

}
