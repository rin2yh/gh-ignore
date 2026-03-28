package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// GhCommand is the gh CLI command name, overridable in tests
var GhCommand = "gh"

func runGhAPI(path string) ([]byte, error) {
	out, err := exec.Command(GhCommand, "api", path).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("gh api failed: %s", exitErr.Stderr)
		}
		return nil, fmt.Errorf("gh api failed: %w", err)
	}
	return out, nil
}

func ListTemplates() ([]string, error) {
	out, err := runGhAPI("/gitignore/templates")
	if err != nil {
		return nil, err
	}

	var templates []string
	if err := json.Unmarshal(out, &templates); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return templates, nil
}

func GetTemplate(name string) (string, error) {
	out, err := runGhAPI("/gitignore/templates/" + name)
	if err != nil {
		return "", err
	}

	var result struct {
		Source string `json:"source"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return result.Source, nil
}
