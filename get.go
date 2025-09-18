package mapx

// Get returns the value of the key in the map.
//
// If the key does not exist, the second return value is false.
//
// Example:
//
//	m := map[string]any{
//		"def": map[string]any{
//			"abc": 1,
//			"xyz": 2,
//		},
//	}
//	v, ok := Get(m, []string{"def", "abc"})
//	// v = 1, ok = true
func Get(m map[string]any, key []string) (any, bool) {
	if len(key) == 0 {
		return nil, false
	}

	if len(key) == 1 {
		v, ok := m[key[0]]

		return v, ok
	}

	v, ok := m[key[0]]
	if !ok {
		return nil, false
	}

	if m, ok = v.(map[string]any); !ok {
		return nil, false
	}

	return Get(m, key[1:])
}
