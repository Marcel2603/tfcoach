package format_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/format"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
)

var (
	issues1 = []types.Issue{
		{
			File:    "main.tf",
			Range:   rng("main.tf", 0, 1), // prints as line 1 due to +1
			Message: `resource name must be "this"`,
			RuleID:  "core.naming.require_this",
		},
	}
	issues2 = []types.Issue{
		{File: "a.tf", Range: rng("a.tf", 4, 7), Message: "m1", RuleID: "r1"},
		{File: "b.tf", Range: rng("b.tf", 9, 2), Message: "m2", RuleID: "r2"},
	}
)

func rng(file string, line0, col int) hcl.Range {
	// line0 is zero-based (as in hcl.Pos)
	return hcl.Range{
		Filename: file,
		Start:    hcl.Pos{Line: line0, Column: col},
		End:      hcl.Pos{Line: line0, Column: col},
	}
}

func TestWriteResults_TextSingle(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "raw")
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `main.tf:0:1: resource name must be "this" (core.naming.require_this)
Summary:
 Issues: 1
`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteResults_TextMultiple(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues2, &buf, "raw")
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `a.tf:4:7: m1 (r1)
b.tf:9:2: m2 (r2)
Summary:
 Issues: 2
`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_JsonSingle(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "json")
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `{
  "issue_count": 1,
  "issues": [
	{
	  "file": "main.tf",
	  "line": 0,
	  "column": 1,
	  "message": "resource name must be \"this\"",
	  "rule_id": "core.naming.require_this",
	  "severity": "UNKNOWN",
	  "category": "",
	  "docs_url": "https://github.com/Marcel2603/tfcoach/tree/main/docs/pages/rules/about:blank.md"
	}
  ]
}
`

	var gotJ, wantJ interface{}

	if err = json.Unmarshal([]byte(want), &wantJ); err != nil {
		t.Fatalf("Unexpected unmarshalling error in test setup: %v, want none", err)
	}

	got := buf.Bytes()
	if err = json.Unmarshal(got, &gotJ); err != nil {
		t.Fatalf("Unexpected unmarshalling error: %v, want none", err)
	}

	if !reflect.DeepEqual(gotJ, wantJ) {
		t.Fatalf("JSON DeepEqual mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_JsonMultiple(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues2, &buf, "json")
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `{
  "issue_count": 2,
  "issues": [
	{
	  "file": "a.tf",
	  "line": 4,
	  "column": 7,
	  "message": "m1",
	  "rule_id": "r1",
	  "severity": "UNKNOWN",
	  "category": "",
	  "docs_url": "https://github.com/Marcel2603/tfcoach/tree/main/docs/pages/rules/about:blank.md"
	},
	{
	  "file": "b.tf",
	  "line": 9,
	  "column": 2,
	  "message": "m2",
	  "rule_id": "r2",
	  "severity": "UNKNOWN",
	  "category": "",
	  "docs_url": "https://github.com/Marcel2603/tfcoach/tree/main/docs/pages/rules/about:blank.md"
	}
  ]
}`
	var gotJ, wantJ interface{}

	if err = json.Unmarshal([]byte(want), &wantJ); err != nil {
		t.Fatalf("Unexpected unmarshalling error in test setup: %v, want none", err)
	}

	got := buf.Bytes()

	if err = json.Unmarshal(got, &gotJ); err != nil {
		t.Fatalf("Unexpected unmarshalling error: %v, want none", err)
	}

	if !reflect.DeepEqual(gotJ, wantJ) {
		t.Fatalf("JSON DeepEqual mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "abcd")
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	want := "unknown output format: abcd"
	if err.Error() != want {
		t.Fatalf("mismatch:\n got: %q\nwant: %q", err, want)
	}
}
