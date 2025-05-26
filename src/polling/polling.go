package polling

import (
	"context"
	"sync"
	"time"

	"nms-plugin/src/winrm"
)

const timeout = 60 * time.Second

func Execute(request map[string]interface{}) map[string]interface{} {

	metricGroups := request["metric_groups"].([]interface{})

	response := map[string]interface{}{
		"type":          "polling",
		"metric_groups": []map[string]interface{}{},
	}

	results := poll(metricGroups)

	response["metric_groups"] = results

	return response
}

func getProtocolFromCredential(credential map[string]interface{}) string {

	if protocol, ok := credential["protocol"]; ok {

		return protocol.(string)

	}

	return "winrm" // default protocol
}

func poll(metricGroups []interface{}) []map[string]interface{} {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	resultChan := make(chan map[string]interface{})

	var wg sync.WaitGroup

	wg.Add(len(metricGroups))

	for _, mgInterface := range metricGroups {

		go func(mg map[string]interface{}) {

			defer wg.Done()

			protocol := getProtocolFromCredential(mg["credential"].(map[string]interface{}))

			result := map[string]interface{}{
				"monitor_id": mg["monitor_id"],
				"name":       mg["name"],
				"time":       time.Now().Format(time.RFC3339),
			}

			switch protocol {

			case "winrm":
				winrm.Execute(ctx, mg, result)

			default:
				// Default to WinRM if protocol not recognized
				winrm.Execute(ctx, mg, result)

			}

			resultChan <- result

		}(mgInterface.(map[string]interface{}))
	}

	go func() {

		wg.Wait()

		close(resultChan)

	}()

	var results []map[string]interface{}

	for result := range resultChan {

		results = append(results, result)

	}

	return results
}
