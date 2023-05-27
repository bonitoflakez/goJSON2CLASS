package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Schema struct {
	Title      string                 `json:"title"`
	Properties map[string]interface{} `json:"properties"`
	Items      *Schema                `json:"items"`
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

func getFirstWordFromTitle(title string) string {
	titleWords := strings.Split(title, " ")
	return titleWords[0]
}
