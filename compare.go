package mapx

import (
	"fmt"
	"reflect"
	"strings"
)

func IsSubset(set, subset any, opts ...OptionCompare) error {
	setType := reflect.TypeOf(set)

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

// IsMapSubset returns true if mapSubset is a subset of mapSet otherwise false.
func IsMapSubset(mapSet, mapSubset map[string]any, opts ...OptionCompare) error {
	opt := newOptionCompare(opts...)

	if len(mapSubset) > len(mapSet) {
		return fmt.Errorf("mapSubset length is greater than mapSet length")
	}

	for k := range mapSubset {
		var actualKey string
		var ok bool

		if opt.CaseInsensitive {
			actualKey, ok = findKeyCaseInsensitive(mapSet, k)
		} else {
			actualKey = k
			_, ok = mapSet[k]
		}

		if !ok {
			return fmt.Errorf("key %s not found in mapSet", k)
		}

		mapSubsetType := reflect.TypeOf(mapSubset[k])
		mapSetType := reflect.TypeOf(mapSet[actualKey])

		if !opt.WeakType && mapSubsetType != mapSetType {
			return fmt.Errorf("type of key %s is not the same in mapSet and mapSubset", k)
		}

		if mapSubsetType.Comparable() && mapSetType.Comparable() {
			if !valuesEqual(mapSubset[k], mapSet[actualKey], opt) {
				return fmt.Errorf("value of key %s is not equal", k)
			}
			continue
		}

		if mapSubsetType.Kind() == reflect.Map && mapSetType.Kind() == reflect.Map {
			if err := IsMapSubset(mapSet[actualKey].(map[string]any), mapSubset[k].(map[string]any), opts...); err != nil {
				return fmt.Errorf("key %s: %v", k, err)
			}

			continue
		}

		if mapSubsetType.Kind() == reflect.Slice && mapSetType.Kind() == reflect.Slice {
			if err := IsSliceSubset(mapSet[actualKey].([]any), mapSubset[k].([]any), opts...); err != nil {
				return fmt.Errorf("key %s: %v", k, err)
			}

			continue
		}

		return fmt.Errorf("value of key %s is not related", k)
	}

	return nil
}

// IsSliceSubset returns true if sliceSubset is a subset of sliceSet otherwise false.
func IsSliceSubset(sliceSet []any, sliceSubset []any, opts ...OptionCompare) error {
	for _, v := range sliceSubset {
		if err := IsSliceContains(sliceSet, v, opts...); err != nil {
			return fmt.Errorf("value %v not found in sliceSet", v)
		}
	}

	return nil
}

// IsSliceContains returns true if slice contains value otherwise false.
func IsSliceContains(slice []any, value any, opts ...OptionCompare) error {
	opt := newOptionCompare(opts...)

	for _, v := range slice {
		valueType := reflect.TypeOf(value)
		vType := reflect.TypeOf(v)

		// Handle comparable values with options
		if valueType.Comparable() && vType.Comparable() {
			if valuesEqual(value, v, opt) {
				return nil
			}

			continue
		}

		// Handle maps with weak typing
		if valueType.Kind() == reflect.Map && vType.Kind() == reflect.Map {
			mapValue, ok1 := value.(map[string]any)
			mapV, ok2 := v.(map[string]any)
			if ok1 && ok2 {
				if err := IsMapSubset(mapV, mapValue, opts...); err == nil {
					return nil
				}
			}

			continue
		}

		// Handle slices with weak typing
		if valueType.Kind() == reflect.Slice && vType.Kind() == reflect.Slice {
			sliceValue, ok1 := value.([]any)
			sliceV, ok2 := v.([]any)
			if ok1 && ok2 {
				if err := IsSliceSubset(sliceV, sliceValue, opts...); err == nil {
					return nil
				}
			}

			continue
		}
	}

	return fmt.Errorf("value not found in slice")
}

// //////////////////////////////////////////////////////////////

// findKeyCaseInsensitive finds a key in the map with case-insensitive matching
func findKeyCaseInsensitive(m map[string]any, key string) (string, bool) {
	// First try exact match
	if _, ok := m[key]; ok {
		return key, true
	}

	// Try case-insensitive match
	for k := range m {
		if strings.EqualFold(k, key) {
			return k, true
		}
	}

	return "", false
}

// valuesEqual compares two values with the given options.
func valuesEqual(a, b any, opt *optionCompare) bool {
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)

	if !aType.Comparable() || !bType.Comparable() {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

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
		opt(o)
	}

	return o
}

type OptionCompare func(*optionCompare)

func WithCaseInsensitive(v bool) OptionCompare {
	return func(o *optionCompare) {
		o.CaseInsensitive = v
	}
}

func WithWeakType(v bool) OptionCompare {
	return func(o *optionCompare) {
		o.WeakType = v
	}
}
