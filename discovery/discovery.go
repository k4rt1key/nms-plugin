package discovery

import (
	"context"
	"sync"
	"time"

	"nms-plugin/logger"
	"nms-plugin/models"
	"nms-plugin/parser"
	"nms-plugin/winrm"
)

const (
	// DiscoveryTimeout is the timeout for discovery operations
	DiscoveryTimeout = 15 * time.Second
)

// Result represents the result of a discovery operation
type Result struct {
	IP         string
	Credential models.Credential
	Success    bool
	Message    string
}

// Execute performs discovery on the specified IPs using the provided credentials
func Execute(request models.DiscoveryRequest) models.DiscoveryResponse {
	logger.Info("Starting discovery process for %d IPs with %d credentials", len(request.IPs), len(request.Credentials))

	// Create the response object
	response := models.DiscoveryResponse{
		Type:   models.DiscoveryType,
		ID:     request.ID,
		Result: make([]models.DiscoveryResult, 0),
	}

	// Create a channel to receive results from workers
	resultsChan := make(chan Result)

	// Create a wait group to track workers
	var wg sync.WaitGroup

	// Track the number of workers we'll spawn
	totalWorkers := len(request.IPs) * len(request.Credentials)
	wg.Add(totalWorkers)

	// Create a map to track successful credentials for each IP
	successfulResults := make(map[string]models.DiscoveryResult)
	var resultsMutex sync.Mutex

	// Create a context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), DiscoveryTimeout)
	defer cancel()

	// Start workers to check each IP with each credential
	for _, ip := range request.IPs {
		for _, credential := range request.Credentials {
			go func(ip string, credential models.Credential) {
				defer wg.Done()

				// Create a WinRM client
				client := winrm.NewClient(ip, request.Port, credential.Username, credential.Password)

				// Test the connection
				success, output, err := client.TestConnection(ctx)

				// Send the result
				result := Result{
					IP:         ip,
					Credential: credential,
					Success:    success,
				}

				if err != nil {
					result.Message = err.Error()
				} else {
					result.Message = parser.ParseDiscoveryOutput(output)
				}

				resultsChan <- result

			}(ip, credential)
		}
	}

	// Start a goroutine to close the results channel once all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Process results as they come in
	for result := range resultsChan {
		logger.Debug("Received discovery result for %s with credential ID %d: success=%v",
			result.IP, result.Credential.ID, result.Success)

		if result.Success {
			// We found a working credential for this IP
			resultsMutex.Lock()
			successfulResults[result.IP] = models.DiscoveryResult{
				Success:    true,
				IP:         result.IP,
				Credential: result.Credential,
				Port:       request.Port,
				Message:    parser.ParseDiscoveryOutput(result.Message),
			}
			resultsMutex.Unlock()
		}
	}

	// Process the final results
	ipResults := make(map[string]bool)

	// First, add all successful results
	for ip, result := range successfulResults {
		response.Result = append(response.Result, result)
		ipResults[ip] = true
	}

	// Then, add failed results for IPs that have no successful credentials
	for _, ip := range request.IPs {
		if !ipResults[ip] {
			response.Result = append(response.Result, models.DiscoveryResult{
				Success:    false,
				IP:         ip,
				Credential: models.Credential{},
				Port:       request.Port,
				Message:    "Connection failed or invalid credentials for this IP",
			})
		}
	}

	logger.Info("Discovery completed with %d successful connections out of %d IPs",
		len(successfulResults), len(request.IPs))

	return response
}
