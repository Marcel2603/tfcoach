package core_test

import (
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestEnforceVariableDescription_ExpectedMETA(t *testing.T) {
	rule := core.EnforceVariableDescriptionRule()

	expectedMETA := types.RuleMeta{
		Title:       "Enforce Variable Description",
		Description: "To understand what that variable does (even if it seems trivial), always add a description",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(rule.ID(), ".", "/"),
	}

	if rule.META() != expectedMETA {
		t.Fatalf("meta mismatch; got %s, wanted %s", rule.META(), expectedMETA)
	}
}

func TestEnforceVariableDescription_AllGood(t *testing.T) {
	f := testutil.ParseToHcl(t, "good.tf", `
variable "a" {
  type = string
  description = "a description"
}
resource "aws_s3_bucket" "this" {}
variable "b" {
  type = string
  description = "b"
}
data "test_resource" "some_test" {}
`)

	rule := core.EnforceVariableDescriptionRule()
	issues := rule.Apply("good.tf", f)

	if len(issues) != 0 {
		t.Fatalf("expected 0 issues; got %d: %#v", len(issues), issues)
	}
}

func TestEnforceVariableDescription_ShouldComplainOnMissingDescription(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
variable "a" {
  type = string
}
variable "b" {
  type = string
  description = "b"
}
`)

	rule := core.EnforceVariableDescriptionRule()
	issues := rule.Apply("bad.tf", f)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue; got %d: %#v", len(issues), issues)
	}
	if issues[0].RuleID != rule.ID() {
		t.Fatalf("rule id mismatch; got %s, want %s", issues[0].RuleID, rule.ID())
	}
}

func TestEnforceVariableDescription_ShouldComplainOnEmptyDescription(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
variable "a" {
  type = string
  description = ""
}
variable "b" {
  type = string
  description = "b"
}
`)

	rule := core.EnforceVariableDescriptionRule()
	issues := rule.Apply("bad.tf", f)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue; got %d: %#v", len(issues), issues)
	}
	if issues[0].RuleID != rule.ID() {
		t.Fatalf("rule id mismatch; got %s, want %s", issues[0].RuleID, rule.ID())
	}
}
