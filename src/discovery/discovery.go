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

	totalJobs := len(ips) * len(credentials)

	wg.Add(totalJobs)

	for _, ipInterface := range ips {

		for _, credInterface := range credentials {

			go func(ip string, cred map[string]interface{}) {

				defer wg.Done()

				protocol := getProtocolFromCredential(cred)

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
					"type":       "discovery",
					"id":         request["id"],
					"ip":         ip,
					"credential": cred,
					"success":    success,
					"message":    message,
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

			output, _ := json.Marshal(result)

			fmt.Println(string(output))

		} else {

			output, _ := json.Marshal(map[string]interface{}{
				"id":         request["id"],
				"type":       "discovery",
				"success":    false,
				"ip":         result["ip"],
				"credential": map[string]interface{}{},
				"message":    request["message"],
			})

			fmt.Println(string(output))

		}
	}
}
