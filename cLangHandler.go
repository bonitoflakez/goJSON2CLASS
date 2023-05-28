package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var preprocessorSizeDefinesMap = make(map[string]int)
var typedefStructsList []string

type CType struct {
	Title      string
	Properties map[string]interface{}
	Items      interface{}
}

func writeCCodeToFile(outFile string, CCode string) {
	var generatedCCode string = getHeaderIncludes() + "\n" +
		getPreprocessorDirectives() + "\n" +
		getTypedefStructsList() + "\n" +
		CCode

	err := os.WriteFile(outFile, []byte(generatedCCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

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
				title, ok := p["title"].(string)
				if ok {
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
				title, ok := p["title"].(string)
				if ok {
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

		var structName string = getFirstWordFromTitle(schema.Title)
		builder.WriteString(indent + "struct " + structName + " {\n")

		for _, name := range propertyNames {
			property := schema.Properties[name]
			if isArrayType(property) {
				itemType := getArrayType(property)
				var hashDefineMacro string = addToDefinesMap(structName, name, 50)
				builder.WriteString(indent + "    " + itemType + " " + name + "[" + hashDefineMacro + "]" + ";\n")
			} else {
				propertyType := getCType(property)
				builder.WriteString(indent + "    " + propertyType + " " + name + ";\n")
			}
		}

		builder.WriteString(indent + "};\n")
	}
}

func isArrayType(property interface{}) bool {
	switch p := property.(type) {
	case map[string]interface{}:
		if _, ok := p["type"]; ok {
			return p["type"] == "array"
		}
	}
	return false
}

func isObjectType(property interface{}) bool {
	switch p := property.(type) {
	case map[string]interface{}:
		if _, ok := p["type"]; ok {
			return p["type"] == "object"
		}
	}
	return false
}

func getArrayType(property interface{}) string {
	switch p := property.(type) {
	case map[string]interface{}:
		if items, ok := p["items"]; ok {
			return getItemCType(items)
		}
	}
	return "unknown"
}

func addToDefinesMap(structName string, propertyName string, value int) string {
	hashDefineMacro := fmt.Sprintf("%s_%s_SIZE", strings.ToUpper(structName), strings.ToUpper(propertyName))
	preprocessorSizeDefinesMap[hashDefineMacro] = value
	return hashDefineMacro
}

func getPreprocessorSizeDefinesMap() map[string]int {
	return preprocessorSizeDefinesMap
}

func getHeaderIncludes() string {
	return `#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
`
}

func getPreprocessorDirectives() string {
	var builder strings.Builder

	definesMap := getPreprocessorSizeDefinesMap()
	for define, value := range definesMap {
		builder.WriteString(fmt.Sprintf("#define %s %d\n", define, value))
	}

	return builder.String()
}

func addToTypedefStructsList(structName string) {
	typedefStructsList = append(typedefStructsList, structName)
}

func getTypedefStructsList() string {
	var typedefStructBuilder strings.Builder
	for _, structName := range typedefStructsList {
		typedefStructBuilder.WriteString("typedef struct " + structName + " " + structName + ";\n")
	}
	return typedefStructBuilder.String()
}
