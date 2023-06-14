package main

import (
	"sort"
	"strings"
)

var preprocessorSizeDefinesMap = make(map[string]int)
var typedefStructsList []string

func generateCCode(schema *Schema) string {
	var builder strings.Builder
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
				if title, ok := p["title"].(string); ok {
					return getFirstWordFromTitle(title)
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
				return "char*"
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

func processSchemaForC(builder *strings.Builder, schema *Schema, indent string) {
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
					processSchemaForC(builder, nestedSchema, indent)
				} else if isArrayType(property) {
					propertyMap := property.(map[string]interface{})
					nestedSchema := propertyMap["items"].(map[string]interface{})
					if isObjectType(nestedSchema) {
						nestedTitle := nestedSchema["title"].(string)
						nestedSchema = nestedSchema["properties"].(map[string]interface{})
						nestedPropertyMap := nestedSchema
						itemsSchema := &Schema{
							Title:      nestedTitle,
							Properties: nestedPropertyMap,
						}
						processSchemaForC(builder, itemsSchema, indent)
					}
				}
			}
		}

		structName := getFirstWordFromTitle(schema.Title)
		builder.WriteString(indent + "struct " + structName + " {\n")

		for _, name := range propertyNames {
			property := schema.Properties[name]
			if isArrayType(property) {
				itemType := getArrayType(property)
				hashDefineMacro := addToDefinesMap(structName, name, 50)
				builder.WriteString(indent + "    " + itemType + " " + name + "[" + hashDefineMacro + "]" + ";\n")
			} else {
				propertyType := getCType(property)
				builder.WriteString(indent + "    " + propertyType + " " + name + ";\n")
			}
		}

		builder.WriteString(indent + "};\n")
	}
}
