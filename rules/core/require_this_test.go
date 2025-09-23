package core_test

import (
	"testing"

	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestRequireThis_AllGood(t *testing.T) {
	f := testutil.ParseToHcl(t, "good.tf", `
		resource "aws_s3_bucket" "this" {}
		resource "random_id"     "this" {}
	`)

	r := core.RequireThisRule()
	issues := r.Apply("good.tf", f)

	if len(issues) != 0 {
		t.Fatalf("expected 0 issues, got %d: %#v", len(issues), issues)
	}
}

func TestRequireThis_FlagsNonThis(t *testing.T) {
	f := testutil.ParseToHcl(t, "bad.tf", `
		resource "aws_s3_bucket" "foo" {}
	`)

	r := core.RequireThisRule()
	issues := r.Apply("bad.tf", f)

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %#v", len(issues), issues)
	}
	got := issues[0]
	if got.RuleID != "core.test_rule" {
		t.Fatalf("wrong rule id: %s", got.RuleID)
	}
	if got.Message != `resource name must be "this"` {
		t.Fatalf("wrong message: %q", got.Message)
	}
	if got.File != "bad.tf" {
		t.Fatalf("wrong file in issue: %s", got.File)
	}
	if got.Range.Start.Line != 2 {
		t.Fatalf("expected start line 2, got %d", got.Range.Start.Line)
	}
}

func TestRequireThis_MixedMultiple(t *testing.T) {
	f := testutil.ParseToHcl(t, "mix.tf", `
		variable "env" {}
		resource "aws_s3_bucket" "this" {}
		resource "aws_iam_role"  "role1" {}
		resource "aws_sqs_queue" "q" {}
		output "x" { value = 1 }
	`)

	r := core.RequireThisRule()
	issues := r.Apply("mix.tf", f)

	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %#v", len(issues), issues)
	}
	for i, is := range issues {
		if is.RuleID != "core.test_rule" {
			t.Fatalf("issue %d wrong rule id: %s", i, is.RuleID)
		}
		if is.Message != `resource name must be "this"` {
			t.Fatalf("issue %d wrong message: %q", i, is.Message)
		}
	}
}

func TestRequireThis_IgnoresNonResourceBlocks(t *testing.T) {
	f := testutil.ParseToHcl(t, "nonres.tf", `
		variable "name" {}
		output "name" { value = "x" }
		locals { a = 1 }
	`)

	r := core.RequireThisRule()
	issues := r.Apply("nonres.tf", f)

	if len(issues) != 0 {
		t.Fatalf("expected 0 issues for non-resource blocks, got %d", len(issues))
	}
}
