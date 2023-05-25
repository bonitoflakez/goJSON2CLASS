package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Schema struct {
	Title      string                 `json:"title"`
	Properties map[string]interface{} `json:"properties"`
	Items      *Schema                `json:"items"`
}

type RustType struct {
	Name     string
	DataType string
}

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

func main() {
	helpMsg := flag.Bool("h", false, "show usage message")
	targetLang := flag.String("l", "nil", "set a target language")
	schemaFile := flag.String("s", "schema.json", "path to file containing JSON schema")
	outputFile := flag.String("o", "output.txt", "path to output file")
	publicDef := flag.Bool("p", false, "set values to public in output code")

	flag.Parse()

	// show usage message if no flag is passed
	if flag.NFlag() == 0 {
		usage()
		os.Exit(1)
	}

	// show help message if `-h` flag is used
	if *helpMsg {
		usage()
		os.Exit(1)
	}

	schema, err := readJSONSchema(*schemaFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// for public flag testing
	if *publicDef && checkPublicSupport(*targetLang) {
		fmt.Println("public is on")
	}

	switch *targetLang {
	case "nil":
		fmt.Println("No language specified")
		os.Exit(1)
	case "rust":
		rustCode := generateRustCode(schema)
		outputRustCode(*outputFile, rustCode)
	default:
		fmt.Println(*targetLang + " is not supported :(")
	}
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

func checkPublicSupport(inp string) bool {
	supportedLanguages := map[string]bool{
		"rust": true,
		"java": true,
		"go":   true,
		"cpp":  true,
	}

	return supportedLanguages[inp]
}

/*
* TODO: Use `pub` when `-p` flag is used
 */

func outputRustCode(outFile string, rustCode string) {
	err := os.WriteFile(outFile, []byte(rustCode), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}

func generateRustCode(schema *Schema) string {
	var builder strings.Builder
	processSchemaForRust(&builder, schema, "")
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

func processSchemaForRust(builder *strings.Builder, schema *Schema, indent string) {
	if schema.Properties != nil {
		builder.WriteString("use serde::{Serialize, Deserialize};\n\n")
		builder.WriteString(indent + "#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString(indent + "struct " + getFirstWordFromTitle(schema.Title) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t#[serde(rename = \"" + name + "\")]\n")
			builder.WriteString(indent + "\t" + name + ": " + getRustType(property) + ",\n")
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
					processNestedObjectsForRust(builder, nestedSchema, indent+"", nestedTitle)
				}
			}
		}
	} else if schema.Items != nil {
		// handle array items
		builder.WriteString(indent + "#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString(indent + "struct " + getFirstWordFromTitle(schema.Title) + " {\n")
		builder.WriteString(indent + "\t" + "#[serde(rename = \"items\")]\n")
		builder.WriteString(indent + "\t" + "items: Vec<" + getRustType(schema.Items) + ">,\n")
		builder.WriteString(indent + "}\n\n")

		// handle nested objects within array items
		processNestedObjectsForRust(builder, schema.Items, indent+"", schema.Items.Title)
	}
}

func processNestedObjectsForRust(builder *strings.Builder, schema *Schema, indent string, structName string) {
	if schema.Properties != nil {
		builder.WriteString(indent + "#[derive(Debug, Serialize, Deserialize)]\n")
		builder.WriteString(indent + "struct " + getFirstWordFromTitle(structName) + " {\n")

		var propertyNames []string
		for name := range schema.Properties {
			propertyNames = append(propertyNames, name)
		}
		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := schema.Properties[name]
			builder.WriteString(indent + "\t#[serde(rename = \"" + name + "\")]\n")
			builder.WriteString(indent + "\t" + name + ": " + getRustType(property) + ",\n")
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
					processNestedObjectsForRust(builder, nestedSchema, indent+"", nestedTitle)
				}
			}
		}
	}
}

func getFirstWordFromTitle(title string) string {
	titleWords := strings.Split(title, " ")
	return titleWords[0]
}
