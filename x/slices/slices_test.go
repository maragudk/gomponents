package slices_test

import (
	"strconv"
	"testing"

	"maragu.dev/is"

	"maragu.dev/gomponents/x/slices"
)

func TestMap(t *testing.T) {
	t.Run("maps integers to strings with index", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := slices.Map(input, func(i int, v int) string {
			return strconv.Itoa(i) + ":" + strconv.Itoa(v)
		})
		is.EqualSlice(t, []string{"0:1", "1:2", "2:3"}, result)
	})

	t.Run("maps strings to integers", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		result := slices.Map(input, func(i int, s string) int {
			return len(s) + i
		})
		is.EqualSlice(t, []int{5, 6, 6}, result)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var input []int
		result := slices.Map(input, func(i int, v int) string {
			return strconv.Itoa(v)
		})
		is.EqualSlice(t, []string{}, result)
	})

	t.Run("handles nil slice", func(t *testing.T) {
		var input []int
		result := slices.Map(input, func(i int, v int) string {
			return strconv.Itoa(v)
		})
		is.Nil(t, result)
	})

	t.Run("maps struct to another struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		type Display struct {
			ID   int
			Text string
		}

		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}

		displays := slices.Map(people, func(i int, p Person) Display {
			return Display{
				ID:   i,
				Text: p.Name + " (" + strconv.Itoa(p.Age) + ")",
			}
		})

		expected := []Display{
			{ID: 0, Text: "Alice (30)"},
			{ID: 1, Text: "Bob (25)"},
		}
		is.EqualSlice(t, expected, displays)
	})

	t.Run("preserves order", func(t *testing.T) {
		input := []int{3, 1, 4, 1, 5, 9}
		result := slices.Map(input, func(i int, v int) int {
			return v * 2
		})
		is.EqualSlice(t, []int{6, 2, 8, 2, 10, 18}, result)
	})
}

func TestFilter(t *testing.T) {
	t.Run("filters even numbers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		result := slices.Filter(input, func(v int) bool {
			return v%2 == 0
		})
		is.EqualSlice(t, []int{2, 4, 6}, result)
	})

	t.Run("filters strings by length", func(t *testing.T) {
		input := []string{"a", "ab", "abc", "abcd", "abcde"}
		result := slices.Filter(input, func(s string) bool {
			return len(s) >= 3
		})
		is.EqualSlice(t, []string{"abc", "abcd", "abcde"}, result)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var input []int
		result := slices.Filter(input, func(v int) bool {
			return v > 0
		})
		is.EqualSlice(t, []int{}, result)
	})

	t.Run("handles nil slice", func(t *testing.T) {
		var input []int
		result := slices.Filter(input, func(v int) bool {
			return v > 0
		})
		is.Nil(t, result)
	})

	t.Run("filters all elements when none match", func(t *testing.T) {
		input := []int{1, 3, 5, 7}
		result := slices.Filter(input, func(v int) bool {
			return v%2 == 0
		})
		is.EqualSlice(t, []int{}, result)
	})

	t.Run("keeps all elements when all match", func(t *testing.T) {
		input := []int{2, 4, 6, 8}
		result := slices.Filter(input, func(v int) bool {
			return v%2 == 0
		})
		is.EqualSlice(t, []int{2, 4, 6, 8}, result)
	})

	t.Run("filters structs", func(t *testing.T) {
		type Product struct {
			Name  string
			Price float64
		}

		products := []Product{
			{Name: "Apple", Price: 0.5},
			{Name: "Banana", Price: 0.3},
			{Name: "Orange", Price: 0.8},
			{Name: "Grape", Price: 1.2},
		}

		expensive := slices.Filter(products, func(p Product) bool {
			return p.Price > 0.5
		})

		expected := []Product{
			{Name: "Orange", Price: 0.8},
			{Name: "Grape", Price: 1.2},
		}
		is.EqualSlice(t, expected, expensive)
	})

	t.Run("preserves order", func(t *testing.T) {
		input := []int{5, 2, 8, 1, 9, 3, 6}
		result := slices.Filter(input, func(v int) bool {
			return v > 4
		})
		is.EqualSlice(t, []int{5, 8, 9, 6}, result)
	})
}

