# Compare

Go compare library.

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/compare.svg)](https://pkg.go.dev/github.com/pierrre/compare)

## Features

- Compares any two values, like `reflect.DeepEqual()`
- Returns a detailed diff result with paths, messages, and values
- Supports all Go types: structs, maps, slices, arrays, pointers, interfaces, channels, functions, and more
- Handles recursive data structures without infinite loops
- Supports custom comparison methods: `.Equal()`, `.Eq()`, `.Cmp()`, or your own
- Highly configurable: per-type comparators, filters, max depth, and max differences
- Optimized for performance

## Usage

- [Compare](https://pkg.go.dev/github.com/pierrre/compare#example-package)
