package slices_test

import (
	"reflect"
	"strconv"
	"testing"

	"maragu.dev/gomponents/x/slices"
)

func TestMap(t *testing.T) {
	t.Run("transforms elements using index and value", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := slices.Map(input, func(i int, v int) string {
			return strconv.Itoa(i) + ":" + strconv.Itoa(v)
		})
		expected := []string{"0:1", "1:2", "2:3"}
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		var input []int
		result := slices.Map(input, func(_ int, v int) int {
			return v
		})
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("keeps elements matching predicate", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		result := slices.Filter(input, func(_ int, v int) bool {
			return v%2 == 0
		})
		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("returns nil for nil input", func(t *testing.T) {
		var input []int
		result := slices.Filter(input, func(_ int, v int) bool {
			return true
		})
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("returns empty when no elements match", func(t *testing.T) {
		input := []int{1, 3, 5}
		result := slices.Filter(input, func(_ int, v int) bool {
			return v%2 == 0
		})
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("can filter by index", func(t *testing.T) {
		input := []string{"a", "b", "c", "d", "e"}
		result := slices.Filter(input, func(i int, _ string) bool {
			return i%2 == 0
		})
		expected := []string{"a", "c", "e"}
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("accumulates values", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := slices.Reduce(input, 0, func(acc int, v int) int {
			return acc + v
		})
		if result != 15 {
			t.Errorf("expected 15, got %v", result)
		}
	})

	t.Run("returns initial value for nil slice", func(t *testing.T) {
		var input []int
		result := slices.Reduce(input, 42, func(acc int, v int) int {
			return acc + v
		})
		if result != 42 {
			t.Errorf("expected 42, got %v", result)
		}
	})

	t.Run("can reduce to different type", func(t *testing.T) {
		input := []string{"a", "bb", "ccc"}
		result := slices.Reduce(input, 0, func(acc int, v string) int {
			return acc + len(v)
		})
		if result != 6 {
			t.Errorf("expected 6, got %v", result)
		}
	})
}
