//revive:disable:var-naming For now it's okay to have a generic name
package types

import (
	"cmp"
	"encoding/json"

	"github.com/fatih/color"
)

type Severity struct {
	Str      string
	Priority int
}

func (s Severity) Cmp(other Severity) int {
	return cmp.Compare(s.Priority, other.Priority)
}

func (s Severity) String() string {
	return s.Str
}

func (s Severity) Color() color.Attribute {
	switch s.Priority {
	case 1:
		return color.FgHiRed
	case 2:
		return color.FgHiYellow
	case 3:
		return color.FgHiWhite
	default:
		return color.Reset
	}
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Str)
}
