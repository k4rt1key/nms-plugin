package polling

import (
	"context"
	"nms-plugin/commands"
	"sync"
	"time"

	"nms-plugin/logger"
	"nms-plugin/models"
	"nms-plugin/winrm"
)

const (
	// PollingTimeout is the timeout for polling operations
	PollingTimeout = 60 * time.Second
)

// Result represents the result of a polling operation
type Result struct {
	MetricGroup models.MetricGroup
	Success     bool
	Data        string
	Message     string
	Time        time.Time
}

// Execute performs polling for the specified metric groups
func Execute(request models.PollingRequest) models.PollingResponse {
	logger.Info("Starting polling process for %d metric groups", len(request.MetricGroups))

	// Create the response object
	response := models.PollingResponse{
		Type:         models.PollingType,
		MetricGroups: make([]models.PollingResult, 0, len(request.MetricGroups)),
	}

	// Create a channel to receive results from workers
	resultsChan := make(chan Result)

	// Create a wait group to track workers
	var wg sync.WaitGroup

	// Track the number of workers we'll spawn
	totalWorkers := len(request.MetricGroups)
	wg.Add(totalWorkers)

	// Create a context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), PollingTimeout)
	defer cancel()

	// Start a worker for each metric group
	for _, metricGroup := range request.MetricGroups {
		go func(mg models.MetricGroup) {
			defer wg.Done()

			startTime := time.Now()

			// Create a WinRM client
			client := winrm.NewClient(
				mg.IP,
				mg.Port,
				mg.Credentials.Username,
				mg.Credentials.Password,
			)

			// Get data for the metric group
			command := commands.GetCommand(mg.Name)
			data, err := client.ExecuteCommand(ctx, command)

			// Create result
			result := Result{
				MetricGroup: mg,
				Time:        startTime,
			}

			if err != nil {
				result.Success = false
				result.Message = err.Error()
				logger.Error("Failed to get metrics for group %s (ID: %d) on %s: %v",
					mg.Name, mg.ProvisionProfileID, mg.IP, err)
			} else {
				result.Success = true
				result.Data = data
				result.Message = "success"
				logger.Debug("Successfully retrieved metrics for group %s (ID: %d) on %s",
					mg.Name, mg.ProvisionProfileID, mg.IP)
			}

			resultsChan <- result

		}(metricGroup)
	}

	// Start a goroutine to close the results channel once all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Process results as they come in
	for result := range resultsChan {
		pollingResult := models.PollingResult{
			Success:            result.Success,
			ProvisionProfileID: result.MetricGroup.ProvisionProfileID,
			Name:               result.MetricGroup.Name,
			Data:               result.Data,
			Message:            result.Message,
			Time:               result.Time.Format(time.RFC3339),
		}

		response.MetricGroups = append(response.MetricGroups, pollingResult)
	}

	logger.Info("Polling completed for %d metric groups", len(request.MetricGroups))

	return response
}
