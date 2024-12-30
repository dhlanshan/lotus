package ops

// Ternary 三目运算
func Ternary[T any](exp bool, e1, e2 T) T {
	if exp {
		return e1
	}
	return e2
}
