package main

import (
	"sort"
	"strings"
)

func generateRustCode(schema *Schema, pubFlag bool) string {
	var builder strings.Builder
	processSchemaForRust(&builder, schema, "\t", pubFlag)
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
		builder.WriteString("#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString("pub struct " + getFirstWordFromTitle(schema.Title) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			serdeAnnotation := "#[serde(rename = \"" + name + "\")]\n"
			declaration := getPropertyDeclaration(name, getRustType(property), pubFlag)
			builder.WriteString(indent + serdeAnnotation + indent + declaration + ",\n")
		}
		builder.WriteString("}\n\n")

		for _, name := range propertyNames {
			property := schema.Properties[name]
			if propertyMap, ok := property.(map[string]interface{}); ok {
				if nestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
					nestedTitle := propertyMap["title"].(string)
					if nestedTitle == "" {
						nestedTitle = name
					}
					nestedPropertyMap := nestedSchema
					nestedSchema := &Schema{
						Title:      nestedTitle,
						Properties: nestedPropertyMap,
					}
					processNestedObjectsForRust(builder, nestedSchema, indent, nestedTitle, pubFlag)
				}
			}
		}
	} else if schema.Items != nil {
		builder.WriteString("#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString("pub struct " + getFirstWordFromTitle(schema.Title) + " {\n")
		serdeAnnotation := "#[serde(rename = \"items\")]\n"
		declaration := getPropertyDeclaration("items", "Vec<"+getRustType(schema.Items)+">", pubFlag)
		builder.WriteString(indent + serdeAnnotation + indent + declaration + ",\n")
		builder.WriteString("}\n\n")

		processNestedObjectsForRust(builder, schema.Items, indent, schema.Items.Title, pubFlag)
	}
}

func processNestedObjectsForRust(builder *strings.Builder, schema *Schema, indent string, structName string, pubFlag bool) {
	if schema.Properties != nil {
		builder.WriteString("#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString("pub struct " + getFirstWordFromTitle(structName) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			serdeAnnotation := "#[serde(rename = \"" + name + "\")]\n"
			declaration := getPropertyDeclaration(name, getRustType(property), pubFlag)
			builder.WriteString(indent + serdeAnnotation + indent + declaration + ",\n")
		}
		builder.WriteString("}\n\n")

		for _, name := range propertyNames {
			property := schema.Properties[name]
			if propertyMap, ok := property.(map[string]interface{}); ok {
				if nestedSchema, ok := propertyMap["properties"].(map[string]interface{}); ok {
					nestedTitle := propertyMap["title"].(string)
					if nestedTitle == "" {
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
