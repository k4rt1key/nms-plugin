package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DiscoveryResult struct {
	Type   string `json:"type"`
	ID     int    `json:"id"`
	Result []struct {
		Success    bool   `json:"success"`
		IP         string `json:"ip"`
		Port       int    `json:"port"`
		Message    string `json:"message"`
		Credential struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"credential"`
	} `json:"result"`
}

func ParseDiscoveryOutput(s string) string {
	var discovery DiscoveryResult
	err := json.Unmarshal([]byte(s), &discovery)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %v", err)
	}

	var output strings.Builder
	output.WriteString("Discovered Windows Systems:\n")

	for _, result := range discovery.Result {
		if !result.Success {
			continue
		}

		// Extract hostname
		hostname := "Unknown"
		if hostLine := findLineStartingWith(result.Message, "Host Name:"); hostLine != "" {
			hostname = strings.TrimSpace(strings.SplitN(hostLine, ":", 2)[1])
		}

		// Extract OS name
		osName := "Unknown Windows OS"
		if osLine := findLineStartingWith(result.Message, "OS Name:"); osLine != "" {
			osName = strings.TrimSpace(strings.SplitN(osLine, ":", 2)[1])
		}

		output.WriteString(fmt.Sprintf("- %s (%s) - IP: %s\n", hostname, osName, result.IP))
	}

	return output.String()
}

func findLineStartingWith(text, prefix string) string {
	lines := strings.Split(text, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}
