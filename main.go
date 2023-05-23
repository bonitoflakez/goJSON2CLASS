package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Schema struct {
	Title      string                 `json:"title"`
	Properties map[string]interface{} `json:"properties"`
	Items      *Schema                `json:"items"`
}

type RustType struct {
	Name     string
	DataType string
}

func main() {
	// Read command-line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: goJSON2TYPES <Schema.json> <Output.rs>")
		os.Exit(1)
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]

	schema, err := readJSONSchema(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	rustCode := generateRustCode(schema)

	err = os.WriteFile(outputFile, []byte(rustCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}

func generateRustCode(schema *Schema) string {
	var builder strings.Builder
	processSchema(&builder, schema, "")
	return builder.String()
}

func readJSONSchema(filePath string) (*Schema, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var schema Schema
	err = json.Unmarshal(data, &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}

	return &schema, nil
}

func getRustType(data interface{}) string {
	switch t := data.(type) {
	case *Schema:
		if t.Properties != nil {
			return t.Title
		} else if t.Items != nil {
			return "Vec<" + getRustType(t.Items) + ">"
		}
	case map[string]interface{}:
		dataType, ok := t["type"].(string)
		if !ok {
			return "unknown"
		}
		switch dataType {
		case "integer":
			return "i64"
		case "number":
			return "f64"
		case "boolean":
			return "bool"
		case "string":
			return "String"
		case "array":
			items, ok := t["items"].(map[string]interface{})
			if ok {
				return "Vec<" + getRustType(items) + ">"
			}
		case "object":
			title, ok := t["title"].(string)
			if ok {
				return title
			}
		}
	}

	return "unknown"
}

func processSchema(builder *strings.Builder, schema *Schema, indent string) {
	if schema.Properties != nil {
		builder.WriteString(indent + "struct " + schema.Title + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t" + name + ": " + getRustType(property) + ",\n")
		}
		builder.WriteString(indent + "}\n\n")

		// handle nested objects within object properties
		for _, name := range propertyNames {
			property := schema.Properties[name]
			if propertyMap, ok := property.(map[string]interface{}); ok {
				if nestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
					nestedTitle, ok := propertyMap["title"].(string)
					if !ok {
						nestedTitle = name
					}
					nestedPropertyMap := nestedSchema
					nestedSchema := &Schema{
						Title:      nestedTitle,
						Properties: nestedPropertyMap,
					}
					processNestedObjects(builder, nestedSchema, indent+"\t", nestedTitle)
				}
			}
		}
	} else if schema.Items != nil {
		// handle array items
		builder.WriteString(indent + "struct " + schema.Title + " {\n")
		builder.WriteString(indent + "\t" + "items: Vec<" + getRustType(schema.Items) + ">,\n")
		builder.WriteString(indent + "}\n\n")

		// handle nested objects within array items
		processNestedObjects(builder, schema.Items, indent+"\t", schema.Items.Title)
	}
}

func processNestedObjects(builder *strings.Builder, schema *Schema, indent string, structName string) {
	if schema.Properties != nil {
		builder.WriteString(indent + "struct " + structName + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t" + name + ": " + getRustType(property) + ",\n")
		}
		builder.WriteString(indent + "}\n\n")

		// Process nested objects within nested properties
		for _, name := range propertyNames {
			property := schema.Properties[name]
			if propertyMap, ok := property.(map[string]interface{}); ok {
				if nestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
					nestedTitle, ok := propertyMap["title"].(string)
					if !ok {
						nestedTitle = name
					}
					nestedPropertyMap := nestedSchema
					nestedSchema := &Schema{
						Title:      nestedTitle,
						Properties: nestedPropertyMap,
					}
					processNestedObjects(builder, nestedSchema, indent+"\t", nestedTitle)
				}
			}
		}
	}
}
