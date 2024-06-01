package lod

import "slices"

// 集合映射
func Map[S any, R any](s []S, f func(S) R) []R {
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// 集合映射
func MapIndex[S any, R any](s []S, f func(S, int) R) []R {
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = f(v, i)
	}
	return r
}

// Flatten returns a new slice concatenating the passed in slices.
func Flatten[T any](s [][]T) []T { return Concat(s...) }

// Concat returns a new slice concatenating the passed in slices.
func Concat[S ~[]E, E any](ss ...S) S { return slices.Concat(ss...) }

// Contains reports whether v is present in s.
func Contains[E comparable](s []E, v E) bool { return slices.Contains(s, v) }

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if j > len(s) or s[i:j] is not a valid slice of s.
// Delete is O(len(s)-i), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete zeroes the elements s[len(s)-(j-i):len(s)].
func Delete[S ~[]E, E any](s S, i, j int) S { return slices.Delete(s, i, j) }

// DeleteFunc removes any elements from s for which del returns true,
// returning the modified slice.
// DeleteFunc zeroes the elements between the new length and the original length.
func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S { return slices.DeleteFunc(s, del) }

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// The elements at s[i:] are shifted up to make room.
// In the returned slice r, r[i] == v[0],
// and r[i+len(v)] == value originally at r[i].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
func Insert[S ~[]E, E any](s S, i int, v ...E) S { return slices.Insert(s, i, v...) }
