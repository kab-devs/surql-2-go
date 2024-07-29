package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var surqlFile string

type Field struct {
	Name string
	Type string
}

type Table struct {
	Name   string
	Fields []Field
}

func main() {
	flag.StringVar(&surqlFile, "filename", "", "Which surgl file should I use to generate the Go struct code?")
	flag.Parse()

	if surqlFile == "" {
		fmt.Println("No .surql file specified.")
		return
	}

	surqlData, err := os.ReadFile(surqlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tables := parseSurql(string(surqlData))

	structBuilder := buildStructStringsFromTables(tables)

	if !writeStructsToFile(structBuilder, false) {
		fmt.Println("Failed to write structs to structs file")
		os.Exit(1)
	}

	fmt.Println("Wrote structs to structs file successfully")
}

func parseSurql(schema string) []Table {
	tableRegex := regexp.MustCompile(`DEFINE TABLE (\w+)`)
	fieldRegex := regexp.MustCompile(`DEFINE FIELD (\w+)\s+ON(?:\s+TABLE)?\s+(\w+)\s+TYPE (\w+)`)

	tableMatches := tableRegex.FindAllStringSubmatch(schema, -1)
	fieldMatches := fieldRegex.FindAllStringSubmatch(schema, -1)

	tables := make(map[string]*Table)
	for _, match := range tableMatches {
		tableName := match[1]
		tables[tableName] = &Table{Name: tableName}
	}

	for _, match := range fieldMatches {
		fieldName := match[1]
		tableName := match[2]
		fieldType := match[3]
		field := Field{Name: fieldName, Type: mapSurqlTypeToGoType(fieldType)}

		if table, ok := tables[tableName]; ok {
			table.Fields = append(table.Fields, field)
		}
	}

	result := []Table{}
	for _, table := range tables {
		result = append(result, *table)
	}

	return result
}

func buildStructStringsFromTables(tables []Table) strings.Builder {
	var structBuilder strings.Builder
	for _, table := range tables {
		structBuilder.WriteString(generateGoStruct(table))
	}
	return structBuilder
}

func writeStructsToFile(structBuilder strings.Builder, isTest bool) bool {
	path := "generated_structs"
	if isTest {
		path = "testdata/generated_structs"
	}

	if err := os.WriteFile(path, []byte(structBuilder.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	return true
}

var surqlTypeToGoType = map[string]string{
	"string":  "string",
	"int":     "int",
	"float":   "float64",
	"bool":    "bool",
	"default": "interface{}",
}

func mapSurqlTypeToGoType(surqlType string) string {
	return surqlTypeToGoType[surqlType]
}

func generateGoStruct(table Table) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "type %s struct {\n", strings.Title(table.Name))

	for _, field := range table.Fields {
		fmt.Fprintf(&sb, "\t%s %s `json:\"%s\"`\n", strings.Title(field.Name), field.Type, field.Name)
	}

	sb.WriteString("}\n")
	return sb.String()
}
