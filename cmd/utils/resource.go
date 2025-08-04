package utils

import (
	t "github.com/open-ug/conveyor/pkg/types"
)

// Pipeline Resource Definition. It contains an array of steps, each with a name and command.

var PipelineResourceDefinition = &t.ResourceDefinition{
	Name:        "pipeline",
	Description: "Pipeline resource definition",
	Version:     "1.0.0",
	Schema: map[string]interface{}{
		"properties": map[string]interface{}{
			"image": map[string]interface{}{
				"type": "string",
			},
			"steps": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"command": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"name", "command"},
				},
			},
		},
		"required": []string{"steps"},
	},
}

type PipelineResource struct {
	Name     string               `json:"name"`
	Resource string               `json:"resource"`
	Spec     PipelineResourceSpec `json:"spec"`
}

type PipelineResourceSpec struct {
	Image string         `json:"image"`
	Steps []PipelineStep `json:"steps"`
}

type PipelineStep struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}
