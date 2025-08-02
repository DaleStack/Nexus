package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Constants for parsing
const (
	FileExtension = ".nx"

	// Keywords
	KeywordModule   = "module"
	KeywordState    = "state"
	KeywordAction   = "action"
	KeywordView     = "view"
	KeywordTemplate = "template"

	// Delimiters
	BlockStart = "{"
	BlockEnd   = "}"
)

// Regular expressions for validation
var (
	identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	moduleNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
)

// ParseError represents a parsing error with context
type ParseError struct {
	Line    int
	Column  int
	Message string
	Context string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s\nContext: %s",
		e.Line, e.Column, e.Message, e.Context)
}

// Parser holds the parsing state
type Parser struct {
	filename string
	lines    []string
	current  int
	line     int
	column   int
}

// NewParser creates a new parser instance
func NewParser(filename string) *Parser {
	return &Parser{
		filename: filename,
		current:  0,
		line:     1,
		column:   1,
	}
}

// ParseFile parses a Nexus file and returns a Module
func ParseFile(filename string) (*Module, error) {
	parser := NewParser(filename)
	return parser.parse()
}

// parse is the main parsing entry point
func (p *Parser) parse() (*Module, error) {
	if err := p.validateFile(); err != nil {
		return nil, err
	}

	if err := p.loadFile(); err != nil {
		return nil, fmt.Errorf("failed to load file %s: %w", p.filename, err)
	}

	if len(p.lines) == 0 {
		return nil, p.newError("file is empty", "")
	}

	module, err := p.parseModule()
	if err != nil {
		return nil, err
	}

	return module, nil
}

// validateFile checks if the file has the correct extension and exists
func (p *Parser) validateFile() error {
	if !strings.HasSuffix(p.filename, FileExtension) {
		return fmt.Errorf("invalid file extension: expected %s, got file %s",
			FileExtension, p.filename)
	}

	if _, err := os.Stat(p.filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", p.filename)
	}

	return nil
}

// loadFile reads the file content into lines
func (p *Parser) loadFile() error {
	file, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p.lines = append(p.lines, scanner.Text())
	}

	return scanner.Err()
}

// parseModule parses the entire module with block syntax
func (p *Parser) parseModule() (*Module, error) {
	// Parse module declaration with opening brace
	moduleName, err := p.parseModuleDeclaration()
	if err != nil {
		return nil, err
	}

	module := &Module{
		Name:     moduleName,
		State:    []Property{},
		Views:    []View{},
		Actions:  []Action{},
		Template: []string{},
	}

	// Parse module body until closing brace
	blockDepth := 1
	for p.hasMoreLines() && blockDepth > 0 {
		line := strings.TrimSpace(p.currentLine())

		// Check for closing brace
		if line == BlockEnd {
			blockDepth--
			if blockDepth == 0 {
				break
			}
		}

		if err := p.parseModuleElement(module); err != nil {
			return nil, err
		}
	}

	if blockDepth > 0 {
		return nil, p.newError("unclosed module block", "")
	}

	return module, nil
}

// parseModuleDeclaration parses the module declaration with opening brace
func (p *Parser) parseModuleDeclaration() (string, error) {
	if !p.hasMoreLines() {
		return "", p.newError("expected module declaration", "")
	}

	line := p.currentLine()
	trimmed := strings.TrimSpace(line)

	// Handle "module Name {" syntax
	if strings.HasPrefix(trimmed, KeywordModule+" ") && strings.HasSuffix(trimmed, " "+BlockStart) {
		// Extract module name between "module " and " {"
		content := strings.TrimSpace(trimmed[len(KeywordModule):])
		moduleName := strings.TrimSpace(strings.TrimSuffix(content, BlockStart))

		if moduleName == "" {
			return "", p.newError("empty module name", trimmed)
		}

		if !p.isValidModuleName(moduleName) {
			return "", p.newError("invalid module name: must start with letter and contain only letters, numbers, and underscores", moduleName)
		}

		p.advance()
		return moduleName, nil
	}

	return "", p.newError("expected module declaration with opening brace: 'module Name {'", trimmed)
}

// parseModuleElement parses individual module elements (state, action, view, template)
func (p *Parser) parseModuleElement(module *Module) error {
	line := strings.TrimSpace(p.currentLine())

	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "//") {
		p.advance()
		return nil
	}

	// Skip closing braces (handled by parent)
	if line == BlockEnd {
		p.advance()
		return nil
	}

	switch {
	case strings.HasPrefix(line, KeywordState+" "):
		return p.parseState(module)
	case strings.HasPrefix(line, KeywordAction+" "):
		return p.parseAction(module)
	case strings.HasPrefix(line, KeywordView+" "):
		return p.parseView(module)
	case strings.HasPrefix(line, KeywordTemplate+" "+BlockStart):
		return p.parseTemplate(module)
	default:
		// Better error message with more context
		return p.newError(fmt.Sprintf("unexpected token - expected 'state', 'action', 'view', or 'template', got: '%s'", line), line)
	}
}

