package sliceutil

import (
	"reflect"
	"testing"
)

func TestAppend(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		element  int
		expected []int
	}{
		{"append to empty slice", []int{}, 1, []int{1}},
		{"append to non-empty slice", []int{1, 2}, 3, []int{1, 2, 3}},
		{"append negative number", []int{1, 2}, -1, []int{1, 2, -1}},
		{"append to large slice", make([]int, 1000), 42, append(make([]int, 1000), 42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Append(tt.input, tt.element)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Append(%v, %v) = %v; want %v", tt.input, tt.element, got, tt.expected)
			}
		})
	}
}

func BenchmarkAppend(b *testing.B) {
	slice := make([]int, 0, 1000) // Pre-allocate to minimize reallocations
	element := 42

	b.Run("small slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Append([]int{}, element)
		}
	})

	b.Run("large slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Append(slice, element)
		}
	})
}

func TestExtend(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		elements []int
		expected []int
	}{
		{"extend with non-empty slice", []int{1, 2, 3}, []int{4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{"extend with empty slice", []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{"extend empty slice with non-empty", []int{}, []int{4, 5, 6}, []int{4, 5, 6}},
		{"extend empty slice with empty slice", []int{}, []int{}, []int{}},
		{"extend with overlapping slices", []int{1, 2, 3}, []int{2, 3, 4}, []int{1, 2, 3, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Extend(tt.slice, tt.elements)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Extend(%v, %v) = %v; want %v", tt.slice, tt.elements, got, tt.expected)
			}
		})
	}
}

func BenchmarkExtend(b *testing.B) {
	original := make([]int, 1000) // A slice with 1000 elements
	elements := make([]int, 1000) // Another slice with 1000 elements

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Extend(original, elements)
	}
}

func TestPop(t *testing.T) {
	// Test pop last element (default index)
	slice := []int{1, 2, 3, 4, 5}
	element, newSlice, err := Pop(slice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if element != 5 {
		t.Errorf("expected 5, got %d", element)
	}
	if !reflect.DeepEqual(newSlice, []int{1, 2, 3, 4}) {
		t.Errorf("expected [1 2 3 4], got %v", newSlice)
	}

	// Test pop with explicit index
	element, newSlice, err = Pop(slice, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if element != 3 {
		t.Errorf("expected 3, got %d", element)
	}
	if !reflect.DeepEqual(newSlice, []int{1, 2, 4, 5}) {
		t.Errorf("expected [1 2 4 5], got %v", newSlice)
	}

	// Test empty slice
	emptySlice := []int{}
	_, _, err = Pop(emptySlice)
	if err == nil {
		t.Error("expected error for empty slice, got nil")
	}

	// Test index out of range
	_, _, err = Pop(slice, -1)
	if err == nil {
		t.Error("expected error for negative index, got nil")
	}

	_, _, err = Pop(slice, 10)
	if err == nil {
		t.Error("expected error for out-of-range index, got nil")
	}
}

func BenchmarkPopLast(b *testing.B) {
	slice := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		slice[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, newSlice, _ := Pop(slice)
		slice = newSlice
	}
}

func BenchmarkPopMiddle(b *testing.B) {
	slice := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		slice[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, newSlice, _ := Pop(slice, len(slice)/2)
		slice = newSlice
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		value     int
		wantSlice []int
		wantErr   bool
	}{
		{"Remove existing value", []int{1, 2, 3, 4}, 3, []int{1, 2, 4}, false},
		{"Remove non-existing value", []int{1, 2, 3, 4}, 5, []int{1, 2, 3, 4}, true},
		{"Remove from empty slice", []int{}, 1, []int{}, true},
		{"Remove first occurrence", []int{1, 2, 3, 2}, 2, []int{1, 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Remove(&tt.slice, tt.value)
			if !reflect.DeepEqual(tt.slice, tt.wantSlice) {
				t.Errorf("Remove() = %v, want %v", tt.slice, tt.wantSlice)
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	slice := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Remove(&slice, 500) // Remove an element from the middle
	}
}

func TestClear(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"Non-empty slice", []int{1, 2, 3}, []int{}},
		{"Empty slice", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clear(tt.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkClear(b *testing.B) {
	slice := make([]int, 1000) // A slice with 1000 elements

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Clear(slice)
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"Non-empty slice", []int{1, 2, 3}, []int{1, 2, 3}},
		{"Empty slice", []int{}, []int{}},
		{"Single element slice", []int{42}, []int{42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Copy(tt.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}

			// Ensure the slices are independent
			if len(got) > 0 {
				got[0] = 99
				if reflect.DeepEqual(got, tt.slice) {
					t.Errorf("Copy() did not create an independent slice")
				}
			}
		})
	}
}

func BenchmarkCopy(b *testing.B) {
	slice := make([]int, 1000) // A slice with 1000 elements

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Copy(slice)
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		element  int
		expected int
	}{
		{"Empty slice", []int{}, 1, 0},
		{"Single element match", []int{1}, 1, 1},
		{"Single element no match", []int{1}, 2, 0},
		{"Multiple elements match", []int{1, 2, 3, 2, 1, 2}, 2, 3},
		{"No match in slice", []int{1, 3, 5, 7}, 2, 0},
		{"All elements match", []int{2, 2, 2, 2}, 2, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Count(tt.slice, tt.element)
			if result != tt.expected {
				t.Errorf("Count(%v, %v) = %v; want %v", tt.slice, tt.element, result, tt.expected)
			}
		})
	}

	// Additional test for strings
	t.Run("String slice", func(t *testing.T) {
		slice := []string{"a", "b", "c", "b", "b"}
		element := "b"
		expected := 3
		result := Count(slice, element)
		if result != expected {
			t.Errorf("Count(%v, %q) = %v; want %v", slice, element, result, expected)
		}
	})
}

func BenchmarkCount(b *testing.B) {
	slice := make([]int, 10000)
	for i := 0; i < len(slice); i++ {
		slice[i] = i % 10 // Create a slice with repeating patterns
	}

	element := 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Count(slice, element)
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		predicate func(index int, item string) bool
		expected  []string
	}{
		{
			name:  "过滤空字符串",
			input: []string{"apple", "", "banana", "cherry", ""},
			predicate: func(index int, item string) bool {
				return item != ""
			},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:  "按索引过滤偶数索引的元素",
			input: []string{"零", "壹", "贰", "叁", "肆", "伍"},
			predicate: func(index int, item string) bool {
				return index%2 == 0
			},
			expected: []string{"零", "贰", "肆"},
		},
		{
			name:  "按内容过滤",
			input: []string{"Go", "Rust", "Python", "Java"},
			predicate: func(index int, item string) bool {
				return len(item) > 2
			},
			expected: []string{"Rust", "Python", "Java"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	// 示例数据
	slice := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		slice = append(slice, i)
	}

	// 基准测试
	b.Run("过滤偶数", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Filter(slice, func(index int, item int) bool {
				return item%2 == 0
			})
		}
	})

	b.Run("过滤大于 500 的值", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Filter(slice, func(index int, item int) bool {
				return item > 500
			})
		}
	})
}
