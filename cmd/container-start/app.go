package containerstart

import (
	//"encoding/json"

	"context"
	"fmt"

	"github.com/docker/docker/client"
	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	logger "github.com/open-ug/conveyor/pkg/driver-runtime/log"
	"github.com/open-ug/conveyor/pkg/types"
	"github.com/open-ug/runner/cmd/utils"
)

// Listen for messages from the runtime
func Reconcile(payload string, event string, runID string, logger *logger.DriverLogger) types.DriverResult {

	// Initialize Docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return types.DriverResult{
			Success: false,
			Message: fmt.Sprintf("failed to create Docker client: %v", err),
		}
	}

	// 1. Create and start the container
	containerID, err := utils.CreateAndStartContainer(ctx, cli, "ghcr.io/cirruslabs/flutter:3.38.5", nil)
	if err != nil {
		return types.DriverResult{
			Success: false,
			Message: fmt.Sprintf("failed to create/start container: %v", err),
		}
	}

	return types.DriverResult{
		Success: true,
		Message: "Sample Driver Reconciled Successfully",
		Data: map[string]interface{}{
			"containerID": containerID,
		},
	}
}

func Listen() {
	driver := &runtime.Driver{
		Reconcile: Reconcile,
		Name:      "container-start",
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
