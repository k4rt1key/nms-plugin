// discovery/discovery.go
package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"nms-plugin/src/winrm"
)

const timeout = 30 * time.Second

func getProtocolFromCredential(credential map[string]interface{}) string {

	if protocol, ok := credential["protocol"]; ok {

		return protocol.(string)

	}

	return "winrm"
}

func Discover(request map[string]interface{}) {

	ips := request["ips"].([]interface{})

	credentials := request["credentials"].([]interface{})

	port := int(request["port"].(float64))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	resultChan := make(chan map[string]interface{})

	var wg sync.WaitGroup

	successMap := make(map[string]map[string]interface{})

	var mu sync.Mutex

	totalJobs := len(ips) * len(credentials)

	wg.Add(totalJobs)

	for _, ipInterface := range ips {

		for _, credInterface := range credentials {

			go func(ip string, cred map[string]interface{}) {

				defer wg.Done()

				protocol, ok := cred["protocol"]

				if !ok {
					protocol = "winrm"
				}

				var success bool

				var message string

				switch protocol {

				case "winrm":
					success, message = winrm.TestConnection(ctx, ip, port, cred)

				default:
					// Default to WinRM if protocol not recognized
					success, message = winrm.TestConnection(ctx, ip, port, cred)
				}

				result := map[string]interface{}{
					"ip":         ip,
					"credential": cred,
					"success":    success,
					"message":    message,
					"port":       port,
					"protocol":   protocol,
				}

				resultChan <- result

			}(ipInterface.(string), credInterface.(map[string]interface{}))
		}
	}

	go func() {

		wg.Wait()

		close(resultChan)

	}()

	for result := range resultChan {

		if result["success"].(bool) {

			mu.Lock()

			// Store the successful result for this IP, but prefer the first successful one
			ip := result["ip"].(string)

			if _, exists := successMap[ip]; !exists {

				successMap[ip] = result

				output, _ := json.Marshal(result)

				fmt.Println(string(output))

			}

			mu.Unlock()

		}
	}

	for _, ipInterface := range ips {

		ip := ipInterface.(string)

		if _, exists := successMap[ip]; !exists {

			output, _ := json.Marshal(map[string]interface{}{
				"success":    false,
				"ip":         ip,
				"credential": map[string]interface{}{},
				"port":       port,
				"protocol":   "",
				"message":    "No valid credential found",
			})

			fmt.Println(string(output))

		}
	}

}
