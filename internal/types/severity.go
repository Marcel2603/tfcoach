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

func (s Severity) ColoredString() string {
	str := s.String()
	switch s.Priority {
	case 1:
		return color.HiRedString(str)
	case 2:
		return color.HiYellowString(str)
	case 3:
		return color.HiWhiteString(str)
	default:
		return str
	}
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Str)
}
