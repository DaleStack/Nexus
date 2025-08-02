package parser

import (
	"fmt"
	"os"
	"strings"
)

type Module struct {
	Name      string
	ViewCount int
	Actions   []string
}

// ParseFile reads a .nx file, extracts the module name, view blocks, and actions
func ParseFile(filename string) (*Module, error) {
	if !strings.HasSuffix(filename, ".nx") {
		return nil, fmt.Errorf("invalid file extension: %s", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	firstLine := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(firstLine, "module ") {
		return nil, fmt.Errorf("missing module declaration")
	}
	modName := strings.TrimPrefix(firstLine, "module ")

	views := 0
	actions := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "view ") {
			views++
		}

		if strings.HasPrefix(line, "action ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				actions = append(actions, parts[1]) // grab action name
			}
		}
	}

	return &Module{
		Name:      modName,
		ViewCount: views,
		Actions:   actions,
	}, nil
}
