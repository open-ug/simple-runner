package main

import (
	"context"
	"fmt"
	"time"

	c "github.com/open-ug/conveyor/pkg/client"
	runtime "github.com/open-ug/conveyor/pkg/driver-runtime"
	cmd "github.com/open-ug/runner/cmd"
	utils "github.com/open-ug/runner/cmd/utils"
)

func main() {
	client, err := c.NewClient(
		"http://localhost:8080",
		false,
		"",
		"",
		30*time.Second)
	if err != nil {
		panic(err)
	}

	// Register the Pipeline resource definition with the client
	_, err = client.CreateOrUpdateResourceDefinition(context.Background(), utils.PipelineResourceDefinition)

	if err != nil {
		panic(err)
	}

	// Create a new driver instance
	driver := &runtime.Driver{
		Name: "command-runner",
		Resources: []string{
			utils.PipelineResourceDefinition.Name},
		Reconcile: cmd.Reconcile,
	}

	// Create a new driver manager with the driver
	driverManager, err := runtime.NewDriverManager(driver, []string{"*"})
	if err != nil {
		fmt.Println("Error creating driver manager: ", err)
		return
	}

	// Start the driver manager
	err = driverManager.Run()
	if err != nil {
		fmt.Println("Error running driver manager: ", err)
	}
}
