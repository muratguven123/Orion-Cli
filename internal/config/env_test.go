package config

import "testing"

func TestSanitizeEnvValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain", input: "AIzaSy123", want: "AIzaSy123"},
		{name: "single quoted", input: "'AIzaSy123'", want: "AIzaSy123"},
		{name: "double quoted", input: "\"ghp_token\"", want: "ghp_token"},
		{name: "backtick quoted", input: "`gemini-2.5-flash`", want: "gemini-2.5-flash"},
		{name: "dangling single quote", input: "AIzaSy123'", want: "AIzaSy123"},
		{name: "spaces and quote", input: "  ' ghp_token '  ", want: "ghp_token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeEnvValue(tt.input)
			if got != tt.want {
				t.Fatalf("sanitizeEnvValue(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
