package main

import (
	"encoding/json"
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
	
	switch request["type"] {

	case "discovery":

		discovery.Execute(request)

	case "polling":

		polling.Execute(request)

	default:

		log.Fatal("Unknown request type:", request["type"])

	}

}
