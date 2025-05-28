package polling

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"nms-plugin/src/winrm"
)

const timeout = 60 * time.Second

func Poll(request map[string]interface{}) {

	metricGroups := request["metric_groups"].([]interface{})

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	/*
		TODO: TESTING
	*/

	//for _ = range 2 {
	//	metricGroups = append(metricGroups, metricGroups...)
	//}

	resultChan := make(chan map[string]interface{})

	var wg sync.WaitGroup

	wg.Add(len(metricGroups))

	for _, mgInterface := range metricGroups {

		go func(mg map[string]interface{}) {

			defer wg.Done()

			protocol, ok := (mg["credential"].(map[string]interface{}))["protocol"]

			if !ok {
				protocol = "winrm"
			}

			result := map[string]interface{}{
				"monitor_id": mg["monitor_id"],
				"name":       mg["name"],
				"time":       time.Now().Format(time.RFC3339),
			}

			switch protocol {

			case "winrm":
				winrm.Poll(ctx, mg, result)

			default:
				// Default to WinRM if protocol not recognized
				winrm.Poll(ctx, mg, result)

			}

			resultChan <- result

		}(mgInterface.(map[string]interface{}))

	}

	go func() {

		wg.Wait()

		close(resultChan)

	}()

	for result := range resultChan {

		output, _ := json.Marshal(result)

		fmt.Println(string(output))

	}

}
