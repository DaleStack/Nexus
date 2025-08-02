package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nexus/internal/parser"
)

const (
	// Application constants
	AppName       = "Nexus"
	Version       = "1.0.0"
	FileExtension = ".nx"

	// Exit codes
	ExitSuccess      = 0
	ExitInvalidArgs  = 1
	ExitInvalidFile  = 2
	ExitParsingError = 3
)

// Config holds application configuration
type Config struct {
	Filename string
	Verbose  bool
}

// Application represents the main application state
type Application struct {
	config Config
}

func main() {
	app := &Application{}

	if err := app.parseArgs(); err != nil {
		app.printError("Argument error", err)
		app.printUsage()
		os.Exit(ExitInvalidArgs)
	}

	if err := app.run(); err != nil {
		app.printError("Runtime error", err)
		os.Exit(ExitParsingError)
	}
}

// parseArgs parses and validates command line arguments
func (app *Application) parseArgs() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("missing required filename argument")
	}

	filename := os.Args[1]

	// Validate file extension
	if !strings.HasSuffix(filename, FileExtension) {
		return fmt.Errorf("invalid file extension: expected %s, got %s",
			FileExtension, filepath.Ext(filename))
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	app.config.Filename = filename
	return nil
}

// run executes the main application logic
func (app *Application) run() error {
	module, err := parser.ParseFile(app.config.Filename)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", app.config.Filename, err)
	}

	app.displayModuleSummary(module)
	return nil
}

// displayModuleSummary outputs a formatted summary of the parsed module
func (app *Application) displayModuleSummary(module *parser.Module) {
	fmt.Printf("\n%s Module Analysis\n", AppName)
	fmt.Println(strings.Repeat("=", 50))

	// Module name
	fmt.Printf("Module Name: %s\n", module.Name)
	fmt.Printf("Source File: %s\n\n", app.config.Filename)

	// State variables section
	app.displayStateVariables(module.State)

	// Template section
	app.displayTemplate(module.Template)

	// Actions section
	app.displayActions(module.Actions)

	fmt.Println(strings.Repeat("=", 50))
}

// displayStateVariables formats and displays state variables
func (app *Application) displayStateVariables(state []parser.StateProperty) {
	fmt.Println("State Variables:")
	if len(state) == 0 {
		fmt.Println("  No state variables defined")
		fmt.Println()
		return
	}

	for i, prop := range state {
		fmt.Printf("  [%d] %s: %s", i+1, prop.Name, prop.Type)
		if prop.Value != "" {
			fmt.Printf(" = %q", prop.Value)
		}
		fmt.Println()
	}
	fmt.Println()
}

// displayTemplate formats and displays template content
func (app *Application) displayTemplate(template []string) {
	fmt.Println("Template:")
	if len(template) == 0 {
		fmt.Println("  No template defined")
		fmt.Println()
		return
	}

	fmt.Printf("  %d line(s) of template code:\n", len(template))
	for i, line := range template {
		fmt.Printf("  %3d: %s\n", i+1, line)
	}
	fmt.Println()
}

// displayActions formats and displays available actions
func (app *Application) displayActions(actions []parser.Action) {
	fmt.Println("Actions:")
	if len(actions) == 0 {
		fmt.Println("  No actions defined")
		fmt.Println()
		return
	}

	fmt.Printf("  %d action(s) available:\n", len(actions))
	for i, action := range actions {
		fmt.Printf("  [%d] %s\n", i+1, action.Name)
	}
	fmt.Println()
}

// printError outputs formatted error messages
func (app *Application) printError(context string, err error) {
	fmt.Fprintf(os.Stderr, "Error: %s - %v\n", context, err)
}

// printUsage displays usage information
func (app *Application) printUsage() {
	fmt.Printf("Usage: %s <filename%s>\n", strings.ToLower(AppName), FileExtension)
	fmt.Printf("Version: %s\n", Version)
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s app.nx\n", strings.ToLower(AppName))
	fmt.Printf("  %s components/header.nx\n", strings.ToLower(AppName))
}
