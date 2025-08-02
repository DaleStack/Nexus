package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nexus <filename.nx>")
		return
	}

	filename := os.Args[1]

	if !strings.HasSuffix(filename, ".nx") {
		fmt.Printf("🚫 Invalid file extension: expected .nx, got %s\n", filename)
		return
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	source := string(content)
	fmt.Println("🔹 Nexus source code:")
	fmt.Println(source)

	lines := strings.Split(source, "\n")
	firstLine := strings.TrimSpace(lines[0])

	if strings.HasPrefix(firstLine, "module ") {
		fmt.Println("✅ Valid Nexus module syntax")
	} else {
		fmt.Println("⚠️ Missing `module` declaration at the top")
	}
}
