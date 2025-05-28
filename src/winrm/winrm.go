// Package winrm winrm/winrm.go
package winrm

import (
	"context"
	"fmt"
	"time"

	"github.com/masterzen/winrm"
)

const defaultTimeout = 10 * time.Second

func NewClient(ip string, port int, username, password string) *winrm.Client {

	endpoint := winrm.NewEndpoint(ip, port, false, false, nil, nil, nil, defaultTimeout)

	client, _ := winrm.NewClient(endpoint, username, password)

	return client
}

func TestConnection(ctx context.Context, ip string, port int, cred map[string]interface{}) (bool, string) {

	client := NewClient(ip, port,
		cred["credential"].(map[string]interface{})["username"].(string),
		cred["credential"].(map[string]interface{})["password"].(string))

	stdout, stderr, exitCode, err := client.RunCmdWithContext(ctx, "hostname")

	if err != nil || exitCode != 0 {

		if stderr != "" {

			fmt.Println(stderr)

			return false, stderr

		}

		if err != nil {

			fmt.Println(err)

			return false, err.Error()

		}

		return false, fmt.Sprintf("exit code %d", exitCode)
	}

	return true, stdout
}

func Poll(ctx context.Context, mg map[string]interface{}, result map[string]interface{}) {

	cred := mg["credential"].(map[string]interface{})

	client := NewClient(

		mg["ip"].(string),

		int(mg["port"].(float64)),

		cred["credential"].(map[string]interface{})["username"].(string),

		cred["credential"].(map[string]interface{})["password"].(string),
	)

	fmt.Println(mg["ip"].(string), int(mg["port"].(float64)), cred["credential"].(map[string]interface{})["username"].(string), cred["credential"].(map[string]interface{})["password"].(string))

	command := getWinRMCommand(mg["name"].(string))

	stdout, _, exitCode, err := client.RunCmdWithContext(ctx, command)

	if err != nil {

		result["success"] = false

		result["message"] = err.Error()

		result["data"] = ""
	} else if exitCode != 0 {

		result["success"] = false

		result["message"] = fmt.Errorf("exit code %d", exitCode)

		result["data"] = ""

	} else {

		result["success"] = true

		result["message"] = "success"

		result["data"] = stdout
	}
}

func getWinRMCommand(name string) string {

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
