package config

import (
	"fmt"
	"os"
)

type Env struct {
	GeminiAPIKey string
	GitHubToken  string
	GeminiModel  string
}

func LoadEnv() (Env, error) {
	env := Env{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GitHubToken:  os.Getenv("GITHUB_TOKEN"),
		GeminiModel:  os.Getenv("GEMINI_MODEL"),
	}

	if env.GeminiModel == "" {
		env.GeminiModel = "gemini-2.5-flash"
	}

	if env.GeminiAPIKey == "" {
		return Env{}, fmt.Errorf("GEMINI_API_KEY eksik. Örnek: export GEMINI_API_KEY='your_api_key'")
	}
	if env.GitHubToken == "" {
		return Env{}, fmt.Errorf("GITHUB_TOKEN eksik. Örnek: export GITHUB_TOKEN='ghp_xxx'")
	}

	return env, nil
}
