package engine_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
)

func TestEngine_WithStubRule(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{"a.tf": `resource "test" "test" {}`}}
	e := engine.New(src)
	e.Register(&testutil.AlwaysFlag{RuleID: "t.id", Message: "m"})
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 {
		t.Fatalf("wanted 1, got %d", len(issues))
	}
}

func TestEngine_WithMultipleStubRules(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{"a.tf": `terraform {}`}}
	e := engine.New(src)
	e.RegisterMany([]types.Rule{
		&testutil.AlwaysFlag{RuleID: "t.id1", Message: "m1"},
		&testutil.NeverFlag{RuleID: "t.x", Message: "x"},
		&testutil.AlwaysFlag{RuleID: "t.id2", Message: "m2"},
	})
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 2 {
		t.Fatalf("wanted 2, got %d", len(issues))
	}
	for _, issue := range issues {
		if !strings.HasPrefix(issue.RuleID, "t.id") {
			t.Fatalf("wanted prefix t.id, got %s", issue.RuleID)
		}
	}
}

func TestEngine_WithHclParsingError(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{"a.tf": `x`}}
	e := engine.New(src)
	e.RegisterMany([]types.Rule{
		&testutil.AlwaysFlag{RuleID: "t.id1", Message: "m1"},
		&testutil.AlwaysFlag{RuleID: "t.id2", Message: "m2"},
		&testutil.AlwaysFlag{RuleID: "t.id3", Message: "m3"},
	})
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 {
		t.Fatalf("wanted 1, got %d", len(issues))
	}
	issue := issues[0]
	if issue.RuleID != "parser" {
		t.Fatalf("wanted parser, got %s", issue.RuleID)
	}
}

func TestEngine_WithMultipleFilesAndManyStubRules(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{
		"a.tf": `locals {}`,
		"b.tf": `resource "test" "test"{}`,
		"c.tf": `terraform {}`,
	}}
	e := engine.New(src)
	var rules []types.Rule
	for i := range 100 {
		rules = append(rules, &testutil.AlwaysFlag{RuleID: strconv.Itoa(i), Message: "m"})
	}
	e.RegisterMany(rules)
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 300 {
		t.Fatalf("wanted 300, got %d", len(issues))
	}
}

func TestEngine_WithRuleThatPublishesIssuesOnFinish(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{
		"a.tf": `# empty file`,
		"b.tf": `# empty file`,
		"c.tf": `# empty file`,
	}}
	e := engine.New(src)
	e.Register(&testutil.FlagOnFinish{RuleID: "t.id", Message: "m"})
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 {
		t.Fatalf("wanted 1, got %d", len(issues))
	}
}

func TestEngine_WithRuleThatShouldBeIgnored(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{
		"a.tf": `# tfcoach-ignore:rule-0
resource "test" "test" {}`,
		"b.tf": `# tfcoach-ignore-file:rule-1
resource "test" "test" {}`,
	}}
	e := engine.New(src)
	var rules []types.Rule
	for i := range 2 {
		rules = append(rules, &testutil.AlwaysFlag{RuleID: "rule-" + strconv.Itoa(i), Message: "m"})
	}
	e.RegisterMany(rules)
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 2 {
		// a.tf violates rule-1 and b.tf violates rule-0
		t.Fatalf("wanted 2, got %d", len(issues))
	}
}
