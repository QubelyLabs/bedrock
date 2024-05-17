package util

// Get attempts to retrieve a value from a KVStore and assert it to a specific type
func GetFromMap[T any](m map[string]any, key string) T {
	value, ok := m[key].(T)
	if !ok {
		return *new(T)
	}

	return value
}

// SetToMap attempts to store a value in a KVStore
func SetToMap[T any](m map[string]any, key string, value T) {
	m[key] = value
}

// RemoveFromMap removes a key-value pair from a KVStore
func RemoveFromMap(m map[string]any, key string) {
	delete(m, key)
}
