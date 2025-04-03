package llm

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

func GenerateTool[T any](name string, description string) (*anthropic.ToolParam, *anthropic.ToolChoiceToolParam) {
	tool := &anthropic.ToolParam{
		Name:        name,
		Description: anthropic.String(description),
		InputSchema: generateSchema[T](),
	}

	toolChoice := &anthropic.ToolChoiceToolParam{
		Type: "tool",
		Name: name,
	}

	return tool, toolChoice
}

func generateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}
