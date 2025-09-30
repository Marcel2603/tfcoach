package runner_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/runner"
	"github.com/Marcel2603/tfcoach/internal/testutil"
)

func TestRunLint_NoIssues(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{"ok.tf": `# nothing`}}
	var rules []engine.Rule // no rules -> no issues
	var out bytes.Buffer
	code := runner.Lint(".", src, rules, &out)
	if code != 0 {
		t.Fatalf("want 0, got %d", code)
	}
	if out.Len() != 0 {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestRunLint_Issues(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{"bad.tf": `resource "x" "y" {}`}}
	rules := []engine.Rule{testutil.AlwaysFlag{
		RuleID: "test.always.flag", Message: "failed", Match: "", // always emits
	}}
	var out bytes.Buffer
	code := runner.Lint(".", src, rules, &out)
	if code != 2 {
		t.Fatalf("want 2, got %d", code)
	}
	if !strings.Contains(out.String(), "failed") {
		t.Fatalf("missing message")
	}
	if !strings.Contains(out.String(), "test.always.flag") {
		t.Fatalf("missing id")
	}
}
