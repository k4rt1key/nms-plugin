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

func TestConnection(ctx context.Context, client *winrm.Client) (bool, string) {

	stdout, stderr, exitCode, err := client.RunCmdWithContext(ctx, "hostname")

	if err != nil || exitCode != 0 {

		if stderr != "" {

			return false, stderr

		}

		if err != nil {

			return false, err.Error()

		}

		return false, fmt.Sprintf("exit code %d", exitCode)
	}

	return true, stdout
}

func ExecuteCommand(ctx context.Context, client *winrm.Client, command string) (string, error) {

	stdout, stderr, exitCode, err := client.RunCmdWithContext(ctx, command)

	if err != nil {

		return stderr, err

	}

	if exitCode != 0 {

		return stderr, fmt.Errorf("exit code %d", exitCode)

	}

	return stdout, nil

}
