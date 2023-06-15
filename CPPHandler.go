package main

import (
	"sort"
	"strings"
)

var CPPtypedefStructsList []string

func generateCPPCode(schema *Schema) string {
	var builder strings.Builder

	processSchemaForCPP(&builder, schema, "")
	return builder.String()
}

func getCPPType(property interface{}) string {
	switch p := property.(type) {
	case map[string]interface{}:
		if pType, ok := p["type"].(string); ok {
			switch pType {
			case "string":
				return "std::string"
			case "number":
				return "double"
			case "integer":
				return "int"
			case "boolean":
				return "bool"
			case "decimal":
				return "float"
			case "object":
				if title, ok := p["title"].(string); ok {
					return getFirstWordFromTitle(title)
				}
			}
		}
	}
	return "unknown"
}

func getItemCPPType(property interface{}) string {
	switch p := property.(type) {
	case map[string]interface{}:
		if pType, ok := p["type"].(string); ok {
			switch pType {
			case "string":
				return "std::string"
			case "number":
				return "double"
			case "integer":
				return "int"
			case "decimal":
				return "float"
			case "object":
				if title, ok := p["title"].(string); ok {
					return getFirstWordFromTitle(title)
				}
			}
		}
	}
	return "unknown"
}

func processSchemaForCPP(builder *strings.Builder, schema *Schema, indent string) {
	if schema.Properties != nil {
		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		if schema.Title != "" {
			firstStructName := getFirstWordFromTitle(schema.Title)
			addToTypedefStructsList(firstStructName)
		}

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
					processSchemaForCPP(builder, nestedSchema, indent)
				} else if isCPPArrayType(property) {
					propertyMap := property.(map[string]interface{})
					nestedSchema := propertyMap["items"].(map[string]interface{})
					if isCPPObjectType(nestedSchema) {
						nestedTitle := nestedSchema["title"].(string)
						nestedSchema = nestedSchema["properties"].(map[string]interface{})
						nestedPropertyMap := nestedSchema
						itemsSchema := &Schema{
							Title:      nestedTitle,
							Properties: nestedPropertyMap,
						}
						processSchemaForCPP(builder, itemsSchema, indent)
					}
				}
			}
		}

		structName := getFirstWordFromTitle(schema.Title)
		builder.WriteString(indent + "struct " + structName + " {\n")

		for _, name := range propertyNames {
			property := schema.Properties[name]
			if isCPPArrayType(property) {
				itemType := getCPPArrayType(property)
				builder.WriteString(indent + "    std::vector<" + itemType + "> " + name + ";\n")
			} else {
				propertyType := getCPPType(property)
				builder.WriteString(indent + "    " + propertyType + " " + name + ";\n")
			}
		}

		builder.WriteString(indent + "};\n\n")
	}
}
