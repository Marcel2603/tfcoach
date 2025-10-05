//go:build tools

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <file.hcl>\n", os.Args[0])
		os.Exit(2)
	}
	filename := os.Args[1]

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		fmt.Fprintln(os.Stderr, diags.Error())
		os.Exit(1)
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		fmt.Fprintf(os.Stderr, "file %q is not native HCL syntax (maybe JSON?)\n", filename)
		os.Exit(1)
	}

	fmt.Printf("FILE %s\n", filename)
	printBody(body, 0)
}

type item struct {
	kind string
	rng  hcl.Range
	attr *hclsyntax.Attribute
	blk  *hclsyntax.Block
}

func printBody(body *hclsyntax.Body, indent int) {
	spaces := strings.Repeat(" ", indent)

	items := make([]item, 0, len(body.Attributes)+len(body.Blocks))

	for _, a := range body.Attributes {
		items = append(items, item{
			kind: "attr",
			rng:  a.Range(),
			attr: a,
		})
	}
	for _, b := range body.Blocks {
		items = append(items, item{
			kind: "block",
			rng:  b.TypeRange,
			blk:  b,
		})
	}

	ctx := &hcl.EvalContext{}

	for _, it := range items {
		switch it.kind {
		case "attr":
			val, diags := it.attr.Expr.Value(ctx)
			if diags.HasErrors() {
				fmt.Printf("%sAttribute %s = <unevaluable: %s>\n", spaces, it.attr.Name, diags.Error())
			} else {
				fmt.Printf("%sAttribute %s = %#v\n", spaces, it.attr.Name, val)
			}

		case "block":
			fmt.Printf("%sBlock %s\n", spaces, it.blk.Type)
			if len(it.blk.Labels) > 0 {
				fmt.Printf("%s  Labels %#v\n", spaces, it.blk.Labels)
			}
			printBody(it.blk.Body, indent+2)
		}
	}
}
