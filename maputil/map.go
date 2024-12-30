package maputil

// GetMapDefault 获取Map value,含默认值
func GetMapDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}
