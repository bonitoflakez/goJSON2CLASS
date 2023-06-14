package main

import (
	"sort"
	"strings"
)

func generateTSCode(schema *Schema) string {
	var builder strings.Builder
	processSchemaForTS(&builder, schema, "")
	return builder.String()
}

func getTSType(data interface{}) string {
	switch t := data.(type) {
	case *Schema:
		if t.Properties != nil {
			return t.Title
		} else if t.Items != nil {
			return getTSType(t.Items) + "[]"
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
				return getTSType(items) + "[]"
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

func processSchemaForTS(builder *strings.Builder, schema *Schema, indent string) {
	if schema.Properties != nil {
		builder.WriteString(indent + "interface " + getFirstWordFromTitle(schema.Title) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t" + name + ": " + getTSType(property) + ",\n")
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
					processNestedObjectsForTS(builder, nestedSchema, indent+"", nestedTitle)
				}
			}
		}
	} else if schema.Items != nil {
		// handle array items
		builder.WriteString(indent + "interface " + getFirstWordFromTitle(schema.Title) + " {\n")
		builder.WriteString(indent + "\t" + getTSType(schema.Items) + "[]" + "\n")
		builder.WriteString(indent + "}\n\n")

		// handle nested objects within array items
		processNestedObjectsForTS(builder, schema.Items, indent+"", schema.Items.Title)
	}
}

func processNestedObjectsForTS(builder *strings.Builder, schema *Schema, indent string, structName string) {
	if schema.Properties != nil {
		builder.WriteString(indent + "interface " + getFirstWordFromTitle(structName) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t" + name + ": " + getTSType(property) + ",\n")
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
					processNestedObjectsForTS(builder, nestedSchema, indent+"", nestedTitle)
				}
			}
		}
	}
}
