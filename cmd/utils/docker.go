package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/open-ug/conveyor/pkg/driver-runtime/log"
)

// ExecInContainer executes a command inside a Docker container.
// Parameters:
//   - ctx: context for cancellation
//   - cli: Docker client
//   - containerName: Name or ID of the container
//   - command: Slice of command arguments (e.g., []string{"ls", "-l", "/"})
//
// Returns:
//   - exitCode: Exit status of the executed command
//   - error: Any error encountered during execution
func ExecInContainer(ctx context.Context, cli *client.Client, containerName string, command []string, logger log.DriverLogger) (int, error) {
	// Create an exec instance configuration
	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	}

	// Create the exec instance in the container
	execIDResp, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return -1, fmt.Errorf("failed to create exec instance: %w", err)
	}

	// Attach to the exec instance to capture logs
	resp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return -1, fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer resp.Close()

	// Read logs line by line and print to stdout
	reader := bufio.NewReader(resp.Reader)
	for {
		line, err := reader.ReadString('\n')
		if len(strings.TrimSpace(line)) > 0 {
			fmt.Print(line) // print log line
			logger.Log(map[string]string{
				"command": command[0],
			}, line) // log the line using the provided logger
		}
		if err != nil {
			if err == io.EOF {
				break // End of logs
			}
			return -1, fmt.Errorf("error reading logs: %w", err)
		}
	}

	// Inspect exec to get the exit code
	inspectResp, err := cli.ContainerExecInspect(ctx, execIDResp.ID)
	if err != nil {
		return -1, fmt.Errorf("failed to inspect exec instance: %w", err)
	}

	return inspectResp.ExitCode, nil
}

// CreateAndStartContainer pulls an image if needed, creates and starts a container, and returns its container ID.
// Parameters:
//   - ctx: context for cancellation
//   - cli: Docker client
//   - image: The Docker image to run (e.g., "alpine:latest")
//   - cmd: Optional command to override the image's default command
//
// Returns:
//   - containerID: ID of the created container
//   - err: Any error encountered
func CreateAndStartContainer(ctx context.Context, cli *client.Client, image string, cmd []string) (string, error) {

	// Create the container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   cmd,
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// StopContainer stops a running container by ID or name.
func StopContainer(ctx context.Context, cli *client.Client, containerID string) error {
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}
