package sleego

import (
	"reflect"
	"sort"
	"testing"
)

func sortedCopy(s []string) []string {
	c := append([]string(nil), s...)
	sort.Strings(c)
	return c
}

func equalUnordered(a, b []string) bool {
	return reflect.DeepEqual(sortedCopy(a), sortedCopy(b))
}

func TestCategoryOperator_Empty(t *testing.T) {
	op := newCategoryOperator()
	got := op.GetCategoriesOf("any")
	if len(got) != 0 {
		t.Errorf("expected empty slice, but got %v", got)
	}
}

func TestCategoryOperator_SetAndGet(t *testing.T) {
	op := newCategoryOperator()
	input := map[string][]string{
		"cat1": {"p1", "p2"},
		"cat2": {"p1"},
	}
	op.SetProcessByCategories(input)

	tests := []struct {
		proc string
		want []string
	}{
		{"p1", []string{"cat1", "cat2"}},
		{"p2", []string{"cat1"}},
		{"p3", nil},
	}

	for _, tt := range tests {
		got := op.GetCategoriesOf(tt.proc)
		if !equalUnordered(got, tt.want) {
			t.Errorf("GetCategoriesOf(%q): expected %v, got %v", tt.proc, tt.want, got)
		}
	}
}

func TestCategoryOperator_ResetOverridesPrevious(t *testing.T) {
	op := newCategoryOperator()

	op.SetProcessByCategories(map[string][]string{
		"first": {"x"},
	})
	if got := op.GetCategoriesOf("x"); !equalUnordered(got, []string{"first"}) {
		t.Fatalf("setup failed, expected [first] but got %v", got)
	}

	op.SetProcessByCategories(map[string][]string{
		"second": {"y"},
	})
	if got := op.GetCategoriesOf("x"); len(got) != 0 {
		t.Errorf("after reset, expected empty slice for 'x', but got %v", got)
	}
	if got := op.GetCategoriesOf("y"); !equalUnordered(got, []string{"second"}) {
		t.Errorf("after reset, expected [second] for 'y', but got %v", got)
	}
}
