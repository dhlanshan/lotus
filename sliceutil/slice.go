package sliceutil

import (
	"errors"
	"fmt"
	"slices"
)

// Append appends an element to the end of the slice.
// It mimics Python's list.append method.
func Append[T any](slice []T, element ...T) []T {
	return append(slice, element...)
}

// Clear removes all elements from the given slice and returns an empty slice.
func Clear[T any](slice []T) []T {
	return slice[:0]
}

// Copy creates a shallow copy of the given slice and returns it.
func Copy[T any](slice []T) []T {
	// Create a new slice with the same length and capacity as the input slice
	newSlice := make([]T, len(slice))
	copy(newSlice, slice) // Use the built-in copy function to copy elements
	return newSlice
}

// Count returns the number of occurrences of `element` in the slice `slice`.
func Count[T comparable](slice []T, element T) int {
	count := 0
	for _, v := range slice {
		if v == element {
			count++
		}
	}
	return count
}

// Extend appends all elements from the `elements` slice to the `slice`.
// It modifies the original slice in place.
func Extend[T any](slice []T, elements []T) []T {
	return append(slice, elements...)
}

// Pop removes and returns the element at the specified index from the slice.
// If no index is provided, it removes and returns the last element.
// It panics if the slice is empty or the index is out of range.
func Pop[T any](slice []T, index ...int) (T, []T, error) {
	var zeroValue T // Zero value of the generic type T

	// Check if the slice is empty
	if len(slice) == 0 {
		return zeroValue, slice, errors.New("cannot pop from an empty slice")
	}

	// Determine the index to pop
	popIndex := len(slice) - 1 // Default to the last element
	if len(index) > 0 {
		popIndex = index[0]
	}

	// Validate the index
	if popIndex < 0 || popIndex >= len(slice) {
		return zeroValue, slice, errors.New(fmt.Sprintf("index out of range: %d", popIndex))
	}

	// Extract the element
	element := slice[popIndex]

	// Remove the element by slicing
	slice = append(slice[:popIndex], slice[popIndex+1:]...)

	return element, slice, nil
}

// Remove removes the first occurrence of `value` from the slice `slice`.
// If the value is not found, it returns an error.
func Remove[T comparable](slice *[]T, value T) {
	if slice == nil || len(*slice) == 0 {
		return
	}

	// 使用 DeleteFunc 删除符合条件的元素
	*slice = slices.DeleteFunc(*slice, func(n T) bool {
		return n == value
	})
}

// Filter returns a new slice containing elements from the input slice
// that satisfy the given predicate function.
//
// The predicate function receives the index and value of each element.
// If it returns true, the element is included in the result.
//
// Example:
//
//	nums := []int{1, 2, 3, 4}
//	evens := Filter(nums, func(i, v int) bool { return v%2 == 0 })
//	// evens: [2, 4]
func Filter[S ~[]T, T any](slice S, predicate func(index int, item T) bool) []T {
	result := make([]T, 0)
	for i, item := range slice {
		if predicate(i, item) {
			result = append(result, item)
		}
	}

	return result
}

// InSlice checks whether a given value exists within a slice.
func InSlice[S ~[]E, E comparable](slice S, value E) bool {
	return slices.Contains(slice, value)
}

// Set returns a new slice with duplicate elements removed.
// The order of elements is preserved based on their first occurrence.
//
// Note: Set is NOT concurrency-safe. The caller must ensure that
// the input slice is not being modified concurrently.
func Set[T comparable](arr []T) []T {
	seen := make(map[T]struct{}, len(arr))
	result := make([]T, 0, len(arr))

	for _, v := range arr {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// SetBy returns a new slice with duplicate elements removed,
// using keyFn to extract a comparable key for each element.
// The order of elements is preserved based on their first occurrence.
//
// Note: SetBy is NOT concurrency-safe. The caller must ensure that
// the input slice is not being modified concurrently.
func SetBy[T any, K comparable](arr []T, keyFn func(T) K) []T {
	seen := make(map[K]struct{}, len(arr))
	result := make([]T, 0, len(arr))

	for _, v := range arr {
		key := keyFn(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, v)
	}
	return result
}
