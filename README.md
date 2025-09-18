# mapX

[![License](https://img.shields.io/github/license/rakunlabs/mapx?color=red&style=flat-square)](https://raw.githubusercontent.com/rakunlabs/mapx/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/rakunlabs_mapx?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=rakunlabs_mapx)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/rakunlabs/mapx/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/rakunlabs/mapx/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/rakunlabs/mapx?style=flat-square)](https://goreportcard.com/report/github.com/rakunlabs/mapx)
[![Go PKG](https://raw.githubusercontent.com/rakunlabs/.github/main/assets/badges/gopkg.svg)](https://pkg.go.dev/github.com/rakunlabs/mapx)

Map functions to help compare and merge of maps.

```sh
go get github.com/rakunlabs/mapx
```

## Usage

Options:
- `mapx.WithWeakType(true/false)` to enable/disable weak type comparison. Default is `true`.
- `mapx.WithCaseInsensitive(true/false)` to enable/disable case insensitive key comparison. Default is `false`.

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
