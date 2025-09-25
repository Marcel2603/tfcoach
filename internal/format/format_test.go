package format_test

import (
	"bytes"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/format"
	"github.com/hashicorp/hcl/v2"
)

func rng(file string, line0, col int) hcl.Range {
	// line0 is zero-based (as in hcl.Pos)
	return hcl.Range{
		Filename: file,
		Start:    hcl.Pos{Line: line0, Column: col},
		End:      hcl.Pos{Line: line0, Column: col},
	}
}

func TestWriteIssues_Single(t *testing.T) {
	issues := []engine.Issue{
		{
			File:    "main.tf",
			Range:   rng("main.tf", 0, 1), // prints as line 1 due to +1
			Message: `resource name must be "this"`,
			RuleID:  "core.naming.require_this",
		},
	}

	var buf bytes.Buffer
	format.WriteIssues(issues, &buf)

	want := "main.tf:0:1: resource name must be \"this\" (core.naming.require_this)\n"
	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteIssues_Multiple(t *testing.T) {
	issues := []engine.Issue{
		{File: "a.tf", Range: rng("a.tf", 4, 7), Message: "m1", RuleID: "r1"},
		{File: "b.tf", Range: rng("b.tf", 9, 2), Message: "m2", RuleID: "r2"},
	}

	var buf bytes.Buffer
	format.WriteIssues(issues, &buf)

	want := "" +
		"a.tf:4:7: m1 (r1)\n" + // 4+1 = 5
		"b.tf:9:2: m2 (r2)\n" // 9+1 = 10
	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}
