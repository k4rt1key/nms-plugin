package commands

func GetCommand(name string) string {
	switch name {
	case "CPU":
		// Grouped CPU-related metrics
		return `powershell -Command "Get-CimInstance Win32_Processor | Select-Object Name, NumberOfCores, NumberOfLogicalProcessors, MaxClockSpeed, Manufacturer, LoadPercentage, @{Name='CPUUserTime';Expression={([System.Diagnostics.PerformanceCounter]::new('Processor', '% User Time', '_Total').NextValue())}}, @{Name='CPUSystemTime';Expression={([System.Diagnostics.PerformanceCounter]::new('Processor', '% Privileged Time', '_Total').NextValue())}}, @{Name='CPUIdleTime';Expression={([System.Diagnostics.PerformanceCounter]::new('Processor', '% Idle Time', '_Total').NextValue())}} | ConvertTo-Json -Depth 3"`

	case "MEMORY":
		// Grouped Memory-related metrics
		return `powershell -Command "Get-CimInstance Win32_OperatingSystem | Select-Object TotalVisibleMemorySize, FreePhysicalMemory, TotalVirtualMemorySize, FreeVirtualMemory, CommittedBytes, TotalSwapSpaceSize, FreeSwapSpaceSize | ConvertTo-Json -Depth 3"`

	case "DISK":
		// Grouped Disk-related metrics
		return `powershell -Command "Get-CimInstance Win32_LogicalDisk | Select-Object DeviceID, VolumeName, FileSystem, VolumeSerialNumber, DriveType, @{Name='SizeMB';Expression={[math]::Round($_.Size / 1MB, 2)}}, @{Name='FreeSpaceMB';Expression={[math]::Round($_.FreeSpace / 1MB, 2)}}, @{Name='UsedSpaceMB';Expression={[math]::Round(($_.Size - $_.FreeSpace) / 1MB, 2)}}, @{Name='FreeSpacePercent';Expression={[math]::Round(($_.FreeSpace / $_.Size) * 100, 2)}}, @{Name='DriveLetter';Expression={$_.DeviceID}}, @{Name='FileSystemType';Expression={$_.FileSystem}}, @{Name='DriveMediaType';Expression={$_.MediaType}} | ConvertTo-Json -Depth 3"`

	case "NETWORK":
		// Grouped Network-related metrics
		return `powershell -Command "$interfaces = Get-NetIPAddress | Select-Object InterfaceAlias,IPAddress,AddressFamily,PrefixLength,InterfaceIndex,Type; $adapters = Get-NetAdapter | Select-Object Name,Status; $output = foreach ($i in $interfaces) { $status = ($adapters | Where-Object { $_.Name -eq $i.InterfaceAlias }).Status; [PSCustomObject]@{ InterfaceAlias=$i.InterfaceAlias; IPAddress=$i.IPAddress; AddressFamily=$i.AddressFamily; PrefixLength=$i.PrefixLength; InterfaceIndex=$i.InterfaceIndex; Type=$i.Type; Status=$status } }; $output | ConvertTo-Json -Depth 3"`

	case "PROCESS":
		// Grouped Process-related metrics
		return `powershell -Command "$procs = Get-Process; $perfData = Get-CimInstance -ClassName Win32_PerfFormattedData_PerfProc_Process; $procs | ForEach-Object { $id = $_.Id; $ctx = $perfData | Where-Object { $_.IDProcess -eq $id }; [PSCustomObject]@{ Name = $_.Name; Id = $id; CPU = $_.CPU; Threads = $_.Threads.Count; Handles = $_.Handles; WorkingSetMB = [math]::Round($_.WorkingSet64 / 1MB, 2); VirtualMemoryMB = [math]::Round($_.VirtualMemorySize64 / 1MB, 2); StartTime = $_.StartTime; Path = $_.Path; Responding = $_.Responding; PriorityClass = $_.PriorityClass } } | ConvertTo-Json -Depth 3"`

	case "SYSTEMINFO":
		// Grouped System-related metrics
		return `powershell -Command "Get-CimInstance Win32_ComputerSystem | Select-Object Name, Manufacturer, Model, , TotalPhysicalMemory, Domain, UserName, PowerState | ConvertTo-Json -Depth 3"`

	default:
		// Default command for operating system memory details
		return `
		powershell -Command "Get-CimInstance Win32_OperatingSystem | Select-Object Caption, Version, BuildNumber, @{Name='SystemUpTime';Expression={(New-TimeSpan -Start $_.LastBootUpTime -End (Get-Date)).ToString()}} | ConvertTo-Json -Depth 3"
`
	}
}
