package mapx

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestComparisonEdgeCases(t *testing.T) {
	type box struct{ Value any }
	type namedString string
	tests := []struct {
		name string
		run  func() error
		want bool
	}{
		{"nil set", func() error { return IsSubset(nil, nil) }, false},
		{"scalar set", func() error { return IsSubset(1, 1) }, false},
		{"typed top-level map", func() error { return IsSubset(map[string]int{"x": 1}, map[string]int{"x": 1}) }, false},
		{"typed top-level slice", func() error { return IsSubset([]string{"x"}, []string{"x"}) }, false},
		{"top-level map", func() error { return IsSubset(map[string]any{"x": 1}, map[string]any{"x": 1}) }, true},
		{"top-level slice", func() error { return IsSubset([]any{1, 2}, []any{2}) }, true},
		{"nil maps", func() error { return IsMapSubset(nil, nil) }, true},
		{"nil slices", func() error { return IsSliceSubset(nil, nil) }, true},
		{"nil values", func() error { return IsMapSubset(map[string]any{"x": nil}, map[string]any{"x": nil}) }, true},
		{"missing versus nil", func() error { return IsMapSubset(map[string]any{"y": nil}, map[string]any{"x": nil}) }, false},
		{"nil versus value", func() error { return IsMapSubset(map[string]any{"x": nil}, map[string]any{"x": 1}) }, false},
		{"value versus nil", func() error { return IsMapSubset(map[string]any{"x": 1}, map[string]any{"x": nil}) }, false},
		{"nil slice match", func() error { return IsSliceContains([]any{1, nil}, nil) }, true},
		{"nil before slice match", func() error { return IsSliceContains([]any{nil, 1}, 1) }, true},
		{"empty slice contains nil", func() error { return IsSliceContains(nil, nil) }, false},
		{"typed nil is not nil", func() error { return IsSliceContains([]any{(*int)(nil)}, nil) }, false},
		{"typed nil match", func() error { return IsSliceContains([]any{(*int)(nil)}, (*int)(nil)) }, true},
		{"nested typed map", func() error {
			return IsMapSubset(map[string]any{"x": map[string]int{"a": 1}}, map[string]any{"x": map[string]int{"a": 1}})
		}, false},
		{"nested typed slice", func() error {
			return IsMapSubset(map[string]any{"x": []string{"a"}}, map[string]any{"x": []string{"a"}})
		}, false},
		{"dynamic uncomparable struct", func() error { return IsSliceContains([]any{box{[]int{1}}}, box{[]int{1}}) }, false},
		{"dynamic uncomparable array", func() error { return IsSliceContains([]any{[1]any{[]int{1}}}, [1]any{[]int{1}}) }, false},
		{"uncomparable before match", func() error { return IsSliceContains([]any{box{[]int{1}}, 1}, 1) }, true},
		{"comparable struct", func() error { return IsSliceContains([]any{box{1}}, box{1}) }, true},
		{"nested slice subset", func() error { return IsSliceContains([]any{[]any{1, 2}}, []any{2}) }, true},
		{"nested map subset", func() error { return IsSliceContains([]any{map[string]any{"x": 1, "y": 2}}, map[string]any{"x": 1}) }, true},
		{"nested strict comparison", func() error {
			return IsSliceContains([]any{map[string]any{"x": 1}}, map[string]any{"x": "1"}, WithWeakType(false))
		}, false},
		{"duplicate subset", func() error { return IsSliceSubset([]any{1}, []any{1, 1}) }, true},
		{"weak number", func() error { return IsSliceContains([]any{1}, "1") }, true},
		{"weak json number", func() error { return IsSliceContains([]any{1}, json.Number("1")) }, true},
		{"weak no numeric normalization", func() error { return IsSliceContains([]any{1}, "1.0") }, false},
		{"case-insensitive value", func() error { return IsSliceContains([]any{"Go"}, "go", WithCaseInsensitive(true)) }, true},
		{"strict named string", func() error {
			return IsSliceContains([]any{namedString("Go")}, "go", WithCaseInsensitive(true), WithWeakType(false))
		}, false},
		{"nil option", func() error { return IsSliceContains([]any{1}, 1, nil) }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); (err == nil) != tt.want {
				t.Fatalf("error = %v, want match = %v", err, tt.want)
			}
		})
	}
}

func TestCaseInsensitiveKeyCollisions(t *testing.T) {
	set := map[string]any{"Name": "a", "NAME": "b"}
	for i := 0; i < 100; i++ {
		err := IsMapSubset(set, map[string]any{"name": "a"}, WithCaseInsensitive(true))
		if err == nil || !strings.Contains(err.Error(), "ambiguous") {
			t.Fatalf("expected ambiguity error, got %v", err)
		}
	}
	if err := IsMapSubset(set, map[string]any{"Name": "a"}, WithCaseInsensitive(true)); err != nil {
		t.Fatalf("exact key should take precedence: %v", err)
	}
	if err := IsMapSubset(set, map[string]any{"Name": "b"}, WithCaseInsensitive(true)); err == nil {
		t.Fatal("exact key must not fall back to another case variant")
	}
	if err := IsMapSubset(map[string]any{"": 1}, map[string]any{"": 1}, WithCaseInsensitive(true)); err != nil {
		t.Fatalf("empty key should match: %v", err)
	}
}

func FuzzJSONSubsetReflexive(f *testing.F) {
	f.Add(`{"x":null,"items":[null,1,{"name":"Go"}]}`)
	f.Add(`[null,{"x":[1,2]},true]`)
	f.Add(`{}`)
	f.Fuzz(func(t *testing.T, input string) {
		var value any
		if json.Unmarshal([]byte(input), &value) != nil {
			return
		}
		switch value.(type) {
		case map[string]any, []any:
			if err := IsSubset(value, value); err != nil {
				t.Fatalf("JSON collection must be its own subset: %v", err)
			}
		}
	})
}
