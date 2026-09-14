package mapx_test

import (
	"fmt"

	"github.com/rakunlabs/mapx"
)

func ExampleGet() {
	config := map[string]any{
		"database": map[string]any{
			"host":     "localhost",
			"password": nil,
		},
	}

	host, ok := mapx.Get(config, []string{"database", "host"})
	fmt.Println(host, ok)

	password, ok := mapx.Get(config, []string{"database", "password"})
	fmt.Println(password, ok)

	missing, ok := mapx.Get(config, []string{"database", "missing"})
	fmt.Println(missing, ok)

	// Output:
	// localhost true
	// <nil> true
	// <nil> false
}

func ExampleMerge() {
	destination := map[string]any{
		"database": map[string]any{"host": "localhost", "port": 5432},
		"tags":     []any{"default", "development"},
	}
	source := map[string]any{
		"database": map[string]any{"host": "db.internal"},
		"tags":     []any{"production"},
	}

	result := mapx.Merge(source, destination)
	database := result["database"].(map[string]any)
	fmt.Println(database["host"], database["port"])
	fmt.Println(result["tags"])

	// Merge updates the destination in place.
	host, _ := mapx.Get(destination, []string{"database", "host"})
	fmt.Println(host)

	// Output:
	// db.internal 5432
	// [production]
	// db.internal
}

func ExampleMergeAny() {
	// Maps merge recursively; other values replace the destination.
	result := mapx.MergeAny(
		map[string]any{"host": "db.internal"},
		map[string]any{"host": "localhost", "port": 5432},
	).(map[string]any)
	fmt.Println(result["host"], result["port"])
	fmt.Println(mapx.MergeAny([]any{"production"}, []any{"development"}))
	fmt.Println(mapx.MergeAny("new", "old"))

	// Output:
	// db.internal 5432
	// [production]
	// new
}

func ExampleIsSubset() {
	// IsSubset dispatches to map or slice subset comparison.
	fmt.Println(mapx.IsSubset(
		map[string]any{"name": "mapx", "language": "Go"},
		map[string]any{"name": "mapx"},
	) == nil)
	fmt.Println(mapx.IsSubset([]any{"Go", "Rust"}, []any{"Go"}) == nil)
	fmt.Println(mapx.IsSubset([]any{"Go", "Rust"}, []any{"Python"}) == nil)

	// Output:
	// true
	// true
	// false
}

func ExampleIsMapSubset() {
	set := map[string]any{
		"service": map[string]any{"name": "api", "port": 8080},
		"enabled": true,
	}
	subset := map[string]any{
		"service": map[string]any{"name": "api"},
	}

	fmt.Println(mapx.IsMapSubset(set, subset) == nil)
	fmt.Println(mapx.IsMapSubset(set, map[string]any{"enabled": false}) == nil)

	// Output:
	// true
	// false
}

func ExampleIsSliceSubset() {
	set := []any{"Go", "Rust", "Python"}

	// Order and duplicate counts do not affect subset membership.
	fmt.Println(mapx.IsSliceSubset(set, []any{"Python", "Go"}) == nil)
	fmt.Println(mapx.IsSliceSubset(set, []any{"Go", "Go"}) == nil)
	fmt.Println(mapx.IsSliceSubset(set, []any{"Java"}) == nil)

	// Output:
	// true
	// true
	// false
}

func ExampleIsSliceContains() {
	services := []any{
		map[string]any{"name": "api", "port": 8080},
		map[string]any{"name": "worker", "port": 9090},
	}

	// Map elements use subset matching: specifying every field is unnecessary.
	fmt.Println(mapx.IsSliceContains(services, map[string]any{"name": "api"}) == nil)
	fmt.Println(mapx.IsSliceContains(services, map[string]any{"name": "web"}) == nil)

	// Output:
	// true
	// false
}

func ExampleWithWeakType() {
	set := map[string]any{"port": 8080}
	subset := map[string]any{"port": "8080"}

	// Weak typing is enabled by default and compares text representations.
	fmt.Println(mapx.IsMapSubset(set, subset) == nil)
	fmt.Println(mapx.IsMapSubset(set, subset, mapx.WithWeakType(false)) == nil)

	// Numeric strings are not normalized.
	fmt.Println(mapx.IsSliceContains([]any{1}, "1.0", mapx.WithWeakType(true)) == nil)

	// Output:
	// true
	// false
	// false
}

func ExampleWithCaseInsensitive() {
	set := map[string]any{"Name": "MapX"}
	subset := map[string]any{"name": "mapx"}

	// The option applies to both map keys and string values.
	fmt.Println(mapx.IsMapSubset(set, subset) == nil)
	fmt.Println(mapx.IsMapSubset(set, subset, mapx.WithCaseInsensitive(true)) == nil)

	// Without an exact key, multiple case-insensitive candidates are ambiguous.
	ambiguous := map[string]any{"Name": "MapX", "NAME": "MapX"}
	err := mapx.IsMapSubset(ambiguous, subset, mapx.WithCaseInsensitive(true))
	fmt.Println(err)

	// Output:
	// false
	// true
	// ambiguous case-insensitive key "name"
}
