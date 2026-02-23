package config

import (
	"fmt"
	"os"
	"strings"
)

type Env struct {
	GeminiAPIKey string
	GitHubToken  string
	GeminiModel  string
}

func LoadEnv() (Env, error) {
	env := Env{
		GeminiAPIKey: sanitizeEnvValue(os.Getenv("GEMINI_API_KEY")),
		GitHubToken:  sanitizeEnvValue(os.Getenv("GITHUB_TOKEN")),
		GeminiModel:  sanitizeEnvValue(os.Getenv("GEMINI_MODEL")),
	}

	if env.GeminiModel == "" {
		env.GeminiModel = "gemini-2.5-flash"
	}

	if env.GeminiAPIKey == "" {
		return Env{}, fmt.Errorf("GEMINI_API_KEY ortam değişkeni tanımlı değil. Lütfen 'set GEMINI_API_KEY=<key>' ile ayarlayın")
	}
	if env.GitHubToken == "" {
		return Env{}, fmt.Errorf("GITHUB_TOKEN ortam değişkeni tanımlı değil. Lütfen 'set GITHUB_TOKEN=<token>' ile ayarlayın")
	}

	return env, nil
}

func sanitizeEnvValue(v string) string {
	trimmed := strings.TrimSpace(v)
	if len(trimmed) >= 2 {
		if (trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\'') ||
			(trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"') ||
			(trimmed[0] == '`' && trimmed[len(trimmed)-1] == '`') {
			return strings.TrimSpace(trimmed[1 : len(trimmed)-1])
		}
	}
	return strings.Trim(trimmed, "'\"")
}
