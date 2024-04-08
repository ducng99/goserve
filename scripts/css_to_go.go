package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		println("Usage: css-to-go <file.css>")
		return
	}

	inputFile := args[0]
	inputAbsPath, err := filepath.Abs(inputFile)
	if err != nil {
		panic(err)
	}

	baseDir := filepath.Dir(inputAbsPath)
	outputName, _ := strings.CutSuffix(filepath.Base(inputAbsPath), filepath.Ext(inputAbsPath))
	constName := stringUptoRune(outputName, '.')
	constName = strings.ToUpper(constName[:1]) + constName[1:]

	packageName := filepath.Base(baseDir)
	outputFile := filepath.Join(baseDir, outputName+".go")

	input, err := os.ReadFile(inputAbsPath)
	if err != nil {
		panic(err)
	}

	// Escape backticks
	input = bytes.ReplaceAll(input, []byte("`"), []byte("`+\"`\"+`"))

	output, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	fmt.Fprintf(output, "package %s\n\nconst %s = `%s`\n", packageName, constName, input)
}

func stringUptoRune(s string, c rune) string {
	for i, char := range s {
		if char == c {
			return s[:i]
		}
	}

	return s
}
