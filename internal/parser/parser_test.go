package parser

import (
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectErr   bool
		wantName    string
		wantViews   int
		wantActions []string
	}{
		{
			name:        "valid file",
			content:     "module App\nview Home\naction login\naction logout",
			expectErr:   false,
			wantName:    "App",
			wantViews:   1,
			wantActions: []string{"login", "logout"},
		},
		{
			name:      "missing module",
			content:   "view Dashboard\naction init",
			expectErr: true,
		},
		{
			name:      "empty file",
			content:   "",
			expectErr: true,
		},
		{
			name:      "invalid extension",
			content:   "module Site",
			expectErr: true, // we're manually passing a wrong extension
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "temp_test.nx"
			if tt.name == "invalid extension" {
				filename = "bad.txt"
			}

			err := os.WriteFile(filename, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("❌ Failed to write temp file: %v", err)
			}
			defer os.Remove(filename)

			mod, err := ParseFile(filename)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if mod.Name != tt.wantName {
				t.Errorf("Module name = %s; want %s", mod.Name, tt.wantName)
			}
		})
	}
}
