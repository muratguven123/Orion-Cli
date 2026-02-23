package gitrepo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Detector struct {
	workdir string
}

func NewDetector(workdir string) *Detector {
	if workdir == "" {
		workdir, _ = os.Getwd()
	}
	return &Detector{workdir: workdir}
}

func (d *Detector) DetectOwnerRepo() (string, string, error) {
	configPath := filepath.Join(d.workdir, ".git", "config")
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("bu klasörde git reposu bulunamadı (.git/config yok)")
		}
		return "", "", fmt.Errorf(".git/config okunamadı: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	insideOrigin := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") {
			insideOrigin = strings.EqualFold(line, "[remote \"origin\"]")
			continue
		}
		if !insideOrigin || !strings.HasPrefix(strings.ToLower(line), "url") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		owner, repo, err := parseRemoteURL(strings.TrimSpace(parts[1]))
		if err != nil {
			return "", "", err
		}
		return owner, repo, nil
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf(".git/config taranamadı: %w", err)
	}

	return "", "", fmt.Errorf("origin remote bulunamadı")
}

var patterns = []*regexp.Regexp{
	regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+?)(?:\.git)?$`),
	regexp.MustCompile(`^git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`),
	regexp.MustCompile(`^ssh://git@github\.com/([^/]+)/([^/]+?)(?:\.git)?$`),
}

func parseRemoteURL(remote string) (string, string, error) {
	for _, p := range patterns {
		match := p.FindStringSubmatch(remote)
		if len(match) == 3 {
			return match[1], match[2], nil
		}
	}
	return "", "", fmt.Errorf("origin URL GitHub formatında değil: %s", remote)
}
