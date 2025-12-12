package core_test

import (
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestAvoidTypeInName_ExpectedMETA(t *testing.T) {
	rule := core.AvoidTypeInNameRule()

	expectedMETA := types.RuleMeta{
		Title:       "Avoid Type in Name",
		Description: "Names shouldn't repeat their type.",
		Severity:    constants.SeverityHigh,
		DocsURI:     strings.ReplaceAll(rule.ID(), ".", "/"),
	}

	if rule.META() != expectedMETA {
		t.Fatalf("meta mismatch; got %s, wanted %s", rule.META(), expectedMETA)
	}
}

func TestAvoidTypeInName_AllGood(t *testing.T) {
	f := testutil.ParseToHcl(t, "good.tf", `
		resource "test_resource" "foo_bar_9" {}
		data "test_resource" "this" {}
		variable "bucket_id" {}
		output "foo_bar_9" {}
		locals { test = "boo" }
	`)

	rule := core.AvoidTypeInNameRule()
	issues := rule.Apply("good.tf", f)

	if len(issues) != 0 {
		t.Fatalf("expected 0 issues; got %d: %#v", len(issues), issues)
	}
}

func TestAvoidTypeInName_NonCompliantNames(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		resource "test_resource" "test_bucket" {}
		resource "test_resource" "test_resource" {}
		resource "test_resource" "this_resource" {}
		data "test_resource" "bucket_test" {}
		data "aws_test_resource" "foo_aws_test" {}
		data "test_resource" "resource_test_bucket" {}
	`)

	rule := core.AvoidTypeInNameRule()
	issues := rule.Apply("bad.tf", f)

	if len(issues) != 9 {
		t.Fatalf("expected 9 issues; got %d: %#v", len(issues), issues)
	}
}

func TestAvoidTypeInNam_ShouldIssueSameID(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		resource "test_resource" "test_this" {}
	`)

	rule := core.AvoidTypeInNameRule()
	issues := rule.Apply("bad.tf", f)

	if issues[0].RuleID != rule.ID() {
		t.Fatalf("rule id mismatch; expected %s; got %s", rule.ID(), issues[0].RuleID)
	}
}
