package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type CType struct {
	Title      string
	Properties map[string]interface{}
	Items      interface{}
}

func writeCCodeToFile(outFile string, CCode string) {
	err := os.WriteFile(outFile, []byte(CCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func generateCCode(schema *Schema) string {
	var builder strings.Builder

	builder.WriteString("#include <stdio.h>\n")
	if checkBool(schema) {
		builder.WriteString("#include <stdbool.h>\n\n")
	}

	processSchemaForC(&builder, schema, "")
	return builder.String()
}

func getCType(property interface{}) string {
	switch p := property.(type) {
	case map[string]interface{}:
		if pType, ok := p["type"].(string); ok {
			switch pType {
			case "string":
				return "char*"
			case "number":
				return "double"
			case "integer":
				return "int"
			case "boolean":
				return "bool"
			case "decimal":
				return "float"
			case "object":
				title, ok := p["title"].(string)
				if ok {
					return "struct " + title
				}
			}
		}
	}
	return "unknown"
}

func getItemCType(property interface{}) string {
	switch p := property.(type) {
	case map[string]interface{}:
		if pType, ok := p["type"].(string); ok {
			switch pType {
			case "string":
				return "char"
			case "number":
				return "double"
			case "integer":
				return "int"
			case "decimal":
				return "float"
			}
		}
	}
	return "unknown"
}

func checkBool(schema *Schema) bool {
	if schema.Properties != nil {
		for _, property := range schema.Properties {
			if propertyMap, ok := property.(map[string]interface{}); ok {
				if propertyType, ok := propertyMap["type"].(string); ok && propertyType == "boolean" {
					return true
				} else if nestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
					if checkNestedBool(&Schema{Properties: nestedSchema}) {
						return true
					}
				}
			}
		}
	}
	return false
}

func checkNestedBool(schema *Schema) bool {
	if schema.Properties != nil {
		for _, property := range schema.Properties {
			if propertyMap, ok := property.(map[string]interface{}); ok {
				if propertyType, ok := propertyMap["type"].(string); ok && propertyType == "boolean" {
					return true
				} else if nestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
					if checkNestedBool(&Schema{Properties: nestedSchema}) {
						return true
					}
				}
			}
		}
	}
	return false
}

func isArrayType(property interface{}) bool {
	if propertyMap, ok := property.(map[string]interface{}); ok {
		if propertyType, ok := propertyMap["type"].(string); ok && propertyType == "array" {
			return true
		}
	}
	return false
}

func getArrayType(property interface{}) interface{} {
	if propertyMap, ok := property.(map[string]interface{}); ok {
		if items, ok := propertyMap["items"]; ok {
			return items
		}
	}
	return nil
}

func processSchemaForC(builder *strings.Builder, schema *Schema, indent string) {
	if schema.Properties != nil {

		builder.WriteString(indent + "struct " + getFirstWordFromTitle(schema.Title) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			// builder.WriteString(indent + "\t" + getCType(property) + " " + name + ";\n")
			if isArrayType(property) {
				itemType := getArrayType(property)
				// default size 50
				// TODO: define size inside schema and replace value by that
				builder.WriteString(indent + "\t" + getItemCType(itemType) + " " + name + "[50];\n")
			} else {
				builder.WriteString(indent + "\t" + getCType(property) + " " + name + ";\n")
			}
		}
		builder.WriteString(indent + "};\n\n")

		// Handle nested objects within object properties
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
					processNestedObjectsForC(builder, nestedSchema, indent+"")
				}
			}
		}
	} else if schema.Items != nil {
		// Handle array items
		fmt.Println("using schema.items from processSchemaForC()")
		builder.WriteString(indent + "struct " + getFirstWordFromTitle(schema.Title) + " {\n")
		builder.WriteString(indent + "\t" + getCType(schema.Items) + " items;\n")
		builder.WriteString(indent + "};\n\n")

		// Handle nested objects within array items
		processNestedObjectsForC(builder, schema.Items, indent+"")
	}
}

func processNestedObjectsForC(builder *strings.Builder, schema interface{}, indent string) {
	if nestedSchema, ok := schema.(*Schema); ok {
		if nestedSchema.Properties != nil {
			builder.WriteString(indent + "struct " + getFirstWordFromTitle(nestedSchema.Title) + " {\n")

			var propertyNames []string
			for name := range nestedSchema.Properties {
				propertyNames = append(propertyNames, name)
			}
			sort.Strings(propertyNames)

			for _, name := range propertyNames {
				property := nestedSchema.Properties[name]
				// builder.WriteString(indent + "\t" + getCType(property) + " " + name + ";\n")
				if isArrayType(property) {
					itemType := getArrayType(property)
					// default size 50
					// TODO: define size inside schema and replace value by that
					builder.WriteString(indent + "\t" + getItemCType(itemType) + " " + name + "[50];\n")
				} else {
					builder.WriteString(indent + "\t" + getCType(property) + " " + name + ";\n")
				}
			}
			builder.WriteString(indent + "};\n\n")

			// handle nested objects within nested properties
			for _, name := range propertyNames {
				property := nestedSchema.Properties[name]
				if propertyMap, ok := property.(map[string]interface{}); ok {
					if subNestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
						subNestedTitle, ok := propertyMap["title"].(string)
						if !ok {
							subNestedTitle = name
						}
						subNestedPropertyMap := subNestedSchema
						subNestedSchema := &Schema{
							Title:      subNestedTitle,
							Properties: subNestedPropertyMap,
						}
						processNestedObjectsForC(builder, subNestedSchema, indent+"")
					}
				}
			}
		}
	}
}
