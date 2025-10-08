package utils_test

import (
	"reflect"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/utils"
)

func TestSortAndDeduplicate_Strings(t *testing.T) {
	slice := []string{"abcd", "defg", "abcd", "defg"}
	result := utils.SortAndDeduplicate(slice)
	expected := []string{"abcd", "defg"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortAndDeduplicate() did not dedupe the slice correctly %#v %#v", expected, result)
	}
	if !reflect.DeepEqual(slice, []string{"abcd", "defg", "abcd", "defg"}) {
		t.Error("SortAndDeplicate() should not manipulate the original slice")
	}
}

func TestSortAndDeduplicate_Int(t *testing.T) {
	slice := []int{55, 35, 55, 35, 98, 55}
	result := utils.SortAndDeduplicate(slice)
	expected := []int{35, 55, 98}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortAndDeduplicate() did not dedupe the slice correctly %#v %#v", expected, result)
	}
	if !reflect.DeepEqual(slice, []int{55, 35, 55, 35, 98, 55}) {
		t.Error("SortAndDeplicate() should not manipulate the original slice")
	}
}
