package main

import (
	"context"
	"time"

	c "github.com/open-ug/conveyor/pkg/client"
	"github.com/open-ug/runner/cmd/cli"
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
	_, err = client.CreateOrUpdateResourceDefinition(context.Background(), utils.FlutterBuilderResourceDefinition)

	if err != nil {
		panic(err)
	}

	cli.Execute()
}
