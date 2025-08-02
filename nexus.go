package main

import (
	"fmt"
	"os"
	"strings"

	"nexus/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nexus <filename.nx>")
		return
	}

	filename := os.Args[1]

	// Optional: centralize this check in parser for reusability
	if !strings.HasSuffix(filename, ".nx") {
		fmt.Printf("🚫 Invalid file extension: %s\n", filename)
		return
	}

	module, err := parser.ParseFile(filename)
	if err != nil {
		fmt.Printf("❌ Error parsing file: %v\n", err)
		return
	}

	// 💡 Summary Output
	fmt.Println("🔍 Nexus Module Summary")
	fmt.Printf("📦 Name: %s\n", module.Name)
	fmt.Printf("👁️ View Blocks: %d\n", module.ViewCount)
	fmt.Printf("🎯 Actions: %v\n", module.Actions)
}
