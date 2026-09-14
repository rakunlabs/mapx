package mapx

// MergeAny calls Merge when both arguments are map[string]any; otherwise it
// returns value directly. It does not make a deep copy.
func MergeAny(value any, to any) any {
	src, srcOK := value.(map[string]any)
	dst, dstOK := to.(map[string]any)
	if srcOK && dstOK {
		return Merge(src, dst)
	}
	return value
}

// Merge recursively merges value into to, with value taking precedence.
// It modifies and returns to, allocating a map if to is nil. Nested maps are
// merged only when both values are map[string]any; other values replace the
// destination. Assigned maps and slices share storage with value.
// Inputs must not contain cycles.
func Merge(value map[string]any, to map[string]any) map[string]any {
	if to == nil {
		to = make(map[string]any, len(value))
	}

	for key, incoming := range value {
		src, srcOK := incoming.(map[string]any)
		dst, dstOK := to[key].(map[string]any)
		if srcOK && dstOK {
			to[key] = Merge(src, dst)
			continue
		}
		to[key] = incoming
	}

	return to
}
