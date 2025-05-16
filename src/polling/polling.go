package polling

import (
	"context"
	"nms-plugin/src/logger"
	"nms-plugin/src/types"
	"nms-plugin/src/winrm"
	"sync"
	"time"
)

const (
	// Timeout is the timeout for polling operations
	Timeout = 60 * time.Second
)

// Result represents the result of a polling operation
type Result struct {
	MetricGroup types.MetricGroup
	Success     bool
	Data        string
	Message     string
	Time        time.Time
}

// Execute performs polling for the specified metric groups
func Execute(request types.PollingRequest) types.PollingResponse {
	logger.Info("Starting polling process for %d metric groups", len(request.MetricGroups))

	// Create the response object
	response := types.PollingResponse{
		Type:         types.PollingType,
		MetricGroups: make([]types.PollingResult, 0, len(request.MetricGroups)),
	}

	// Create a channel to receive results from workers
	resultsChan := make(chan Result)

	// Create a wait group to track workers
	var wg sync.WaitGroup

	// Track the number of workers we'll spawn
	totalWorkers := len(request.MetricGroups)

	wg.Add(totalWorkers)

	// Create a context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	// Start a worker for each metric group
	for _, metricGroup := range request.MetricGroups {
		go func(mg types.MetricGroup) {
			defer wg.Done()

			startTime := time.Now()

			// Create a WinRM client
			client := winrm.NewClient(
				mg.IP,
				mg.Port,
				mg.Credential.Username,
				mg.Credential.Password,
			)

			// Get data for the metric group
			command := getCommand(mg.Name)
			data, err := client.ExecuteCommand(ctx, command)

			// Create result
			result := Result{
				MetricGroup: mg,
				Time:        startTime,
			}

			if err != nil {
				result.Success = false
				result.Message = err.Error()
				logger.Error("Failed to get metrics for group %s (ID: %d) on %s: %v", mg.Name, mg.MonitorID, mg.IP, err)
			} else {
				result.Success = true
				result.Data = data
				result.Message = "success"
				logger.Debug("Successfully retrieved metrics for group %s (ID: %d) on %s",
					mg.Name, mg.MonitorID, mg.IP)
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
		pollingResult := types.PollingResult{
			Success:            result.Success,
			ProvisionProfileID: result.MetricGroup.MonitorID,
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

func getCommand(name string) string {

	switch name {

	case "CPUINFO":

		return `powershell -Command "Get-CimInstance Win32_Processor | Select-Object Name, NumberOfCores, NumberOfLogicalProcessors, MaxClockSpeed, Manufacturer, LoadPercentage | ConvertTo-Json -Depth 3"`

	case "UPTIME":

		return `powershell -Command "([pscustomobject]@{ uptime = ((get-date) - (gcim Win32_OperatingSystem).LastBootUpTime).TotalSeconds }) | ConvertTo-Json"`

	case "CPUUSAGE":

		return `powershell -Command "$c = Get-Counter '\Processor(_Total)\% User Time','\Processor(_Total)\% Privileged Time','\Processor(_Total)\% Idle Time','\Processor(_Total)\% Processor Time'; $d = $c.CounterSamples; [pscustomobject]@{ user = ($d | ? {$_.Path -like '*User*'}).CookedValue; system = ($d | ? {$_.Path -like '*Privileged*'}).CookedValue; idle = ($d | ? {$_.Path -like '*Idle*'}).CookedValue; total = ($d | ? {$_.Path -like '*Processor Time'}).CookedValue } | ConvertTo-Json"`

	case "MEMORY":

		return `powershell -Command "Get-CimInstance Win32_OperatingSystem | Select-Object TotalVisibleMemorySize, FreePhysicalMemory, TotalVirtualMemorySize, FreeVirtualMemory| ConvertTo-Json -Depth 3"`

	case "DISK":

		return `powershell -Command "Get-CimInstance Win32_LogicalDisk | Select-Object DeviceID, VolumeName, FileSystem, VolumeSerialNumber, DriveType, @{Name='SizeMB';Expression={[math]::Round($_.Size / 1MB, 2)}}, @{Name='FreeSpaceMB';Expression={[math]::Round($_.FreeSpace / 1MB, 2)}}, @{Name='UsedSpaceMB';Expression={[math]::Round(($_.Size - $_.FreeSpace) / 1MB, 2)}}, @{Name='FreeSpacePercent';Expression={[math]::Round(($_.FreeSpace / $_.Size) * 100, 2)}}, @{Name='DriveLetter';Expression={$_.DeviceID}}, @{Name='FileSystemType';Expression={$_.FileSystem}}, @{Name='DriveMediaType';Expression={$_.MediaType}} | ConvertTo-Json -Depth 3"`

	case "NETWORK":

		return `powershell -Command "Get-CimInstance Win32_PerfRawData_Tcpip_TCPv4 | Select-Object ConnectionsEstablished, ConnectionsActive, ConnectionFailures, ConnectionsPassive | ConvertTo-Json -Depth 3"`

	case "PROCESS":

		return `powershell -Command "Write-Output (@{ TotalProcesses = (Get-Process).Count; IdleProcesses = (Get-Process -Name Idle -ErrorAction SilentlyContinue).Count; SystemProcesses = (Get-Process -Name System -ErrorAction SilentlyContinue).Count; RunningProcesses = (Get-Process | Where-Object { $_.CPU -gt 0 }).Count; SleepProcesses = (Get-Process | Where-Object { $_.CPU -eq $null -or $_.CPU -eq 0 } | Where-Object { $_.Name -ne 'Idle' -and $_.Name -ne 'System' }).Count } | ConvertTo-Json)"`

	case "SYSTEMINFO":

		return `powershell -Command "Get-CimInstance Win32_ComputerSystem | Select-Object Name, Manufacturer, Model, , TotalPhysicalMemory, Domain, UserName, PowerState | ConvertTo-Json -Depth 3"`

	default:

		return `powershell -Command "([pscustomobject]@{ uptime = ((get-date) - (gcim Win32_OperatingSystem).LastBootUpTime).TotalSeconds }) | ConvertTo-Json"`
	}
}
