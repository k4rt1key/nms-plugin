package winrm

import (
	"context"
	"fmt"
	"nms-plugin/src/logger"
	"time"

	"github.com/masterzen/winrm"
)

const (
	// DefaultTimeout is the default timeout for WinRM operations
	DefaultTimeout = 10 * time.Second
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
	endpoint := winrm.NewEndpoint(c.IP, c.Port, false, false, nil, nil, nil, DefaultTimeout)

	client, err := winrm.NewClient(endpoint, c.Username, c.Password)

	if err != nil {
		logger.Error("Failed to create WinRM client for %s:%d: %v", c.IP, c.Port, err)
		return "", fmt.Errorf("failed to create WinRM client: %w", err)
	}

	// Execute the command with context
	stdout, stderr, exitCode, err := client.RunCmdWithContext(ctx, command)

	if err != nil {
		logger.Error("Command execution failed on %s:%d: %v", c.IP, c.Port, err)
		return stderr, err
	}

	if exitCode != 0 {
		logger.Error("Command exited with code %d on %s:%d", exitCode, c.IP, c.Port)
		return stderr, fmt.Errorf("command exited with code %d", exitCode)
	}

	logger.Debug("Command execution succeeded on %s:%d", c.IP, c.Port)

	return stdout, nil
}

// TestConnection tests if the WinRM connection is successful
func (c *Client) TestConnection(ctx context.Context) (bool, string, error) {
	output, err := c.ExecuteCommand(ctx, "hostname")

	if err != nil {
		return false, "", err
	}

	return true, output, nil
}
