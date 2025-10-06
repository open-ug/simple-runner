package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	log "github.com/open-ug/conveyor/pkg/driver-runtime/log"
	types "github.com/open-ug/conveyor/pkg/types"
	utils "github.com/open-ug/runner/cmd/utils"
)

func Reconcile(payload string, event string, driverName string, logger *log.DriverLogger) types.DriverResult {
	// Unmarshal the payload into a PipelineResource
	var pipelineResource utils.PipelineResource
	if err := json.Unmarshal([]byte(payload), &pipelineResource); err != nil {
		fmt.Println("Error unmarshalling payload: ", err)
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

	// 1. Create and start the container
	containerID, err := utils.CreateAndStartContainer(ctx, cli, pipelineResource.Spec.Image, nil)
	if err != nil {
		return types.DriverResult{
			Success: false,
			Message: fmt.Sprintf("failed to create/start container: %v", err),
		}
	}
	logger.Log(nil, fmt.Sprintf("Started container %s for pipeline %s", containerID, pipelineResource.Name))

	// 2. Execute all pipeline steps sequentially
	for _, step := range pipelineResource.Spec.Steps {
		logger.Log(map[string]string{
			"step": step.Name,
		}, fmt.Sprintf("Executing step: %s", step.Name))

		// Split the command string into args (simple split, adjust if needed for complex commands)
		command := strings.Fields(step.Command)

		exitCode, err := utils.ExecInContainer(ctx, cli, containerID, command, *logger)
		if err != nil {
			logger.Log(map[string]string{"step": step.Name}, fmt.Sprintf("Error executing step: %v", err))

			return types.DriverResult{
				Success: false,
				Message: fmt.Sprintf("Error executing step %s: %v", step.Name, err),
			}
		}

		if exitCode != 0 {
			logger.Log(map[string]string{"step": step.Name}, fmt.Sprintf("Step failed with exit code %d", exitCode))
			//_ = cli.ContainerStop(ctx, containerID, nil) // ensure container is stopped
			return types.DriverResult{
				Success: false,
				Message: fmt.Sprintf("step %s failed with exit code %d", step.Name, exitCode),
			}
		}

		logger.Log(map[string]string{"step": step.Name}, "Step completed successfully")
	}

	// 3. Stop the container after all steps are done

	logger.Log(nil, fmt.Sprintf("Stopped container %s", containerID))

	return types.DriverResult{
		Success: true,
		Message: "Pipeline executed successfully",
	}
}
