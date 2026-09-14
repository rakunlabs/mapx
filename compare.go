package mapx

import (
	"fmt"
	"reflect"
	"strings"
)

// IsSubset returns nil when subset is contained in set, otherwise an error.
// Both arguments must be map[string]any or both must be []any.
// Recursive input collections must not contain cycles.
func IsSubset(set, subset any, opts ...OptionCompare) error {
	setType := reflect.TypeOf(set)
	if setType == nil {
		return fmt.Errorf("set is not map or slice")
	}

	// if set is map
	if setType.Kind() == reflect.Map {
		// Safe type assertion for maps
		mapSet, ok1 := set.(map[string]any)
		mapSubset, ok2 := subset.(map[string]any)
		if !ok1 || !ok2 {
			return fmt.Errorf("cannot convert arguments to map[string]any")
		}
		return IsMapSubset(mapSet, mapSubset, opts...)
	}

	// if set is slice
	if setType.Kind() == reflect.Slice {
		// Safe type assertion for slices
		sliceSet, ok1 := set.([]any)
		sliceSubset, ok2 := subset.([]any)
		if !ok1 || !ok2 {
			return fmt.Errorf("cannot convert arguments to []any")
		}
		return IsSliceSubset(sliceSet, sliceSubset, opts...)
	}

	return fmt.Errorf("set is not map or slice")
}

// IsMapSubset returns nil when mapSubset is contained in mapSet, otherwise an error.
// Nested collections must be map[string]any or []any. Nil values match only nil.
func IsMapSubset(mapSet, mapSubset map[string]any, opts ...OptionCompare) error {
	opt := newOptionCompare(opts...)
	return isMapSubset(mapSet, mapSubset, opt)
}

func isMapSubset(mapSet, mapSubset map[string]any, opt *optionCompare) error {
	if len(mapSubset) > len(mapSet) {
		return fmt.Errorf("mapSubset length is greater than mapSet length")
	}

	for k := range mapSubset {
		var actualKey string
		var ok bool

		if opt.CaseInsensitive {
			var err error
			actualKey, ok, err = findKeyCaseInsensitive(mapSet, k)
			if err != nil {
				return err
			}
		} else {
			actualKey = k
			_, ok = mapSet[k]
		}

		if !ok {
			return fmt.Errorf("key %s not found in mapSet", k)
		}

		if err := compareValue(mapSet[actualKey], mapSubset[k], opt); err != nil {
			return fmt.Errorf("key %s: %w", k, err)
		}
	}

	return nil
}

// IsSliceSubset returns nil when every subset element matches a set element.
// Order and duplicate counts are ignored; matching elements are not consumed.
func IsSliceSubset(sliceSet []any, sliceSubset []any, opts ...OptionCompare) error {
	return isSliceSubset(sliceSet, sliceSubset, newOptionCompare(opts...))
}

func isSliceSubset(sliceSet, sliceSubset []any, opt *optionCompare) error {
	for i, v := range sliceSubset {
		if err := isSliceContains(sliceSet, v, opt); err != nil {
			return fmt.Errorf("subset index %d: %w", i, err)
		}
	}

	return nil
}

// IsSliceContains returns nil when slice contains a matching value, otherwise an error.
// Map and slice elements are matched using subset semantics, not exact equality.
func IsSliceContains(slice []any, value any, opts ...OptionCompare) error {
	return isSliceContains(slice, value, newOptionCompare(opts...))
}

func isSliceContains(slice []any, value any, opt *optionCompare) error {
	for _, v := range slice {
		if compareValue(v, value, opt) == nil {
			return nil
		}
	}

	return fmt.Errorf("value not found in slice")
}

func compareValue(set, subset any, opt *optionCompare) error {
	if set == nil || subset == nil {
		if set == nil && subset == nil {
			return nil
		}
		return fmt.Errorf("nil and non-nil values do not match")
	}
	if !opt.WeakType && reflect.TypeOf(set) != reflect.TypeOf(subset) {
		return fmt.Errorf("value types differ: %T and %T", set, subset)
	}
	switch sub := subset.(type) {
	case map[string]any:
		if m, ok := set.(map[string]any); ok {
			return isMapSubset(m, sub, opt)
		}
	case []any:
		if s, ok := set.([]any); ok {
			return isSliceSubset(s, sub, opt)
		}
	default:
		if valuesEqual(subset, set, opt) {
			return nil
		}
	}
	return fmt.Errorf("values do not match or have unsupported types: %T and %T", set, subset)
}

// //////////////////////////////////////////////////////////////

// findKeyCaseInsensitive finds a key in the map with case-insensitive matching
func findKeyCaseInsensitive(m map[string]any, key string) (string, bool, error) {
	// First try exact match
	if _, ok := m[key]; ok {
		return key, true, nil
	}

	// Try case-insensitive match
	var match string
	found := false
	for k := range m {
		if strings.EqualFold(k, key) {
			if found {
				return "", false, fmt.Errorf("ambiguous case-insensitive key %q", key)
			}
			match, found = k, true
		}
	}

	return match, found, nil
}

// valuesEqual compares two values with the given options.
func valuesEqual(a, b any, opt *optionCompare) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	if !va.Comparable() || !vb.Comparable() {
		return false
	}

	switch {
	case opt.CaseInsensitive && aType.Kind() == reflect.String && bType.Kind() == reflect.String:
		return strings.EqualFold(va.String(), vb.String())
	case aType == bType:
		return va.Interface() == vb.Interface()
	case opt.WeakType:
		// Compare string representations for weak typing
		return fmt.Sprintf("%v", va.Interface()) == fmt.Sprintf("%v", vb.Interface())
	default:
		return false
	}
}

// //////////////////////////////////////////////////////////////

type optionCompare struct {
	CaseInsensitive bool
	WeakType        bool
}

func newOptionCompare(opts ...OptionCompare) *optionCompare {
	o := &optionCompare{
		CaseInsensitive: false,
		WeakType:        true,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}

	return o
}

// OptionCompare configures subset and containment comparisons.
type OptionCompare func(*optionCompare)

// WithCaseInsensitive enables case-insensitive map keys and string values.
// Exact keys take precedence; multiple non-exact matches produce an error.
func WithCaseInsensitive(v bool) OptionCompare {
	return func(o *optionCompare) {
		o.CaseInsensitive = v
	}
}

// WithWeakType compares different comparable types using their %v representations.
// It defaults to true and does not parse or normalize numeric strings.
func WithWeakType(v bool) OptionCompare {
	return func(o *optionCompare) {
		o.WeakType = v
	}
}