func TestReduce(t *testing.T) {
	t.Run("sums integers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := slices.Reduce(input, 0, func(acc int, v int) int {
			return acc + v
		})
		is.Equal(t, 15, result)
	})

	t.Run("concatenates strings", func(t *testing.T) {
		input := []string{"hello", " ", "world"}
		result := slices.Reduce(input, "", func(acc string, v string) string {
			return acc + v
		})
		is.Equal(t, "hello world", result)
	})

	t.Run("finds maximum", func(t *testing.T) {
		input := []int{3, 7, 2, 9, 1, 5}
		result := slices.Reduce(input, input[0], func(max int, v int) int {
			if v > max {
				return v
			}
			return max
		})
		is.Equal(t, 9, result)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var input []int
		result := slices.Reduce(input, 100, func(acc int, v int) int {
			return acc + v
		})
		is.Equal(t, 100, result)
	})

	t.Run("handles nil slice", func(t *testing.T) {
		var input []int
		result := slices.Reduce(input, 42, func(acc int, v int) int {
			return acc + v
		})
		is.Equal(t, 42, result)
	})

	t.Run("calculates product", func(t *testing.T) {
		input := []int{2, 3, 4}
		result := slices.Reduce(input, 1, func(acc int, v int) int {
			return acc * v
		})
		is.Equal(t, 24, result)
	})

	t.Run("reduces to different type", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		result := slices.Reduce(input, 0, func(acc int, v string) int {
			return acc + len(v)
		})
		is.Equal(t, 3, result)
	})

	t.Run("builds map from slice", func(t *testing.T) {
		type Item struct {
			Key   string
			Value int
		}

		items := []Item{
			{Key: "a", Value: 1},
			{Key: "b", Value: 2},
			{Key: "c", Value: 3},
		}

		result := slices.Reduce(items, make(map[string]int), func(m map[string]int, item Item) map[string]int {
			m[item.Key] = item.Value
			return m
		})

		expected := map[string]int{"a": 1, "b": 2, "c": 3}
		is.Equal(t, len(expected), len(result))
		for k, v := range expected {
			is.Equal(t, v, result[k])
		}
	})

	t.Run("counts occurrences", func(t *testing.T) {
		input := []string{"apple", "banana", "apple", "orange", "banana", "apple"}
		counts := slices.Reduce(input, make(map[string]int), func(acc map[string]int, v string) map[string]int {
			acc[v]++
			return acc
		})

		is.Equal(t, 3, counts["apple"])
		is.Equal(t, 2, counts["banana"])
		is.Equal(t, 1, counts["orange"])
	})

	t.Run("builds nested structure", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		people := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 30},
		}

		// Group by age
		byAge := slices.Reduce(people, make(map[int][]string), func(acc map[int][]string, p Person) map[int][]string {
			acc[p.Age] = append(acc[p.Age], p.Name)
			return acc
		})

		is.EqualSlice(t, []string{"Alice", "Charlie"}, byAge[30])
		is.EqualSlice(t, []string{"Bob"}, byAge[25])
	})
}

func TestCombinedOperations(t *testing.T) {
	t.Run("map then filter", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		doubled := slices.Map(input, func(_ int, v int) int {
			return v * 2
		})
		result := slices.Filter(doubled, func(v int) bool {
			return v > 5
		})
		is.EqualSlice(t, []int{6, 8, 10}, result)
	})

	t.Run("filter then reduce", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		evens := slices.Filter(input, func(v int) bool {
			return v%2 == 0
		})
		sum := slices.Reduce(evens, 0, func(acc int, v int) int {
			return acc + v
		})
		is.Equal(t, 12, sum)
	})

	t.Run("map filter reduce pipeline", func(t *testing.T) {
		type Sale struct {
			Product  string
			Quantity int
			Price    float64
		}

		sales := []Sale{
			{Product: "Widget", Quantity: 10, Price: 5.0},
			{Product: "Gadget", Quantity: 5, Price: 10.0},
			{Product: "Doohickey", Quantity: 2, Price: 15.0},
			{Product: "Widget", Quantity: 3, Price: 5.0},
		}

		// Calculate total revenue for Widgets
		revenues := slices.Map(sales, func(_ int, s Sale) float64 {
			return float64(s.Quantity) * s.Price
		})

		widgetSales := slices.Filter(sales, func(s Sale) bool {
			return s.Product == "Widget"
		})

		widgetRevenues := slices.Map(widgetSales, func(_ int, s Sale) float64 {
			return float64(s.Quantity) * s.Price
		})

		totalWidgetRevenue := slices.Reduce(widgetRevenues, 0.0, func(acc float64, v float64) float64 {
			return acc + v
		})

		is.Equal(t, 65.0, totalWidgetRevenue)
	})
}