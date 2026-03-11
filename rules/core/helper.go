package core

import "github.com/hashicorp/hcl/v2/hclsyntax"

func nameOf(block *hclsyntax.Block) string {
	if len(block.Labels) == 0 {
		return ""
	}

	// <block_type> "<label1>" "<label2>"
	return block.Labels[len(block.Labels)-1]
}
