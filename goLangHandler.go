package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type GoType struct {
	Name     string
	DataType string
}

func writeGoCodeToFile(outFile string, goCode string) {
	err := os.WriteFile(outFile, []byte(goCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func generateGoCode(schema *Schema) string {
	var builder strings.Builder
	processSchemaForGo(&builder, schema, "")
	return builder.String()
}

func getGoType(data interface{}) string {
	switch t := data.(type) {
	case *Schema:
		if t.Properties != nil {
			return t.Title
		} else if t.Items != nil {
			return "[]" + getGoType(t.Items)
		}
	case map[string]interface{}:
		dataType, ok := t["type"].(string)
		if !ok {
			return "unknown"
		}
		switch dataType {
		case "integer":
			return "int64"
		case "number":
			return "float64"
		case "boolean":
			return "bool"
		case "string":
			return "string"
		case "array":
			items, ok := t["items"].(map[string]interface{})
			if ok {
				return "[]" + getGoType(items)
			}
		case "object":
			title, ok := t["title"].(string)
			if ok {
				return getFirstWordFromTitle(title)
			}
		}
	}

	return "unknown"
}

func processSchemaForGo(builder *strings.Builder, schema *Schema, indent string) {
	if schema.Properties != nil {
		builder.WriteString("package main\n\n")
		builder.WriteString(indent + "type " + getFirstWordFromTitle(schema.Title) + " struct {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t" + name + " " + getGoType(property) + "\n")
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
					processNestedObjectsForGo(builder, nestedSchema, indent+"", nestedTitle)
				}
			}
		}
	} else if schema.Items != nil {
		// handle array items
		builder.WriteString(indent + "type " + getFirstWordFromTitle(schema.Title) + " struct {\n")
		builder.WriteString(indent + "\t" + "[]" + getGoType(schema.Items) + "\n")
		builder.WriteString(indent + "}\n\n")

		// handle nested objects within array items
		processNestedObjectsForGo(builder, schema.Items, indent+"", schema.Items.Title)
	}
}

func processNestedObjectsForGo(builder *strings.Builder, schema *Schema, indent string, structName string) {
	if schema.Properties != nil {
		builder.WriteString(indent + "type " + getFirstWordFromTitle(structName) + " struct {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t" + name + " " + getGoType(property) + "\n")
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
					processNestedObjectsForGo(builder, nestedSchema, indent+"", nestedTitle)
				}
			}
		}
	}
}
