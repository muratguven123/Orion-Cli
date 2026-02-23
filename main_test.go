package main

import (
	"strings"
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	msg := "gemini failed with key AIzaSyBDYDT9xPJkbUYg27BEd5tC01dOXkP7aCQ and token ghp_abcdefghijklmnopqrstuvwxyzABCD"
	got := redactSecrets(msg)

	if got == msg {
		t.Fatalf("expected secrets to be redacted")
	}
	if strings.Contains(got, "AIzaSyBDYDT9xPJkbUYg27BEd5tC01dOXkP7aCQ") || strings.Contains(got, "ghp_abcdefghijklmnopqrstuvwxyzABCD") {
		t.Fatalf("secret still present in output: %s", got)
	}
}
