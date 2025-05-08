package commands

func GetCommand(name string) string {
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
