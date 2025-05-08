package main

import (
	"encoding/json"
	"fmt"
	"nms-plugin/src/discovery"
	"nms-plugin/src/logger"
	"nms-plugin/src/polling"
	"nms-plugin/src/schema"
	"os"
	"strings"
)

func main() {

	// Initialize the logger
	logger.Initialize()
	defer logger.Close()

	logger.Info("WinRM Plugin started")

	// Validate command-line arguments
	if len(os.Args) < 2 {
		logger.Fatal("No input provided. Usage: %s '{\"type\": \"discovery\", ...}'", os.Args[0])
	}

	// Get the JSON input from command-line arguments
	input := strings.Join(os.Args[1:], " ")
	logger.Debug("Received input: %s", input)

	// Parse the input to determine the request type
	var requestType struct {
		Type schema.RequestType `json:"type"`
	}

	if err := json.Unmarshal([]byte(input), &requestType); err != nil {
		logger.Fatal("Failed to parse input JSON: %v", err)
	}

	var response interface{}

	// Process the request based on its type
	switch requestType.Type {
	case schema.DiscoveryType:
		response = handleDiscovery(input)
	case schema.PollingType:
		response = handlePolling(input)
	default:
		logger.Fatal("Unknown request type: %s", requestType.Type)
	}

	// Convert the response to JSON
	outputJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		logger.Fatal("Failed to marshal response to JSON: %v", err)
	}

	// Print the JSON response
	fmt.Println(string(outputJSON))

	logger.Info("WinRM Plugin completed successfully")
}

// handleDiscovery processes a discovery request
func handleDiscovery(input string) schema.DiscoveryResponse {
	var request schema.DiscoveryRequest
	if err := json.Unmarshal([]byte(input), &request); err != nil {
		logger.Fatal("Failed to parse discovery request: %v", err)
	}

	logger.Info("Processing discovery request with ID %d for %d IPs", request.ID, len(request.IPs))
	return discovery.Execute(request)
}

// handlePolling processes a polling request
func handlePolling(input string) schema.PollingResponse {
	var request schema.PollingRequest
	if err := json.Unmarshal([]byte(input), &request); err != nil {
		logger.Fatal("Failed to parse polling request: %v", err)
	}

	logger.Info("Processing polling request for %d metric groups", len(request.MetricGroups))
	return polling.Execute(request)
}
