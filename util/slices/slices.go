// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package slices

import "sort"

// RemoveString remove string from slice if function return true.
func RemoveString(slice []string, remove func(item string) bool) []string {
	for i := 0; i < len(slice); i++ {
		if remove(slice[i]) {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

// FindString return true if target in slice, return false if not.
func FindString(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

// FindInt return true if target in slice, return false if not.
func FindInt(slice []int, target int) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

// FindUint return true if target in slice, return false if not.
func FindUint(slice []uint, target uint) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

func HasString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Map functional conversion for array items
func Map[T any, V any](arr []T, fn func(T) V) []V {
	_arr := make([]V, len(arr))
	for i, item := range arr {
		_arr[i] = fn(item)
	}
	return _arr
}

// Filter functional filter for array items
func Filter[T any](arr []T, fn func(T) bool) []T {
	var _arr []T
	for _, item := range arr {
		if fn(item) {
			_arr = append(_arr, item)
		}
	}
	return _arr
}

// Index search the index of array item with function
func Index[T any](arr []T, fn func(T) bool) int {
	for idx, item := range arr {
		if fn(item) {
			return idx
		}
	}
	return -1
}

// Find search the first item with function
func Find[T any](arr []T, fn func(T) bool) *T {
	if idx := Index(arr, fn); idx >= 0 {
		return &arr[idx]
	}
	return nil
}

// Flatten the given arr
func Flatten[T any](arr [][]T) []T {
	var _arr []T
	for _, items := range arr {
		_arr = append(_arr, items...)
	}
	return _arr
}

// All checks if all items satisfy the condition function
func All[T any](arr []T, fn func(T) bool) bool {
	for _, item := range arr {
		if !fn(item) {
			return false
		}
	}
	return true
}

// Any checks if any item satisfy the condition function
func Any[T any](arr []T, fn func(T) bool) bool {
	for _, item := range arr {
		if fn(item) {
			return true
		}
	}
	return false
}

// Count checks how many items satisfy the condition function
func Count[T any](arr []T, fn func(T) bool) int {
	cnt := 0
	for _, item := range arr {
		if fn(item) {
			cnt++
		}
	}
	return cnt
}

// GroupBy group by array items with given projection function
func GroupBy[T any, V comparable](arr []T, fn func(T) V) map[V][]T {
	m := map[V][]T{}
	for _, item := range arr {
		key := fn(item)
		if _, found := m[key]; !found {
			m[key] = []T{}
		}
		m[key] = append(m[key], item)
	}
	return m
}

// Reduce array
func Reduce[T any, V any](arr []T, reduce func(V, T) V, v V) V {
	for _, item := range arr {
		v = reduce(v, item)
	}
	return v
}

type comparableItem[T any] interface {
	Equal(T) bool
}

// Contains test if target array contains pivot
// If T is a pointer, T needs to implement Equal(T) function, otherwise the
// pointer address
// If T is not a pointer, T could be either
func Contains[T comparable](arr []T, pivot T) bool {
	for _, item := range arr {
		eq := item == pivot
		if obj, ok := any(item).(comparableItem[T]); ok {
			eq = obj.Equal(pivot)
		}
		if obj, ok := any(&item).(comparableItem[*T]); ok {
			eq = obj.Equal(&pivot)
		}
		if eq {
			return true
		}
	}
	return false
}

// Iterable .
type Iterable[T any, V any] interface {
	Next() bool
	Value() V
	*T
}

// IterToArray convert iterable to list. Next() should be called to get the first
// item
func IterToArray[U any, V any, T Iterable[U, V]](iter T) []V {
	var arr []V
	for iter != nil && iter.Next() {
		arr = append(arr, iter.Value())
	}
	return arr
}

// Sort given array
func Sort[T any](arr []T, cmp func(T, T) bool) {
	sort.Slice(arr, func(i, j int) bool {
		return cmp(arr[i], arr[j])
	})
}
