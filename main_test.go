package main

import (
	"os"
	"regexp"
	"testing"
)

var tables []Table
var generatedStructs []byte

const structCount int = 12

func TestSurqlParsing(t *testing.T) {
	// Read the .surql file
	surqlData, err := os.ReadFile("testdata/test_schema.surql")
	if err != nil {
		t.Fatalf("Failed to read .surql file: %v", err)
	}

	// Parse the .surql file into a slice of tables
	tables = parseSurql(string(surqlData))

	// Check that the number of tables is correct
	if len(tables) != structCount {
		t.Fatalf("Expected %d generated structs, got %d", structCount, len(tables))
	}
}

func TestGenerateStructs(t *testing.T) {
	// Generate the Go structs
	structBuilder := buildStructStringsFromTables(tables)

	if !writeStructsToFile(structBuilder, true) {
		t.Fatal("Failed to write structs to structs file")
		os.Exit(1)
	}

	fileContent, err := os.ReadFile("testdata/generated_structs")
	if err != nil {
		t.Fatalf("Failed to read expected generated structs: %v", err)
	}

	if len(fileContent) == 0 {
		t.Fatal("Expected generated structs to be non-empty")
	}

	if string(fileContent) != structBuilder.String() {
		t.Fatal("Generated structs do not match")
	}
}

func TestGeneratedStructsCount(t *testing.T) {
	// Read the generated structs from a file
	generatedStructs, err := os.ReadFile("testdata/generated_structs")
	if err != nil {
		t.Fatalf("Failed to read generated structs: %v", err)
	}

	re := regexp.MustCompile(`(?m)type \w+ struct \{`)
	generatedStructCount := len(re.FindAllIndex(generatedStructs, -1))

	// Check that the number of generated structs is correct
	if generatedStructCount != structCount {
		t.Fatalf("Expected %d generated structs, got %d", structCount, generatedStructCount)
	}
}
