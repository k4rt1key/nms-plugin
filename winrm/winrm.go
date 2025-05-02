package winrm

import (
	"context"
	"fmt"
	"time"

	"github.com/masterzen/winrm"
	"nms-plugin/logger"
	"nms-plugin/models"
)

const (
	// DefaultTimeout is the default timeout for WinRM operations
	DefaultTimeout = 30 * time.Second
)

// Client represents a WinRM client
type Client struct {
	IP       string
	Port     int
	Username string
	Password string
}

// NewClient creates a new WinRM client
func NewClient(ip string, port int, username, password string) *Client {
	return &Client{
		IP:       ip,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// ExecuteCommand executes a command on the remote host
func (c *Client) ExecuteCommand(ctx context.Context, command string) (string, error) {
	logger.Debug("Executing command '%s' on %s:%d with user '%s'", command, c.IP, c.Port, c.Username)

	// Set up WinRM client configuration
	config := &winrm.Endpoint{
		Host:     c.IP,
		Port:     c.Port,
		HTTPS:    false,
		Insecure: true,
		Timeout:  DefaultTimeout,
	}

	// Create WinRM client
	client, err := winrm.NewClient(config, c.Username, c.Password)
	if err != nil {
		logger.Error("Failed to create WinRM client for %s:%d: %v", c.IP, c.Port, err)
		return "", fmt.Errorf("failed to create WinRM client: %w", err)
	}

	// Create a channel to signal command completion
	done := make(chan struct{})
	var output string
	var cmdErr error

	// Execute command in a goroutine
	go func() {
		defer close(done)

		// Run the command and capture stdout and stderr
		var stdout, stderr string
		stdout, stderr, _, cmdErr = client.RunWithString(command, "")

		if cmdErr != nil {
			output = stderr
			logger.Error("Command execution failed on %s:%d: %v", c.IP, c.Port, cmdErr)
			return
		}

		output = stdout
		logger.Debug("Command execution succeeded on %s:%d", c.IP, c.Port)
	}()

	// Wait for command completion or context cancellation
	select {
	case <-done:
		return output, cmdErr
	case <-ctx.Done():
		logger.Error("Command execution timed out on %s:%d", c.IP, c.Port)
		return "", ctx.Err()
	}
}

// TestConnection tests if the WinRM connection is successful
func (c *Client) TestConnection(ctx context.Context) (bool, string, error) {
	output, err := c.ExecuteCommand(ctx, "systeminfo")
	if err != nil {
		return false, "", err
	}
	return true, output, nil
}

// GetCPUInfo gets CPU information from the remote host
func (c *Client) GetCPUInfo(ctx context.Context) (string, error) {
	// PowerShell command to get CPU information
	command := "powershell \"Get-Process | Sort-Object CPU -Descending | Select-Object -First 10 Name, Id, CPU, WorkingSet | Format-Table -AutoSize\""
	return c.ExecuteCommand(ctx, command)
}

// GetMemoryInfo gets memory information from the remote host
func (c *Client) GetMemoryInfo(ctx context.Context) (string, error) {
	// PowerShell command to get memory information
	command := "powershell \"$os = Get-WmiObject Win32_OperatingSystem; Write-Output ('TotalVisibleMemorySize: {0:N2} MB' -f ($os.TotalVisibleMemorySize / 1KB)); Write-Output ('FreePhysicalMemory: {0:N2} MB' -f ($os.FreePhysicalMemory / 1KB)); Write-Output ('TotalVirtualMemorySize: {0:N2} MB' -f ($os.TotalVirtualMemorySize / 1KB)); Write-Output ('FreeVirtualMemory: {0:N2} MB' -f ($os.FreeVirtualMemory / 1KB))\""
	return c.ExecuteCommand(ctx, command)
}

// GetMetricGroupData gets data for a specific metric group
func (c *Client) GetMetricGroupData(ctx context.Context, metricGroup models.MetricGroup) (string, error) {
	var command string

	switch metricGroup.Name {
	case "CPU":
		return c.GetCPUInfo(ctx)
	case "MEMORY":
		return c.GetMemoryInfo(ctx)
	default:
		command = "systeminfo"
	}

	return c.ExecuteCommand(ctx, command)
}
