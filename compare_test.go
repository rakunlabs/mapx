package mapx

import (
	"encoding/json"
	"testing"
)

func TestIsMapSubset(t *testing.T) {
	type args struct {
		m   map[string]any
		sub map[string]any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "simple one test",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
				},
				sub: map[string]any{
					"abc": 1,
				},
			},
			want: true,
		},
		{
			name: "mix type",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
					"def": map[string]any{
						"abc": 1,
						"xyz": 2,
					},
				},
				sub: map[string]any{
					"abc": 1,
					"def": map[string]any{
						"abc": 1,
					},
				},
			},
			want: true,
		},
		{
			name: "mix type",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
					"def": map[string]any{
						"abc": 1,
						"xyz": 2,
					},
				},
				sub: map[string]any{
					"abc": 1,
					"def": map[string]any{
						"abc": json.Number("1"),
					},
				},
			},
			want: true,
		},
		{
			name: "mix type",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
					"def": map[string]any{
						"abc": 1,
						"xyz": 2,
					},
				},
				sub: map[string]any{
					"abc": 1,
					"def": map[string]any{
						"abc": "sfdfds",
					},
				},
			},
			want: false,
		},
		{
			name: "mix type in sub different type",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
					"def": map[string]any{
						"abc": "sfdfds",
						"xyz": 2,
					},
				},
				sub: map[string]any{
					"abc": 1,
					"def": map[string]any{
						"abc": map[string]any{
							"abc": 1,
						},
					},
				},
			},
			want: false,
		},
		{
			name: "mix type in sub",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
					"def": map[string]any{
						"abc": []any{
							"abc",
							"xyz",
							"xyz2",
						},
					},
				},
				sub: map[string]any{
					"abc": 1,
					"def": map[string]any{
						"abc": []any{
							"abc",
							"xyz2",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "mix type in sub false",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
					"def": map[string]any{
						"abc": []any{
							"abc",
							"xyz",
							"xyz2",
						},
					},
				},
				sub: map[string]any{
					"abc": 1,
					"def": map[string]any{
						"abc": []any{
							"abc",
							"xyz3",
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMapSubset(tt.args.m, tt.args.sub); (got == nil) != tt.want {
				t.Errorf("IsMapSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
