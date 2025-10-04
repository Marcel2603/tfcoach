package core_test

import (
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestNameFormat_ExpectedMETA(t *testing.T) {
	rule := core.NamingConventionRule()

	expectedMETA := engine.RuleMeta{
		Title:       "Naming Convention",
		Description: "terraform names should only contain lowercase alphanumeric characters and underscores.",
		Severity:    "HIGH",
		DocsURL:     strings.ReplaceAll(rule.ID(), ".", "/"),
	}

	if rule.META() != expectedMETA {
		t.Fatalf("meta mismatch; got %s, wanted %s", rule.META(), expectedMETA)
	}
}

func TestNameFormat_AllGood(t *testing.T) {
	f := testutil.ParseToHcl(t, "good.tf", `
		resource "test_resource" "foo_bar_9" {}
		data "test_resource" "foo_bar_9" {}
		variable "foo_bar_9" {}
		output "foo_bar_9" {}
		locals { test = "boo" }
	`)

	rule := core.NamingConventionRule()
	issues := rule.Apply("good.tf", f)

	if len(issues) != 0 {
		t.Fatalf("expected 0 issues; got %d: %#v", len(issues), issues)
	}
}

func TestNameFormat_FailedResource(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		resource "test_resource" "UpperCasE" {}
		resource "test_resource" "d-ashes" {}
		resource "test_resource" "special_chars$%" {}
	`)

	rule := core.NamingConventionRule()
	issues := rule.Apply("bad.tf", f)

	if len(issues) != 3 {
		t.Fatalf("expected 4 issues; got %d: %#v", len(issues), issues)
	}
}

func TestNameFormat_ShouldIssueSameID(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		resource "test_resource" "UpperCasE" {}
	`)

	rule := core.NamingConventionRule()
	issues := rule.Apply("bad.tf", f)

	if issues[0].RuleID != rule.ID() {
		t.Fatal("rule id mismatch")
	}
}

func TestNameFormat_FailedVariables(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		variable "UpperCasE" {}
		variable "d-ashes" {}
		variable "special_chars$%" {}
	`)

	rule := core.NamingConventionRule()
	issues := rule.Apply("bad.tf", f)

	if len(issues) != 3 {
		t.Fatalf("expected 4 issues; got %d: %#v", len(issues), issues)
	}
}

func TestNameFormat_FailedDataObjects(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		data "test_data" "UpperCasE" {}
		data "test_data" "d-ashes" {}
		data "test_data" "special_chars$%" {}
	`)

	rule := core.NamingConventionRule()
	issues := rule.Apply("bad.tf", f)

	if len(issues) != 3 {
		t.Fatalf("expected 4 issues; got %d: %#v", len(issues), issues)
	}
}

func TestNameFormat_FailedModules(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		module "UpperCasE" {}
		module "d-ashes" {}
		module "special_chars$%" {}
	`)

	rule := core.NamingConventionRule()
	issues := rule.Apply("bad.tf", f)

	if len(issues) != 3 {
		t.Fatalf("expected 4 issues; got %d: %#v", len(issues), issues)
	}
}
