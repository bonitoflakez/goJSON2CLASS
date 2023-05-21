package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func readJSONSchema(filePath string) (*Schema, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
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

func generateRustCode(schema *Schema) string {
	var builder strings.Builder

	processSchema(&builder, schema, "")

	return builder.String()
}

func processSchema(builder *strings.Builder, schema *Schema, indent string) {
	if schema.Properties != nil {
		// Process object properties
		builder.WriteString(indent + "struct " + schema.Title + " {\n")
		for name, property := range schema.Properties {
			builder.WriteString(indent + "\t" + name + ": " + getRustType(property) + ",\n")
		}
		builder.WriteString(indent + "}\n\n")
	} else if schema.Items != nil {
		// Process array items
		builder.WriteString(indent + "struct " + schema.Title + " {\n")
		builder.WriteString(indent + "\t" + "items: Vec<" + getRustType(schema.Items) + ">,\n")
		builder.WriteString(indent + "}\n\n")
	}

	// Process nested objects
	for _, property := range schema.Properties {
		if propertyMap, ok := property.(map[string]interface{}); ok {
			if nestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
				for nestedTitle, nestedProperty := range nestedSchema {
					if nestedPropertyMap, ok := nestedProperty.(map[string]interface{}); ok {
						nestedSchema := &Schema{
							Title:      nestedTitle,
							Properties: nestedPropertyMap,
						}
						processSchema(builder, nestedSchema, indent+"\t")
					}
				}
			}
		}
	}
}

func getRustType(data interface{}) string {
	switch t := data.(type) {
	case map[string]interface{}:
		if t["type"] == "array" {
			items, ok := t["items"].(map[string]interface{})
			if ok {
				return "Vec<" + getRustType(items) + ">"
			}
		} else if t["type"] == "object" {
			// Nested object
			title, ok := t["title"].(string)
			if ok {
				return title
			}
		}
	case string:
		switch t {
		case "integer":
			return "i64"
		case "number":
			return "f64"
		case "boolean":
			return "bool"
		case "string":
			return "String"
		}
	}

	return "unknown"
}

func main() {
	// Read command-line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: goJSON2TYPES <input_schema.json> <output.rs>")
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

	err = ioutil.WriteFile(outputFile, []byte(rustCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}
