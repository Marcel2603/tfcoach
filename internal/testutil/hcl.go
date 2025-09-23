package testutil

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func ParseToHcl(t *testing.T, filename, src string) *hcl.File {
	t.Helper()
	f, diags := hclsyntax.ParseConfig([]byte(src), filename, hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatalf("parse error: %v", diags.Error())
	}
	return f
}
