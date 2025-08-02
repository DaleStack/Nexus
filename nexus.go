package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nexus <filename>")
		return
	}

	filename := os.Args[1]

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	source := string(content)
	fmt.Println("🔹 Nexus source code:")
	fmt.Println(source)

	// Check if file starts with `module`
	lines := strings.Split(source, "\n")
	firstLine := strings.TrimSpace(lines[0])

	if strings.HasPrefix(firstLine, "module ") {
		fmt.Println("✅ Valid Nexus module syntax")
	} else {
		fmt.Println("⚠️ Missing `module` declaration at the top")
	}
}
