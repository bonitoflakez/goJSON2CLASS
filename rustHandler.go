package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type RustType struct {
	Name     string
	DataType string
}

func writeRustCodeToFile(outFile string, rustCode string) {
	err := os.WriteFile(outFile, []byte(rustCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func generateRustCode(schema *Schema, pubFlag bool) string {
	var builder strings.Builder
	processSchemaForRust(&builder, schema, "", pubFlag)
	return builder.String()
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

func processSchemaForRust(builder *strings.Builder, schema *Schema, indent string, pubFlag bool) {
	if schema.Properties != nil {
		builder.WriteString("use serde::{Serialize, Deserialize};\n\n")
		builder.WriteString(indent + "#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString(indent + "pub struct " + getFirstWordFromTitle(schema.Title) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			if pubFlag {
				builder.WriteString(indent + "\t#[serde(rename = \"" + name + "\")]\n")
				builder.WriteString(indent + "\tpub " + name + ": " + getRustType(property) + ",\n")
			} else {
				builder.WriteString(indent + "\t#[serde(rename = \"" + name + "\")]\n")
				builder.WriteString(indent + "\t" + name + ": " + getRustType(property) + ",\n")
			}
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
					processNestedObjectsForRust(builder, nestedSchema, indent+"", nestedTitle, pubFlag)
				}
			}
		}
	} else if schema.Items != nil {
		// handle array items
		if pubFlag {
			builder.WriteString(indent + "#[derive(Debug, Serialize, Deserialize)]\n")
			builder.WriteString(indent + "pub struct " + getFirstWordFromTitle(schema.Title) + " {\n")
			builder.WriteString(indent + "\t" + "#[serde(rename = \"items\")]\n")
			builder.WriteString(indent + "\tpub " + "items: Vec<" + getRustType(schema.Items) + ">,\n")
			builder.WriteString(indent + "}\n\n")
		} else {
			builder.WriteString(indent + "#[derive(Debug, Serialize, Deserialize)]\n")
			builder.WriteString(indent + "pub struct " + getFirstWordFromTitle(schema.Title) + " {\n")
			builder.WriteString(indent + "\t" + "#[serde(rename = \"items\")]\n")
			builder.WriteString(indent + "\t" + "items: Vec<" + getRustType(schema.Items) + ">,\n")
			builder.WriteString(indent + "}\n\n")
		}

		// handle nested objects within array items
		processNestedObjectsForRust(builder, schema.Items, indent+"", schema.Items.Title, pubFlag)
	}
}

func processNestedObjectsForRust(builder *strings.Builder, schema *Schema, indent string, structName string, pubFlag bool) {
	if schema.Properties != nil {
		builder.WriteString(indent + "#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString(indent + "pub struct " + getFirstWordFromTitle(structName) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			if pubFlag {
				builder.WriteString(indent + "\t#[serde(rename = \"" + name + "\")]\n")
				builder.WriteString(indent + "\tpub " + name + ": " + getRustType(property) + ",\n")
			} else {
				builder.WriteString(indent + "\t#[serde(rename = \"" + name + "\")]\n")
				builder.WriteString(indent + "\t" + name + ": " + getRustType(property) + ",\n")
			}
		}
		builder.WriteString(indent + "}\n\n")

		// handle nested objects within nested properties
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
					processNestedObjectsForRust(builder, nestedSchema, indent+"", nestedTitle, pubFlag)
				}
			}
		}
	}
}
