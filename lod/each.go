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

func Concat[S ~[]E, E any](ss ...S) S                   { return slices.Concat(ss...) }
func Contains[S ~[]E, E comparable](s S, v E) bool      { return slices.Contains(s, v) }
func Delete[S ~[]E, E any](s S, i, j int) S             { return slices.Delete(s, i, j) }
func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S { return slices.DeleteFunc(s, del) }
func Insert[S ~[]E, E any](s S, i int, v ...E) S        { return slices.Insert(s, i, v...) }
