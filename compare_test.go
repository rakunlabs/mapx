package mapx

import (
	"encoding/json"
	"testing"
)

func TestIsMapSubset(t *testing.T) {
	type args struct {
		m             map[string]any
		sub           map[string]any
		OptionCompare []OptionCompare
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
					"aBc": "1",
				},
				OptionCompare: []OptionCompare{WithCaseInsensitive(true)},
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
			name: "mix type string and number",
			args: args{
				m: map[string]any{
					"abc": 1,
					"xyz": 2,
					"def": map[string]any{
						"abc": 1,
						"xyz": int64(2),
					},
				},
				sub: map[string]any{
					"abc": 1,
					"def": map[string]any{
						"abc": json.Number("1"),
						"xyz": float32(2),
					},
				},
			},
			want: true,
		},
		{
			name: "mix type false",
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
				OptionCompare: []OptionCompare{WithWeakType(false)},
			},
			want: false,
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
			if got := IsMapSubset(tt.args.m, tt.args.sub, tt.args.OptionCompare...); (got == nil) != tt.want {
				t.Errorf("IsMapSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
