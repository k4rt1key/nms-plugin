package commands

func GetCommand(name string) string {
	switch name {
	case "CPU":
		// Get CPU details
		return `powershell -Command "Get-CimInstance Win32_Processor | Select-Object Name,NumberOfCores,NumberOfLogicalProcessors"`

	case "MEMORY":
		// Get total and available memory
		return `powershell -Command "Get-CimInstance Win32_OperatingSystem | Select-Object TotalVisibleMemorySize,FreePhysicalMemory"`

	case "DISK":
		// Get disk usage info
		return `powershell -Command "Get-CimInstance Win32_LogicalDisk | Where-Object { $_.DriveType -eq 3 } | Select-Object DeviceID,Size,FreeSpace"`

	case "NETWORK":
		// Get IP addresses and adapter names
		return `powershell -Command "Get-NetIPAddress | Select-Object InterfaceAlias,IPAddress"`

	case "FILE":
		// List files in a directory (change C:\Test as needed)
		return `powershell -Command "Get-ChildItem -Path 'C:\' | Select-Object Name,Length"`

	case "PROCCESS":
		// List running processes
		return `powershell -Command "Get-Process | Select-Object Name,Id,CPU"`

	case "SYSTEMINFO":
		// System summary info (legacy but widely available)
		return `systeminfo`

	default:
		return `systeminfo`
	}
}
