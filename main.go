package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	helpMsg := flag.Bool("h", false, "show usage message")
	targetLang := flag.String("l", "nil", "set a target language")
	schemaFile := flag.String("s", "schema.json", "path to file containing JSON schema")
	outputFile := flag.String("o", "output.txt", "path to output file")
	publicDef := flag.Bool("p", false, "set values to public in output code")

	flag.Parse()

	if flag.NFlag() == 0 || *helpMsg {
		usage()
		os.Exit(1)
	}

	schema, err := readJSONSchema(*schemaFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if !checkPublicSupport(*targetLang) && *publicDef {
		fmt.Println("Public is not supported for " + *targetLang)
		fmt.Println("Choosing default settings")
	}

	switch *targetLang {
	case "nil":
		fmt.Println("No language specified")
		os.Exit(1)
	case "rust":
		rustCode := generateRustCode(schema, *publicDef)
		writeRustCodeToFile(*outputFile, rustCode)
	case "c":
		cLangCode := generateCCode(schema)
		writeCCodeToFile(*outputFile, cLangCode)
	case "cpp":
		CPPCode := generateCPPCode(schema)
		writeCPPCodeToFile(*outputFile, CPPCode)
	case "go":
		goLangCode := generateGoCode(schema)
		writeGoCodeToFile(*outputFile, goLangCode)
	case "ts":
		TSCode := generateTSCode(schema)
		writeTSCodeToFile(*outputFile, TSCode)
	case "java":
		JavaCode := generateJavaCode(schema)
		writeJavaCodeToFile(*outputFile, JavaCode)
	default:
		fmt.Println(*targetLang + " is not supported :(")
	}
}
