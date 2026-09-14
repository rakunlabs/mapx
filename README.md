# mapX

[![License](https://img.shields.io/github/license/rakunlabs/mapx?color=red&style=flat-square)](https://raw.githubusercontent.com/rakunlabs/mapx/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/rakunlabs_mapx?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=rakunlabs_mapx)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/rakunlabs/mapx/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/rakunlabs/mapx/actions)
[![Go PKG](https://raw.githubusercontent.com/rakunlabs/.github/main/assets/badges/gopkg.svg)](https://pkg.go.dev/github.com/rakunlabs/mapx)

Map functions to help compare and merge of maps.

```sh
go get github.com/rakunlabs/mapx
```

## Usage

Options:
- `mapx.WithWeakType(true/false)` to enable/disable weak type comparison. Default is `true`.
- `mapx.WithCaseInsensitive(true/false)` to enable/disable case insensitive map key and string value comparison. Default is `false`.

```go
testMap := map[string]any{
    "abc": 1,
    "def": map[string]any{
        "abc": int64(1),
        "xyz": float32(2),
    },
}

subMap := map[string]any{
    "abc": 1,
    "def": map[string]any{
        "abc": "1",
    },
}

err := mapx.IsMapSubset(testMap, subMap, mapx.WithWeakType(true))
// this will return nil error because 1 == "1" with weak type comparison
```

### Comparison rules

- Comparison functions return `nil` for a match and an error otherwise.
- `IsSubset` accepts two `map[string]any` values or two `[]any` values. An untyped `nil` top-level set is rejected.
- Nested collections must also use `map[string]any` and `[]any`. Other collection types, such as `map[string]int` and `[]string`, do not match and return errors rather than panicking.
- A present `nil` value matches another `nil` value. A missing map key is different from a present key with a `nil` value. A typed nil (such as `(*int)(nil)`) is different from an untyped nil.
- Nil maps and slices behave as empty collections. An empty collection is a subset of any supported collection of the same kind.
- Slice comparisons ignore order and duplicate counts: `[]any{1, 1}` is a subset of `[]any{1}`. Nested maps and slices also use subset matching, not exact equality.
- Weak typing compares different comparable types using their Go `%v` text representations. For example, `1`, `int64(1)` and `"1"` match, but `1` and `"1.0"` do not. Numeric strings are not parsed or normalized. Disabling weak typing requires identical Go types, including for named string types.
- Case-insensitive key lookup prefers an exact key. Without an exact key, multiple case-insensitive candidates produce an ambiguity error. For example, looking up `"name"` in a map containing both `"Name"` and `"NAME"` is ambiguous.
- Values that cannot be safely compared, including structs containing slice-valued interface fields, do not match.
- Recursive input collections must not contain cycles.

### Merge

```go
destination := map[string]any{"database": map[string]any{"host": "localhost", "port": 5432}}
source := map[string]any{"database": map[string]any{"host": "db.internal"}}
result := mapx.Merge(source, destination)
// destination and result now contain host=db.internal and port=5432.
```

`Merge(source, destination)` modifies and returns the destination. Source values
take precedence. Nested values are merged recursively only when both are
`map[string]any`; other values, including slices, replace the destination value.
If the destination is nil, a new top-level map is allocated.

This is not a deep copy: assigned nested maps and slices share storage with the
source. Changes to those shared values can affect both maps. Inputs must not
contain cycles.

`MergeAny(source, destination)` calls `Merge` when both arguments are
`map[string]any`; otherwise it returns the source directly.

### Get

```go
value, ok := mapx.Get(result, []string{"database", "host"})
// value == "db.internal", ok == true
```

`Get` traverses `map[string]any` values using a key path. An empty path, a missing
key, or a non-map intermediate value returns `(nil, false)`. A present leaf with
a nil value returns `(nil, true)`. Slice indexing is not supported.
