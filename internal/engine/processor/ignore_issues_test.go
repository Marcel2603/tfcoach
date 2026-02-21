package processor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/engine/processor"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/hashicorp/hcl/v2"
)

func TestIgnoreIssuesProcessor_ShouldRespectTfcoachnoreport(t *testing.T) {
	tempDir := t.TempDir()
	createFile(t, filepath.Join(tempDir, ".tfcoachnoreport"), `
a.tf
`)

	ignored := `
resource "test" "a"{}
resource "test" "b"{}
`
	resource2 := `
resource "test" "c"{}
`
	ignoredFile := testutil.ParseToHcl(t, "a.tf", ignored)
	ignoreIssueProcessor, err := processor.NewIgnoreIssuesProcessor(tempDir)
	if err != nil {
		t.Fatal("Setup error: ", err)
	}
	ignoreIssueProcessor.ScanFile([]byte(ignored), ignoredFile, "a.tf")

	anotherFile := testutil.ParseToHcl(t, "b.tf", resource2)
	ignoreIssueProcessor.ScanFile([]byte(resource2), anotherFile, "b.tf")

	issues := []types.Issue{
		{File: "a.tf", RuleID: "rule-a", Range: hcl.Range{Start: hcl.Pos{Line: 2}}},
		{File: "a.tf", RuleID: "another-rule", Range: hcl.Range{Start: hcl.Pos{Line: 4}}},
		{File: "b.tf", RuleID: "rule-a", Range: hcl.Range{Start: hcl.Pos{Line: 4}}},
	}

	processedIssues := ignoreIssueProcessor.ProcessIssues(issues)

	if len(processedIssues) != 1 {
		t.Fatalf("Wrong number of expected issues; got %d, wanted %d", len(processedIssues), 1)
	}
}

func TestIgnoreIssuesProcessor_ProcessFileIgnore(t *testing.T) {
	ignored := `
# tfcoach-ignore-file:rule-a
resource "test" "ignored"{
}
resource "test" "non_compliant"{}
`
	resource2 := `
resource "test" "non_compliant"{}
`
	ignoredFile := testutil.ParseToHcl(t, "main.tf", ignored)
	ignoreIssueProcessor, err := processor.NewIgnoreIssuesProcessor(".")
	if err != nil {
		t.Fatal("Setup error: ", err)
	}
	ignoreIssueProcessor.ScanFile([]byte(ignored), ignoredFile, "main.tf")

	anotherFile := testutil.ParseToHcl(t, "another.tf", resource2)
	ignoreIssueProcessor.ScanFile([]byte(resource2), anotherFile, "another.tf")
	issues := []types.Issue{
		{File: "main.tf", RuleID: "rule-a", Range: hcl.Range{Start: hcl.Pos{Line: 2}}},
		{File: "main.tf", RuleID: "another-rule", Range: hcl.Range{Start: hcl.Pos{Line: 4}}},
		{File: "another.tf", RuleID: "rule-a", Range: hcl.Range{Start: hcl.Pos{Line: 4}}},
	}

	processedIssues := ignoreIssueProcessor.ProcessIssues(issues)

	if len(processedIssues) != 2 {
		t.Fatalf("Wrong number of expected issues; got %d, wanted %d", len(processedIssues), 2)
	}
}

func TestIgnoreIssuesProcessor_ProcessRuleIgnore(t *testing.T) {
	ignoreSingle := `# tfcoach-ignore:rule-a
resource "test" "ignoreSingle"{
}
`
	ignoreMultiple := `# tfcoach-ignore:rule-a,rule-b
resource "test" "ignoreMultiple"{
}
`
	ignoreShouldNotEffectNextResource := `#tfcoach-ignore:rule-a , rule-b
resource "test" "ignoreMultiple"{
}
resource "test" "notIgnored"{
}
`

	tests := []struct {
		name, resource string
		numberOfIssues int
	}{
		{
			name:           "Ignore Single Rule",
			resource:       ignoreSingle,
			numberOfIssues: 1,
		},
		{
			name:           "Ignore Multiple Rules",
			resource:       ignoreMultiple,
			numberOfIssues: 0,
		},
		{
			name:           "Ignore Should Not Effect Next Resource",
			resource:       ignoreShouldNotEffectNextResource,
			numberOfIssues: 2,
		},
	}
	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hclFile := testutil.ParseToHcl(t, "main.tf", tt.resource)

			ignoreIssueProcessor, err := processor.NewIgnoreIssuesProcessor(".")
			if err != nil {
				t.Fatal("Setup error: ", err)
			}
			ignoreIssueProcessor.ScanFile([]byte(tt.resource), hclFile, "main.tf")

			var issues []types.Issue

			ruleA := testutil.AlwaysFlag{RuleID: "rule-a", Message: "m2"}
			ruleB := testutil.AlwaysFlag{RuleID: "rule-b", Message: "m2"}
			issues = append(issues, ruleA.Apply("main.tf", hclFile)...)
			issues = append(issues, ruleB.Apply("main.tf", hclFile)...)

			processedIssues := ignoreIssueProcessor.ProcessIssues(issues)

			if len(processedIssues) != tt.numberOfIssues {
				t.Fatalf("Wrong number of expected issues; got %d wanted %d", len(issues), tt.numberOfIssues)
			}
		})
	}
}

func createFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir -p %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
