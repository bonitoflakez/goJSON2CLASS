package main

type Schema struct {
	Title      string                 `json:"title"`
	Properties map[string]interface{} `json:"properties"`
	Items      *Schema                `json:"items"`
}

type JavaType struct {
	Title      string
	Properties map[string]interface{}
	Items      interface{}
}

type CType struct {
	Title      string
	Properties map[string]interface{}
	Items      interface{}
}

type CPPType struct {
	Title      string
	Properties map[string]interface{}
	Items      interface{}
}

type GoType struct {
	Name     string
	DataType string
}

type RustType struct {
	Name     string
	DataType string
}

type TSType struct {
	Name     string
	DataType string
}
