package types

import (
	"cmp"
	"encoding/json"
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

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Str)
}
