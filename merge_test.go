package mapx

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	type args struct {
		value map[string]any
		to    map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "merge",
			args: args{
				value: map[string]any{
					"foo": "bar",
					"bar": map[string]any{
						"x": "bar",
					},
				},
				to: map[string]any{
					"foo": "bar",
					"bar": map[string]any{
						"foo": "bar",
					},
				},
			},
			want: map[string]any{
				"foo": "bar",
				"bar": map[string]any{
					"foo": "bar",
					"x":   "bar",
				},
			},
		},
		{
			name: "merge mix",
			args: args{
				value: map[string]any{
					"foo": []any{"bar"},
					"bar": map[string]any{
						"x": map[string]any{
							"foo": "bar",
						},
					},
				},
				to: map[string]any{
					"foo": "bar",
					"bar": map[string]any{
						"x": map[string]any{
							"foo": []string{"bar"},
						},
					},
				},
			},
			want: map[string]any{
				"foo": []any{"bar"},
				"bar": map[string]any{
					"x": map[string]any{
						"foo": "bar",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.args.value, tt.args.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMergeAny(t *testing.T) {
	type args struct {
		value any
		to    any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "merge",
			args: args{
				value: []any{"bar"},
				to:    []any{"foo"},
			},
			want: []any{"bar"},
		},
		{
			name: "maps merge",
			args: args{
				value: map[string]any{"x": 2},
				to:    map[string]any{"x": 1, "y": 3},
			},
			want: map[string]any{"x": 2, "y": 3},
		},
		{
			name: "map replaces scalar",
			args: args{value: map[string]any{"x": 1}, to: "old"},
			want: map[string]any{"x": 1},
		},
		{
			name: "scalar replaces map",
			args: args{value: "new", to: map[string]any{"x": 1}},
			want: "new",
		},
		{
			name: "nil replaces map",
			args: args{value: nil, to: map[string]any{"x": 1}},
			want: nil,
		},
		{
			name: "typed nil maps allocate",
			args: args{value: map[string]any(nil), to: map[string]any(nil)},
			want: map[string]any{},
		},
		{
			name: "typed nil source with untyped nil destination",
			args: args{value: map[string]any(nil), to: nil},
			want: map[string]any(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeAny(tt.args.value, tt.args.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		value map[string]any
		to    map[string]any
		want  map[string]any
	}{
		{"both nil", nil, nil, map[string]any{}},
		{"nil destination", map[string]any{"x": 1}, nil, map[string]any{"x": 1}},
		{"nil source", nil, map[string]any{"x": 1}, map[string]any{"x": 1}},
		{"empty source", map[string]any{}, map[string]any{"x": 1}, map[string]any{"x": 1}},
		{"new key", map[string]any{"x": 1}, map[string]any{"y": 2}, map[string]any{"x": 1, "y": 2}},
		{"scalar replaces map", map[string]any{"x": 1}, map[string]any{"x": map[string]any{"y": 2}}, map[string]any{"x": 1}},
		{"map replaces scalar", map[string]any{"x": map[string]any{"y": 2}}, map[string]any{"x": 1}, map[string]any{"x": map[string]any{"y": 2}}},
		{"slice replaces slice", map[string]any{"x": []any{3}}, map[string]any{"x": []any{1, 2}}, map[string]any{"x": []any{3}}},
		{"nil clears map", map[string]any{"x": nil}, map[string]any{"x": map[string]any{"y": 2}}, map[string]any{"x": nil}},
		{"typed nil source preserves map", map[string]any{"x": map[string]any(nil)}, map[string]any{"x": map[string]any{"y": 2}}, map[string]any{"x": map[string]any{"y": 2}}},
		{"typed nil destination map", map[string]any{"x": map[string]any{"y": 2}}, map[string]any{"x": map[string]any(nil)}, map[string]any{"x": map[string]any{"y": 2}}},
		{"typed map replaces rather than merges", map[string]any{"x": map[string]int{"a": 1}}, map[string]any{"x": map[string]int{"b": 2}}, map[string]any{"x": map[string]int{"a": 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.value, tt.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Merge() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMergeMutatesDestination(t *testing.T) {
	nested := map[string]any{"keep": 1, "replace": "old"}
	destination := map[string]any{"nested": nested}
	result := Merge(map[string]any{"nested": map[string]any{"replace": "new"}}, destination)
	if nested["replace"] != "new" || nested["keep"] != 1 {
		t.Fatalf("nested destination not merged in place: %v", nested)
	}
	result["added"] = true
	if destination["added"] != true {
		t.Fatal("result must be the destination map")
	}
}

func TestMergeSharesAssignedValues(t *testing.T) {
	for _, nilDestination := range []bool{true, false} {
		name := "existing destination"
		if nilDestination {
			name = "nil destination"
		}
		t.Run(name, func(t *testing.T) {
			nested := map[string]any{"x": 1}
			items := []int{1, 2}
			source := map[string]any{"nested": nested, "items": items}
			var destination map[string]any
			if !nilDestination {
				destination = map[string]any{}
			}
			result := Merge(source, destination)
			result["nested"].(map[string]any)["x"] = 2
			result["items"].([]int)[0] = 3
			if nested["x"] != 2 || items[0] != 3 {
				t.Fatal("assigned maps and slices must share storage with the source")
			}
			result["new"] = true
			if _, ok := source["new"]; ok {
				t.Fatal("top-level result must not alias the source")
			}
		})
	}
}
