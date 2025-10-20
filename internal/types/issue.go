//revive:disable:var-naming For now it's okay to have a generic name
package types

import "github.com/hashicorp/hcl/v2"

type Issue struct {
	File    string
	Range   hcl.Range
	Message string
	RuleID  string
}
