package mapx

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
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
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Merge() = %v", diff)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeAny(tt.args.value, tt.args.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeAny() = %v, want %v", got, tt.want)
			}
		})
	}
}
