# orderedmap

A thread-safe, ordered map implementation for Go with generics support.

## Features

- Thread-safe with read-write mutex
- Maintains insertion order
- Supports Go generics
- Simple and easy-to-use API
- JSON support
- Go v1.21+ required
- MIT License

## Installation

```bash
go get github.com/zoppax/orderedmap
```

## API

### Types

- `Pair[K, V]` - Key-value pair with `Key` and `Value` fields

### Constructor

- `New[K comparable, V any]()` - Creates a new empty OrderedMap
- `FromMap[K comparable, V any](m map[K]V, keys []K)` - Creates an OrderedMap from a map with specified key order

### Basic Operations

- `Set(key K, value V)` - Sets a key-value pair. If key exists, updates the value
- `Get(key K) (V, bool)` - Gets value by key, returns (value, true) if exists
- `Delete(key K)` - Deletes a key-value pair

### Query

- `Len() int` - Returns the number of elements
- `Index(key K) int` - Returns the index of key, -1 if not found
- `Keys() []K` - Returns all keys in insertion order
- `Values() []V` - Returns all values in insertion order
- `First() *Pair[K, V]` - Returns the first key-value pair, nil if empty
- `Last() *Pair[K, V]` - Returns the last key-value pair, nil if empty

### Iteration

- `Iterate(fn func(*Pair[K, V]) bool)` - Iterates over all pairs. Return false to stop iteration

### JSON Support

`OrderedMap` implements `json.Marshaler` and `json.Unmarshaler` interfaces, so you can use it directly with the standard library's `encoding/json` package.

### Error Types

- `ErrInvalidJSON` - Returned when JSON format is invalid
- `ErrUnsupportedKeyType` - Returned when key type is not string

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/zoppax/orderedmap"
)

func main() {
	// Create from empty map
	om := orderedmap.New[string, int]()

	om.Set("b", 2)
	om.Set("a", 1)
	om.Set("c", 3)

	fmt.Println("Keys:", om.Keys())   // [b a c]
	fmt.Println("Values:", om.Values()) // [2 1 3]

	val, ok := om.Get("a")
	fmt.Printf("Get a: %d, ok: %v\n", val, ok)

	fmt.Println("Len:", om.Len()) // 3
	fmt.Println("Index of 'c':", om.Index("c")) // 2
	fmt.Println("Index of 'z':", om.Index("z")) // -1, not found

	first := om.First()
	fmt.Printf("First: %s=%d\n", first.Key, first.Value)

	last := om.Last()
	fmt.Printf("Last: %s=%d\n", last.Key, last.Value)

	om.Iterate(func(pair *orderedmap.Pair[string, int]) bool {
		fmt.Printf("%s: %d\n", pair.Key, pair.Value)
		return true
	})

	om.Delete("b")
	fmt.Println("After delete, Keys:", om.Keys()) // [a c]

	// Or create from a map with specified order
	m := map[string]int{"x": 10, "y": 20, "z": 30}
	keys := []string{"z", "x", "y"}
	om2 := orderedmap.FromMap(m, keys)
	fmt.Println("Keys:", om2.Keys())   // [z x y]
	fmt.Println("Values:", om2.Values()) // [30 10 20]

	// JSON Marshal
	data, err := json.Marshal(om)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON:", string(data))

	// JSON Unmarshal
	var om3 orderedmap.OrderedMap[string, int]
	if err := json.Unmarshal(data, &om3); err != nil {
		panic(err)
	}
	fmt.Println("Keys after unmarshal:", om3.Keys()) // [b a c]
}
```

## Authors

- [Macc Liu](https://maccliu.com)

## License

MIT License | Copyright (c) 2026 Zoppax LLC
