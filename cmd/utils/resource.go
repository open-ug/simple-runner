package utils

import (
	t "github.com/open-ug/conveyor/pkg/types"
)

var FlutterBuilderResourceDefinition = &t.ResourceDefinition{
	Name:        "flutter-builder",
	Description: "FFlutter builder resource definition",
	Version:     "1.0.0",
	Schema: map[string]interface{}{
		"properties": map[string]interface{}{
			"repository": map[string]interface{}{
				"type": "string",
			},
			"env": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"value": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"name", "value"},
				},
			},
		},
		"required": []string{"env"},
	},
}

type FlutterBuilderResourceSpec struct {
	Repository string                      `json:"repository"`
	Env        []FlutterBuilderResourceEnv `json:"env"`
}

type FlutterBuilderResourceEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type FlutterBuilderResource struct {
	Name     string                     `json:"name"`
	Resource string                     `json:"resource"`
	Spec     FlutterBuilderResourceSpec `json:"spec"`
	Metadata map[string]interface{}     `json:"metadata"`
}
