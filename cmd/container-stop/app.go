package containerstop

import (
	//"encoding/json"

	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/docker/docker/client"
	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
	"github.com/open-ug/conveyor/pkg/types"
	"github.com/open-ug/runner/cmd/utils"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) types.DriverResult {

	log.SetFlags(log.Ldate | log.Ltime)

	var flutterResource utils.FlutterBuilderResource
	err := json.Unmarshal([]byte(payload), &flutterResource)
	if err != nil {
		return types.DriverResult{
			Success: false,
			Message: fmt.Sprintf("Error unmarshalling payload: %v", err),
		}
	}

	// Initialize Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return types.DriverResult{
			Success: false,
			Message: fmt.Sprintf("failed to create Docker client: %v", err),
		}
	}

	containerStart, ok := flutterResource.Metadata["driverresults.container-start"].(map[string]interface{})
	if !ok {
		return types.DriverResult{
			Success: false,
			Message: "Invalid type assertion for driverresults.container-start",
		}
	}

	data, ok := containerStart["data"].(map[string]interface{})
	if !ok {
		return types.DriverResult{
			Success: false,
			Message: "Invalid type assertion for data",
		}
	}

	containerId, ok := data["containerID"].(string)
	if !ok {
		return types.DriverResult{
			Success: false,
			Message: "Invalid type assertion for containerID",
		}
	}

	err = utils.StopContainer(ctx, cli, containerId)
	if err != nil {
		return types.DriverResult{
			Success: false,
			Message: fmt.Sprintf("failed to stop container: %v", err),
		}
	}

	return types.DriverResult{
		Success: true,
		Message: "Sample Driver Reconciled Successfully",
	}
}

func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
		Name:      "container-stop",
		Resources: []string{utils.FlutterBuilderResourceDefinition.Name},
	}

	driverManager, err := runtime.NewDriverManager(driver, []string{"*"})

	if err != nil {
		fmt.Println("Error creating driver manager: ", err)
		return
	}

	err = driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}

}
