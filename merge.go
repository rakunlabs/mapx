package mapx

func MergeAny(value any, to any) any {
	switch value.(type) {
	case map[string]any:
		switch to.(type) {
		case map[string]any:
			return Merge(value.(map[string]any), to.(map[string]any))
		default:
			return value
		}
	default:
		return value
	}
}

func Merge(value map[string]any, to map[string]any) map[string]any {
	for k := range value {
		if _, ok := to[k]; ok {
			// check if to[k] is map
			if _, ok := to[k].(map[string]any); ok {
				// check if value[k] is map
				if _, ok := value[k].(map[string]any); ok {
					// merge
					to[k] = Merge(value[k].(map[string]any), to[k].(map[string]any))
				} else {
					to[k] = value[k]
				}
			} else {
				to[k] = value[k]
			}
		} else {
			to[k] = value[k]
		}
	}

	return to
}
