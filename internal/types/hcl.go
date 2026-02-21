//revive:disable:var-naming For now it's okay to have a generic name
package types

import "github.com/hashicorp/hcl/v2"

type DetectedBlock struct {
	Name  string
	File  string
	Range hcl.Range
}
