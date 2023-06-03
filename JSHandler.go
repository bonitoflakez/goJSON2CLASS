package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type JSType struct {
	Name     string
	DataType string
}

func writeJSCodeToFile(outFile string, jsCode string) {
	err := os.WriteFile(outFile, []byte(jsCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func generateJSCode(schema *Schema) string {
	var builder strings.Builder
	processSchemaForJS(&builder, schema, "")
	return builder.String()
}

func getJSType(data interface{}) string {
	switch t := data.(type) {
	case *Schema:
		if t.Properties != nil {
			return t.Title
		} else if t.Items != nil {
			return getJSType(t.Items) + "[]"
		}
	case map[string]interface{}:
		dataType, ok := t["type"].(string)
		if !ok {
			return "unknown"
		}
		switch dataType {
		case "integer":
			return "number"
		case "number":
			return "number"
		case "boolean":
			return "boolean"
		case "string":
			return "string"
		case "array":
			items, ok := t["items"].(map[string]interface{})
			if ok {
				return getJSType(items) + "[]"
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

func processSchemaForJS(builder *strings.Builder, schema *Schema, indent string) {
	builder.WriteString("class " + getFirstWordFromTitle(schema.Title) + " {\n")
	builder.WriteString(indent + "  constructor() {\n")

	if schema.Properties != nil {
		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			jsType := getJSType(property)
			builder.WriteString(indent + "    this." + getFirstWordFromTitle(name) + " = ")
			if jsType == "string" {
				builder.WriteString("''")
			} else if jsType == "number" {
				builder.WriteString("0")
			} else if jsType == "boolean" {
				builder.WriteString("false")
			} else {
				builder.WriteString("new " + jsType + "()")
			}
			builder.WriteString("; // " + jsType + " property\n")
		}
	}

	builder.WriteString(indent + "  }\n")

	if schema.Properties != nil {
		builder.WriteString("\n")
		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

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
					builder.WriteString(indent + "  get " + getFirstWordFromTitle(name) + "() {\n")
					builder.WriteString(indent + "    return new " + getFirstWordFromTitle(nestedTitle) + "();\n")
					builder.WriteString(indent + "  }\n\n")
					processSchemaForJS(builder, nestedSchema, indent+"  ")
				}
			}
		}
	}

	builder.WriteString("}\n\n")
}
