package github_test

import (
	"os"
	"testing"

	"github.com/rin2yh/gh-ignore/internal/github"
)

// TestHelperProcess acts as a fake gh binary reusing the test binary itself.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_TEST_HELPER_PROCESS") != "1" {
		return
	}

	// Strip test framework flags that precede "gh" in os.Args.
	args := os.Args
	for i, a := range args {
		if a == "gh" {
			args = args[i:]
			break
		}
	}

	if len(args) < 3 {
		os.Exit(1)
	}

	switch args[2] {
	case "/gitignore/templates":
		os.Stdout.WriteString(`["Go","Python","Ruby"]`) //nolint:errcheck // test helper process, write errors are non-actionable
	case "/gitignore/templates/Go":
		os.Stdout.WriteString(`{"name":"Go","source":"*.exe\n*.dll\n"}`) //nolint:errcheck // test helper process, write errors are non-actionable
	case "/gitignore/templates/NotFound":
		os.Stderr.WriteString("Not Found") //nolint:errcheck // test helper process, write errors are non-actionable
		os.Exit(1)
	default:
		os.Exit(1)
	}
	os.Exit(0)
}

func setupFakeGh(t *testing.T) {
	t.Helper()
	orig := github.GhCommand
	github.GhCommand = os.Args[0]
	t.Setenv("GO_TEST_HELPER_PROCESS", "1")
	t.Cleanup(func() { github.GhCommand = orig })
}

func TestListTemplates(t *testing.T) {
	setupFakeGh(t)

	templates, err := github.ListTemplates()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(templates) != 3 {
		t.Fatalf("expected 3 templates, got %d", len(templates))
	}
	if templates[0] != "Go" || templates[1] != "Python" || templates[2] != "Ruby" {
		t.Fatalf("unexpected templates: %v", templates)
	}
}

func TestGetTemplate(t *testing.T) {
	setupFakeGh(t)

	source, err := github.GetTemplate("Go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if source != "*.exe\n*.dll\n" {
		t.Fatalf("unexpected source: %q", source)
	}
}

func TestGetTemplate_NotFound(t *testing.T) {
	setupFakeGh(t)

	_, err := github.GetTemplate("NotFound")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
