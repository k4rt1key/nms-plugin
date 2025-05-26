package polling

import (
	"context"
	"sync"
	"time"

	"nms-plugin/src/winrm"
)

const timeout = 60 * time.Second

func Execute(request map[string]interface{}) map[string]interface{} {

	metricGroups := request["metric_groups"].([]interface{})

	protocol := getProtocol(request)

	response := map[string]interface{}{
		"type":          "polling",
		"metric_groups": []map[string]interface{}{},
	}

	switch protocol {

	case "winrm":

		results := pollWinRM(metricGroups)

		response["metric_groups"] = results

	default:

		results := pollWinRM(metricGroups)

		response["metric_groups"] = results
	}

	return response
}

func getProtocol(request map[string]interface{}) string {

	if protocol, ok := request["protocol"]; ok {

		return protocol.(string)

	}

	return "winrm"
}

func pollWinRM(metricGroups []interface{}) []map[string]interface{} {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	resultChan := make(chan map[string]interface{})

	var wg sync.WaitGroup

	wg.Add(len(metricGroups))

	for _, mgInterface := range metricGroups {

		go func(mg map[string]interface{}) {

			defer wg.Done()

			cred := mg["credential"].(map[string]interface{})

			client := winrm.NewClient(

				mg["ip"].(string),

				int(mg["port"].(float64)),

				cred["username"].(string),

				cred["password"].(string),
			)

			command := getCommand(mg["name"].(string))

			data, err := winrm.ExecuteCommand(ctx, client, command)

			result := map[string]interface{}{

				"monitor_id": mg["monitor_id"],

				"name": mg["name"],

				"time": time.Now().Format(time.RFC3339),
			}

			if err != nil {

				result["success"] = false

				result["message"] = err.Error()

				result["data"] = ""

			} else {

				result["success"] = true

				result["message"] = "success"

				result["data"] = data

			}

			resultChan <- result

		}(mgInterface.(map[string]interface{}))
	}

	go func() {

		wg.Wait()

		close(resultChan)

	}()

	var results []map[string]interface{}

	for result := range resultChan {

		results = append(results, result)

	}

	return results
}

func getCommand(name string) string {

	commands := map[string]string{

		"CPUINFO": `powershell -Command "Get-CimInstance Win32_Processor | Select-Object Name, NumberOfCores, NumberOfLogicalProcessors, MaxClockSpeed, Manufacturer, LoadPercentage | ConvertTo-Json -Depth 3"`,

		"UPTIME": `powershell -Command "([pscustomobject]@{ uptime = ((get-date) - (gcim Win32_OperatingSystem).LastBootUpTime).TotalSeconds }) | ConvertTo-Json"`,

		"CPUUSAGE": `powershell -Command "$c = Get-Counter '\Processor(_Total)\% User Time','\Processor(_Total)\% Privileged Time','\Processor(_Total)\% Idle Time','\Processor(_Total)\% Processor Time'; $d = $c.CounterSamples; [pscustomobject]@{ user = ($d | ? {$_.Path -like '*User*'}).CookedValue; system = ($d | ? {$_.Path -like '*Privileged*'}).CookedValue; idle = ($d | ? {$_.Path -like '*Idle*'}).CookedValue; total = ($d | ? {$_.Path -like '*Processor Time'}).CookedValue } | ConvertTo-Json"`,

		"MEMORY": `powershell -Command "Get-CimInstance Win32_OperatingSystem | Select-Object TotalVisibleMemorySize, FreePhysicalMemory, TotalVirtualMemorySize, FreeVirtualMemory| ConvertTo-Json -Depth 3"`,

		"DISK": `powershell -Command "Get-CimInstance Win32_LogicalDisk | Select-Object DeviceID, VolumeName, FileSystem, VolumeSerialNumber, DriveType, @{Name='SizeMB';Expression={[math]::Round($_.Size / 1MB, 2)}}, @{Name='FreeSpaceMB';Expression={[math]::Round($_.FreeSpace / 1MB, 2)}}, @{Name='UsedSpaceMB';Expression={[math]::Round(($_.Size - $_.FreeSpace) / 1MB, 2)}}, @{Name='FreeSpacePercent';Expression={[math]::Round(($_.FreeSpace / $_.Size) * 100, 2)}}, @{Name='DriveLetter';Expression={$_.DeviceID}}, @{Name='FileSystemType';Expression={$_.FileSystem}}, @{Name='DriveMediaType';Expression={$_.MediaType}} | ConvertTo-Json -Depth 3"`,

		"NETWORK": `powershell -Command "Get-CimInstance Win32_PerfRawData_Tcpip_TCPv4 | Select-Object ConnectionsEstablished, ConnectionsActive, ConnectionFailures, ConnectionsPassive | ConvertTo-Json -Depth 3"`,

		"PROCESS": `powershell -Command "Write-Output (@{ TotalProcesses = (Get-Process).Count; IdleProcesses = (Get-Process -Name Idle -ErrorAction SilentlyContinue).Count; SystemProcesses = (Get-Process -Name System -ErrorAction SilentlyContinue).Count; RunningProcesses = (Get-Process | Where-Object { $_.CPU -gt 0 }).Count; SleepProcesses = (Get-Process | Where-Object { $_.CPU -eq $null -or $_.CPU -eq 0 } | Where-Object { $_.Name -ne 'Idle' -and $_.Name -ne 'System' }).Count } | ConvertTo-Json)"`,

		"SYSTEMINFO": `powershell -Command "Get-CimInstance Win32_ComputerSystem | Select-Object Name, Manufacturer, Model, TotalPhysicalMemory, Domain, UserName, PowerState | ConvertTo-Json -Depth 3"`,
	}

	if cmd, exists := commands[name]; exists {

		return cmd

	}

	return commands["UPTIME"]
}
