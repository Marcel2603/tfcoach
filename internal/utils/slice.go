//revive:disable:var-naming For now its okay to have a generic name
package utils

import (
	"cmp"
	"slices"
)

func SortAndDeduplicate[S ~[]E, E cmp.Ordered](s S) S {
	//copy to now manipulate source
	sliceCopy := make([]E, len(s))
	copy(sliceCopy, s)
	slices.Sort(sliceCopy)
	return slices.Compact(sliceCopy)
}
