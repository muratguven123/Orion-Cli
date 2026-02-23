package main

import (
	"fmt"
	"os"
	"regexp"

	"orion-cli/cmd"
)

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`AIza[0-9A-Za-z_\-]{20,}`),
	regexp.MustCompile(`gh[pousr]_[0-9A-Za-z]{20,}`),
}

func redactSecrets(msg string) string {
	redacted := msg
	for _, p := range secretPatterns {
		redacted = p.ReplaceAllString(redacted, "[REDACTED]")
	}
	return redacted
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, redactSecrets(err.Error()))
		os.Exit(1)
	}
}
