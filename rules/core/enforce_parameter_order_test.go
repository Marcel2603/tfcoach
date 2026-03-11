package core_test

import (
	"strings"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/constants"
	"github.com/Marcel2603/tfcoach/internal/testutil"
	"github.com/Marcel2603/tfcoach/internal/types"
	"github.com/Marcel2603/tfcoach/rules/core"
)

func TestEnforceParameterOrder_ExpectedMETA(t *testing.T) {
	rule := core.EnforceParameterOrderRule()

	expectedMETA := types.RuleMeta{
		Title:       "Enforce Parameter Order",
		Description: "Enforce parameters should follow a consistent order",
		Severity:    constants.SeverityMedium,
		DocsURI:     strings.ReplaceAll(rule.ID(), ".", "/"),
	}

	if rule.META() != expectedMETA {
		t.Fatalf("meta mismatch; got %s, wanted %s", rule.META(), expectedMETA)
	}
}

func TestEnforceParameterOrder_AllGood(t *testing.T) {
	cases := []struct {
		name     string
		resource string
	}{
		{"no_parameters",
			`resource "aws_instance" "web" {}`,
		},
		{"with_count",
			`resource "aws_instance" "web" {
  count = 1
  ami = data.aws_ami.web.id
}`,
		},
		{
			"with_for_each",
			`resource "aws_instance" "web" {
  for_each = var.instances
  ami = each.key
  availability_zone = each.value
}`,
		},
		{
			"with_lifecycle",
			`resource "aws_instance" "web" {
  ami = 1234
  lifecycle {
    ignore_changes = [tags]
  }
}`,
		},
		{
			"with_depends_on",
			`resource "aws_instance" "web" {
  ami = 1234
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
		{
			"with_non_block_and_block_parameters",
			`resource "aws_instance" "web" {
  ami = 1234
  availability_zone = var.az
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
}`,
		},
		{
			"with_multiple_in_correct_order_1",
			`resource "aws_instance" "web" {
  count = 1
  ami = data.aws_ami.web.id
  instance_market_options {
    spot_options {
      max_price = 0.002
    }
  }
  lifecycle {
    ignore_changes = [tags]
  }
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
		{
			"with_multiple_in_correct_order_2",
			`resource "aws_instance" "web" {
  for_each = var.instances
  ami = each.key
  availability_zone = each.value
  instance_market_options {
    market_type = "spot"
  }
  depends_on = [
    aws_iam_role_policy.test
  ]
}`,
		},
	}

	rule := core.EnforceParameterOrderRule()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			issues := rule.Apply("a.tf", testutil.ParseToHcl(t, "a.tf", tt.resource))
			if len(issues) != 0 {
				t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
			}
		})
	}
}

func TestEnforceParameterOrder_ShouldComplain(t *testing.T) {
	cases := []struct {
		name       string
		issueCount int
		resource   string
	}{
		{"count_not_first",
			1,
			`resource "aws_instance" "web" {
  ami = data.aws_ami.web.id
  count = 1
}`,
		},
		{
			"for_each_not_first",
			1,
			`resource "aws_instance" "web" {
  instance_type = "t4g.nano"
  for_each = var.instances
  ami = each.key
  availability_zone = each.value
}`,
		},
		{
			"lifecycle_too_high",
			1,
			`resource "aws_instance" "web" {
  ami = 1234
  lifecycle {
    ignore_changes = [tags]
  }
  tags = {
    Name = "test"
  }
}`,
		},
		{
			"depends_on_too_high",
			1,
			`resource "aws_instance" "web" {
  ami = 1234
  depends_on = [
    aws_iam_role_policy.test
  ]
  tags = {
    Name = "test"
  }
}`,
		},
		{
			"depends_on_before_lifecycle",
			1,
			`resource "aws_instance" "web" {
  count = 1
  ami = data.aws_ami.web.id
  depends_on = [
    aws_iam_role_policy.test
  ]
  lifecycle {
    ignore_changes = [tags]
  }
}`,
		},
		{
			// TODO #20: do we want to apply the rule for nested blocks?
			"non_block_after_block_parameters",
			1,
			`resource "aws_instance" "web" {
  ami = 1234
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
  availability_zone = var.az
}`,
		},
		{
			"multiple_non_compliant_resources",
			3,
			`resource "aws_instance" "non_compliant_1" {
  ami = 1234
  instance_market_options {
    market_type = "spot"
  }
  availability_zone = var.az
}

resource "aws_instance" "compliant" {
  count = 1
  ami = 5678
  availability_zone = var.az
  instance_market_options {
    market_type = "spot"
  }
  depends_on = [aws_iam_role_policy.test]
}

resource "aws_instance" "non_compliant_2" {
  ami = 4321
  depends_on = [aws_iam_role_policy.test]
  availability_zone = var.az
}

resource "aws_instance" "non_compliant_3" {
  ami = 8765
  availability_zone = var.az
  count = 1
}`,
		},
	}

	rule := core.EnforceParameterOrderRule()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			issues := rule.Apply("a.tf", testutil.ParseToHcl(t, "a.tf", tt.resource))
			if len(issues) != tt.issueCount {
				t.Fatalf("Mismatch in reported issue count; want %d; got %d: %#v", tt.issueCount, len(issues), issues)
			}
		})
	}
}

func TestEnforceParameterOrder_FinishShouldDoNothing(t *testing.T) {
	rule := core.EnforceParameterOrderRule()

	issues := rule.Finish()
	if len(issues) != 0 {
		t.Fatalf("Issues found; expected none; got %d: %#v", len(issues), issues)
	}
}
