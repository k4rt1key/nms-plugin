// discovery/discovery.go
package discovery

import (
	"context"
	"sync"
	"time"

	"nms-plugin/src/winrm"
)

const timeout = 30 * time.Second

func Execute(request map[string]interface{}) map[string]interface{} {

	ips := request["ips"].([]interface{})

	credentials := request["credentials"].([]interface{})

	port := int(request["port"].(float64))

	protocol := getProtocol(request)

	response := map[string]interface{}{
		"type":    "discovery",
		"id":      request["id"],
		"results": []map[string]interface{}{},
	}

	switch protocol {

	case "winrm":

		results := discoverWinRM(ips, credentials, port)

		response["results"] = results

	default:

		results := discoverWinRM(ips, credentials, port)

		response["results"] = results
	}

	return response
}

func getProtocol(request map[string]interface{}) string {

	if protocol, ok := request["protocol"]; ok {

		return protocol.(string)

	}

	return "winrm"
}

func discoverWinRM(ips, credentials []interface{}, port int) []map[string]interface{} {

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

				client := winrm.NewClient(ip, port, cred["username"].(string), cred["password"].(string))

				success, message := winrm.TestConnection(ctx, client)

				result := map[string]interface{}{
					"ip":         ip,
					"credential": cred,
					"success":    success,
					"message":    message,
					"port":       port,
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

			successMap[result["ip"].(string)] = result

			mu.Unlock()

		}
	}

	var results []map[string]interface{}

	for _, ipInterface := range ips {

		ip := ipInterface.(string)

		if result, exists := successMap[ip]; exists {

			results = append(results, result)

		} else {

			results = append(results, map[string]interface{}{
				"success":    false,
				"ip":         ip,
				"credential": map[string]interface{}{},
				"port":       port,
				"message":    "Connection failed",
			})

		}
	}

	return results
}