// parseState parses a state declaration
func (p *Parser) parseState(module *Module) error {
	line := strings.TrimSpace(p.currentLine())
	stateLine := strings.TrimSpace(strings.TrimPrefix(line, KeywordState+" "))

	// Parse: name: type = value
	colonIndex := strings.Index(stateLine, ":")
	if colonIndex == -1 {
		return p.newError("invalid state declaration: missing colon", line)
	}

	name := strings.TrimSpace(stateLine[:colonIndex])
	if !p.isValidIdentifier(name) {
		return p.newError("invalid state variable name", name)
	}

	remainder := strings.TrimSpace(stateLine[colonIndex+1:])

	// Parse type and optional value
	var typ, value string
	if equalIndex := strings.Index(remainder, "="); equalIndex != -1 {
		typ = strings.TrimSpace(remainder[:equalIndex])
		value = strings.TrimSpace(remainder[equalIndex+1:])
		value = p.parseStringValue(value)
	} else {
		typ = remainder
	}

	if typ == "" {
		return p.newError("missing type in state declaration", line)
	}

	module.State = append(module.State, Property{
		Name:  name,
		Type:  typ,
		Value: value,
	})

	p.advance()
	return nil
}

// parseAction parses an action declaration
func (p *Parser) parseAction(module *Module) error {
	line := strings.TrimSpace(p.currentLine())
	actionLine := strings.TrimSpace(strings.TrimPrefix(line, KeywordAction+" "))

	// Parse action name and optional parameters
	parts := strings.Fields(actionLine)
	if len(parts) == 0 {
		return p.newError("missing action name", line)
	}

	actionName := parts[0]

	if !p.isValidIdentifier(actionName) {
		return p.newError("invalid action name", actionName)
	}

	// TODO: Parse parameters if needed
	// For now, we'll just store the action name
	module.Actions = append(module.Actions, Action{
		Name:       actionName,
		Parameters: []Parameter{},
	})

	p.advance()
	return nil
}

// parseView parses a view declaration
func (p *Parser) parseView(module *Module) error {
	line := strings.TrimSpace(p.currentLine())
	viewName := strings.TrimSpace(strings.TrimPrefix(line, KeywordView+" "))

	if viewName == "" {
		return p.newError("missing view name", line)
	}

	if !p.isValidIdentifier(viewName) {
		return p.newError("invalid view name", viewName)
	}

	module.Views = append(module.Views, View{
		Name:    viewName,
		Content: []string{},
	})

	p.advance()
	return nil
}

// parseTemplate parses a template block within a module
func (p *Parser) parseTemplate(module *Module) error {
	line := strings.TrimSpace(p.currentLine())
	if !strings.HasSuffix(line, BlockStart) {
		return p.newError("expected opening brace after template", line)
	}

	p.advance() // Skip template { line

	var templateLines []string
	blockDepth := 1

	for p.hasMoreLines() && blockDepth > 0 {
		line := p.currentLine()

		// Count braces to handle nested JSX elements
		openBraces := strings.Count(line, BlockStart)
		closeBraces := strings.Count(line, BlockEnd)

		// Update block depth
		blockDepth += openBraces - closeBraces

		// If we're still inside the template block, add the line
		if blockDepth > 0 {
			templateLines = append(templateLines, line)
		}

		p.advance()
	}

	if blockDepth > 0 {
		return p.newError("unclosed template block", "")
	}

	module.Template = templateLines
	return nil
}

// Helper methods

// currentLine returns the current line being parsed
func (p *Parser) currentLine() string {
	if p.current >= len(p.lines) {
		return ""
	}
	return p.lines[p.current]
}

// hasMoreLines checks if there are more lines to parse
func (p *Parser) hasMoreLines() bool {
	return p.current < len(p.lines)
}

// advance moves to the next line
func (p *Parser) advance() {
	if p.current < len(p.lines) {
		p.current++
		p.line++
		p.column = 1
	}
}

// newError creates a new ParseError with current position
func (p *Parser) newError(message, context string) *ParseError {
	return &ParseError{
		Line:    p.line,
		Column:  p.column,
		Message: message,
		Context: context,
	}
}

// isValidIdentifier checks if a string is a valid identifier
func (p *Parser) isValidIdentifier(name string) bool {
	return identifierRegex.MatchString(name)
}

// isValidModuleName checks if a string is a valid module name
func (p *Parser) isValidModuleName(name string) bool {
	return moduleNameRegex.MatchString(name)
}

// parseStringValue parses and cleans string values
func (p *Parser) parseStringValue(value string) string {
	// Remove surrounding quotes if present
	if len(value) >= 2 {
		if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
			return value[1 : len(value)-1]
		}
	}

	// Handle numeric values
	if _, err := strconv.Atoi(value); err == nil {
		return value
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return value
	}
	if value == "true" || value == "false" {
		return value
	}

	return value
}

// isWhitespace checks if a character is whitespace
func (p *Parser) isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
