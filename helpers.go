package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// general functions

func usage() {
	fmt.Println("Usage: goJSON2CLASS -l <target-lang> -s <schema.json> -o <output.ext>")
	fmt.Println()
	fmt.Println("\t-l >> choose a language.")
	fmt.Println("\t\tExample: `-l rust` (default: nil)")
	fmt.Println()
	fmt.Println("\t-s >> path to file containing JSON schema. (default: schema.json)")
	fmt.Println("\t\tExample: `-s schema.json`")
	fmt.Println()
	fmt.Println("\t-o >> path to output file with extension. (default: output.txt)")
	fmt.Println("\t\tExample: `-o output.rs`")
	fmt.Println()
	fmt.Println("\t-p >> define public if supported by language (default: false)")
	fmt.Println("\t\tExample: `-p`")
}

func readJSONSchema(filePath string) (*Schema, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
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

func writeCodeToFile(outFile string, generateCode string) {
	err := os.WriteFile(outFile, []byte(generateCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func checkPublicSupport(inp string) bool {
	supportedLanguages := map[string]bool{
		"rust": true,
		// TODO: Implement public definitions of these langs
		// "java": true,
		// "go":   true,
		// "cpp":  true,
	}

	return supportedLanguages[inp]
}

func getFirstWordFromTitle(title string) string {
	titleWords := strings.Split(title, " ")
	return titleWords[0]
}

// functions for C handler

func cHeaderFormat() string {
	return getCHeaderIncludes() + "\n" +
		getPreprocessorDirectives() + "\n" +
		getTypedefStructsList() + "\n"
}

func getPreprocessorDirectives() string {
	var builder strings.Builder

	definesMap := getPreprocessorSizeDefinesMap()
	for define, value := range definesMap {
		builder.WriteString(fmt.Sprintf("#define %s %d\n", define, value))
	}

	return builder.String()
}

func getTypedefStructsList() string {
	var typedefStructBuilder strings.Builder
	for _, structName := range typedefStructsList {
		typedefStructBuilder.WriteString("typedef struct " + structName + " " + structName + ";\n")
	}
	return typedefStructBuilder.String()
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
			return getCDataType(items)
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

func getCHeaderIncludes() string {
	return `#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
`
}

// functions for CPP handler

func getCPPHeaderIncludes() string {
	return `#include <iostream>
#include <vector>
#include <string>`
}

func isCPPArrayType(property interface{}) bool {
	if p, ok := property.(map[string]interface{}); ok {
		if _, ok := p["type"]; ok {
			return p["type"] == "array"
		}
	}
	return false
}

func isCPPObjectType(property interface{}) bool {
	if p, ok := property.(map[string]interface{}); ok {
		if _, ok := p["type"]; ok {
			return p["type"] == "object"
		}
	}
	return false
}

func getCPPArrayType(property interface{}) string {
	if p, ok := property.(map[string]interface{}); ok {
		if items, ok := p["items"]; ok {
			return getItemCPPType(items)
		}
	}
	return "unknown"
}

// used in c, cpp and java handler
func addToTypedefStructsList(structName string) {
	typedefStructsList = append(typedefStructsList, structName)
}

// function for rust handler

func getPropertyDeclaration(name, typ string, pubFlag bool) string {
	if pubFlag {
		return "pub " + name + ": " + typ
	}
	return name + ": " + typ
}

// functions for java handler

func isJavaArrayType(property interface{}) bool {
	switch p := property.(type) {
	case map[string]interface{}:
		if _, ok := p["type"]; ok {
			return p["type"] == "array"
		}
	}
	return false
}

func isJavaObjectType(property interface{}) bool {
	switch p := property.(type) {
	case map[string]interface{}:
		if _, ok := p["type"]; ok {
			return p["type"] == "object"
		}
	}
	return false
}

func getJavaArrayType(property interface{}) string {
	switch p := property.(type) {
	case map[string]interface{}:
		if items, ok := p["items"]; ok {
			return getItemJavaType(items)
		}
	}
	return "unknown"
}
