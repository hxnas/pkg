package lod

// Iif 三元运算
func Iif[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

// Iif 三元运算
func IifF[T any](cond bool, t, f func() T) T { return Iif(cond, t, f)() }

// 选择第一个不为零值的值
func Select[T comparable](vs ...T) (out T) {
	for _, v := range vs {
		if v != out {
			return v
		}
	}
	return
}
