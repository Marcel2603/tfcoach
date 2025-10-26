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
			Message: `Block "a" should be inside of "b.tf"`,
			RuleID:  "core.file_naming",
		},
	}
	issues2 = []types.Issue{
		{File: "a.tf", Range: rng("a.tf", 4, 7), Message: "m1", RuleID: "core.something_something"},
		{File: "b.tf", Range: rng("b.tf", 9, 2), Message: "m2", RuleID: "core.naming_convention"},
	}
	issues3 = []types.Issue{
		{File: "a.tf", Range: rng("a.tf", 4, 7), Message: "m1", RuleID: "core.something_something"},
		{File: "b.tf", Range: rng("b.tf", 9, 2), Message: "m2", RuleID: "core.naming_convention"},
		{File: "a.tf", Range: rng("a.tf", 10, 2), Message: "m3", RuleID: "core.naming_convention"},
		{File: "a.tf", Range: rng("a.tf", 2, 1), Message: "m4", RuleID: "core.file_naming"},
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

func TestWriteResults_CompactSingle(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "compact", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `L main.tf:0:1: Block "a" should be inside of "b.tf" [core.file_naming]
Summary: 1 issue
`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteResults_CompactMultiple(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues2, &buf, "compact", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `H b.tf:9:2: m2 [core.naming_convention]
U a.tf:4:7: m1 [core.something_something]
Summary: 2 issues
`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_JsonSingle(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "json", true)
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
	  "message": "Block \"a\" should be inside of \"b.tf\"",
	  "rule_id": "core.file_naming",
	  "severity": "LOW",
	  "category": "",
	  "docs_url": "https://marcel2603.github.io/tfcoach/rules/core/file_naming"
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
	err := format.WriteResults(issues2, &buf, "json", true)
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
	  "rule_id": "core.something_something",
	  "severity": "UNKNOWN",
	  "category": "",
	  "docs_url": "about:blank"
	},
	{
	  "file": "b.tf",
	  "line": 9,
	  "column": 2,
	  "message": "m2",
	  "rule_id": "core.naming_convention",
	  "severity": "HIGH",
	  "category": "",
	  "docs_url": "https://marcel2603.github.io/tfcoach/rules/core/naming_convention"
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

func TestWriteResults_PrettySingle(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "pretty", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 1 issue found in 1 file

─── main.tf ────────────

  0:1	[core.file_naming]	LOW
	💡  Block "a" should be inside of "b.tf"
	📑  https://marcel2603.github.io/tfcoach/rules/core/file_naming

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteResults_PrettyMultiple(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues2, &buf, "pretty", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 2 issues found in 2 files

─── a.tf ───────────────

  4:7	[core.something_something]	UNKNOWN
	💡  m1
	📑  about:blank

─── b.tf ───────────────

  9:2	[core.naming_convention]	HIGH
	💡  m2
	📑  https://marcel2603.github.io/tfcoach/rules/core/naming_convention

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_PrettySorting(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues3, &buf, "pretty", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 4 issues found in 2 files

─── a.tf ───────────────

  10:2	[core.naming_convention]	HIGH
	💡  m3
	📑  https://marcel2603.github.io/tfcoach/rules/core/naming_convention

  2:1	[core.file_naming]	LOW
	💡  m4
	📑  https://marcel2603.github.io/tfcoach/rules/core/file_naming

  4:7	[core.something_something]	UNKNOWN
	💡  m1
	📑  about:blank

─── b.tf ───────────────

  9:2	[core.naming_convention]	HIGH
	💡  m2
	📑  https://marcel2603.github.io/tfcoach/rules/core/naming_convention

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_PrettyNoEmojis(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues3, &buf, "pretty", false)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 4 issues found in 2 files

─── a.tf ───────────────

  10:2	[core.naming_convention]	HIGH
	m3
	Docs: https://marcel2603.github.io/tfcoach/rules/core/naming_convention

  2:1	[core.file_naming]	LOW
	m4
	Docs: https://marcel2603.github.io/tfcoach/rules/core/file_naming

  4:7	[core.something_something]	UNKNOWN
	m1
	Docs: about:blank

─── b.tf ───────────────

  9:2	[core.naming_convention]	HIGH
	m2
	Docs: https://marcel2603.github.io/tfcoach/rules/core/naming_convention

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_EducationalSingle(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "educational", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 1 rule broken (1 issue total)

─── File Naming (Severity LOW) ─────────

💡  File naming should follow a strict convention.

🆔  [core.file_naming]
📑  https://marcel2603.github.io/tfcoach/rules/core/file_naming

⚠️  Broken at:
🔹 main.tf:0:1 ➡️  Block "a" should be inside of "b.tf"

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteResults_EducationalMultiple(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues2, &buf, "educational", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 2 rules broken (2 issues total)

─── Naming Convention (Severity HIGH) ─────────

💡  Terraform names should only contain lowercase alphanumeric characters and underscores.

🆔  [core.naming_convention]
📑  https://marcel2603.github.io/tfcoach/rules/core/naming_convention

⚠️  Broken at:
🔹 b.tf:9:2 ➡️  m2


─── Unknown (Severity UNKNOWN) ────────────────

💡  Unknown rule

🆔  [core.something_something]
📑  about:blank

⚠️  Broken at:
🔹 a.tf:4:7 ➡️  m1

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_EducationalSorting(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues3, &buf, "educational", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 3 rules broken (4 issues total)

─── Naming Convention (Severity HIGH) ─────────

💡  Terraform names should only contain lowercase alphanumeric characters and underscores.

🆔  [core.naming_convention]
📑  https://marcel2603.github.io/tfcoach/rules/core/naming_convention

⚠️  Broken at:
🔹 a.tf:10:2 ➡️  m3
🔹 b.tf:9:2 ➡️  m2


─── File Naming (Severity LOW) ────────────────

💡  File naming should follow a strict convention.

🆔  [core.file_naming]
📑  https://marcel2603.github.io/tfcoach/rules/core/file_naming

⚠️  Broken at:
🔹 a.tf:2:1 ➡️  m4


─── Unknown (Severity UNKNOWN) ────────────────

💡  Unknown rule

🆔  [core.something_something]
📑  about:blank

⚠️  Broken at:
🔹 a.tf:4:7 ➡️  m1

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_EducationalNoEmojis(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues3, &buf, "educational", false)
	if err != nil {
		t.Fatalf("Unexpected error: %v, want none", err)
	}

	want := `Summary: 3 rules broken (4 issues total)

─── Naming Convention (Severity HIGH) ─────────

Explanation: Terraform names should only contain lowercase alphanumeric characters and underscores.

ID: [core.naming_convention]
Read more: https://marcel2603.github.io/tfcoach/rules/core/naming_convention

Broken at:
- a.tf:10:2 ─ m3
- b.tf:9:2 ─ m2


─── File Naming (Severity LOW) ────────────────

Explanation: File naming should follow a strict convention.

ID: [core.file_naming]
Read more: https://marcel2603.github.io/tfcoach/rules/core/file_naming

Broken at:
- a.tf:2:1 ─ m4


─── Unknown (Severity UNKNOWN) ────────────────

Explanation: Unknown rule

ID: [core.something_something]
Read more: about:blank

Broken at:
- a.tf:4:7 ─ m1

`

	if got := buf.String(); got != want {
		t.Fatalf("mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteResults_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteResults(issues1, &buf, "abcd", true)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	want := "unknown output format: abcd"
	if err.Error() != want {
		t.Fatalf("mismatch:\n got: %q\nwant: %q", err, want)
	}
}
