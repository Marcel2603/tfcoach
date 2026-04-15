package types_test

import (
	"sort"
	"sync"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/types"
)

func TestSet_Len(t *testing.T) {
	set := types.Set[string]{}
	set.Add("test")
	set.Add("test2")
	if set.Len() != 2 {
		t.Errorf("Expected length 2, got %d", set.Len())
	}
}

func TestSet_Has(t *testing.T) {
	tests := []struct {
		name     string
		add      []string
		check    string
		expected bool
	}{
		{"present", []string{"a"}, "a", true},
		{"absent", []string{"a"}, "b", false},
		{"empty set", []string{}, "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := types.Set[string]{}
			for _, v := range tt.add {
				set.Add(v)
			}
			if got := set.Has(tt.check); got != tt.expected {
				t.Errorf("Has(%q) = %v, want %v", tt.check, got, tt.expected)
			}
		})
	}
}

func TestSet_Delete(t *testing.T) {
	t.Run("existing element", func(t *testing.T) {
		set := types.Set[string]{}
		set.Add("a")
		set.Delete("a")
		if set.Has("a") {
			t.Error("expected element to be deleted")
		}
		if set.Len() != 0 {
			t.Errorf("expected Len 0, got %d", set.Len())
		}
	})
	t.Run("non-existing element is no-op", func(t *testing.T) {
		set := types.Set[string]{}
		set.Add("a")
		set.Delete("b")
		if set.Len() != 1 {
			t.Errorf("expected Len 1, got %d", set.Len())
		}
	})
}

func TestSet_Values(t *testing.T) {
	set := types.Set[string]{}
	input := []string{"a", "b", "c"}
	for _, v := range input {
		set.Add(v)
	}
	got := set.Values()
	sort.Strings(got)
	sort.Strings(input)
	if len(got) != len(input) {
		t.Fatalf("expected %v, got %v", input, got)
	}
	for i := range input {
		if got[i] != input[i] {
			t.Errorf("expected %q at index %d, got %q", input[i], i, got[i])
		}
	}
}

func TestSet_ConcurrentAdd(t *testing.T) {
	set := types.Set[int]{}
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			set.Add(v)
		}(i)
	}
	wg.Wait()
	if set.Len() != 100 {
		t.Errorf("expected Len 100, got %d", set.Len())
	}
}
