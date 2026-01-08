package slices_test

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"maragu.dev/gomponents/x/slices"
)

func TestMap(t *testing.T) {
	t.Run("transforms from one type to another with index", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := slices.Map(input, func(i int, v int) string {
			return strconv.Itoa(i) + ":" + strconv.Itoa(v)
		})
		expected := []string{"0:1", "1:2", "2:3"}
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("returns empty for nil input", func(t *testing.T) {
		var input []int
		result := slices.Map(input, func(_ int, v int) int {
			return v
		})
		if len(result) != 0 {
			t.Errorf("expected empty, got %v", result)
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

	t.Run("returns empty for nil input", func(t *testing.T) {
		var input []int
		result := slices.Filter(input, func(_ int, v int) bool {
			return true
		})
		if len(result) != 0 {
			t.Errorf("expected empty, got %v", result)
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

func ExampleMap() {
	items := []string{"party hat", "super hat"}
	e := html.Ul(slices.Map(items, func(_ int, item string) g.Node {
		return html.Li(g.Text(item))
	})...)
	_ = e.Render(os.Stdout)
	// Output: <ul><li>party hat</li><li>super hat</li></ul>
}

func ExampleMap_withIndex() {
	items := []string{"party hat", "super hat"}
	e := html.Ul(slices.Map(items, func(i int, item string) g.Node {
		return html.Li(g.Textf("%v: %v", i, item))
	})...)
	_ = e.Render(os.Stdout)
	// Output: <ul><li>0: party hat</li><li>1: super hat</li></ul>
}

func ExampleFilter() {
	type Product struct {
		Name   string
		InStock bool
	}
	products := []Product{
		{Name: "party hat", InStock: true},
		{Name: "super hat", InStock: false},
		{Name: "silly hat", InStock: true},
	}
	inStock := slices.Filter(products, func(_ int, p Product) bool {
		return p.InStock
	})
	e := html.Ul(slices.Map(inStock, func(_ int, p Product) g.Node {
		return html.Li(g.Text(p.Name))
	})...)
	_ = e.Render(os.Stdout)
	// Output: <ul><li>party hat</li><li>silly hat</li></ul>
}

func ExampleFilter_withIndex() {
	items := []string{"party hat", "super hat", "silly hat", "regular hat"}
	everyOther := slices.Filter(items, func(i int, _ string) bool {
		return i%2 == 0
	})
	e := html.Ul(slices.Map(everyOther, func(_ int, item string) g.Node {
		return html.Li(g.Text(item))
	})...)
	_ = e.Render(os.Stdout)
	// Output: <ul><li>party hat</li><li>silly hat</li></ul>
}

func ExampleReduce() {
	prices := []float64{9.99, 14.99, 24.99}
	total := slices.Reduce(prices, 0.0, func(sum float64, price float64) float64 {
		return sum + price
	})
	e := html.Span(g.Textf("Total: $%.2f", total))
	_ = e.Render(os.Stdout)
	// Output: <span>Total: $49.97</span>
}
