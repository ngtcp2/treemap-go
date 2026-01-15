# treemap-go

The package `github.com/ngtcp2/treemap-go/treemap` provides Map, the
sorted, key-value storage in Go.  The underlying implementation is
based on B+ Tree.  It aims to generally efficient
insertion/removal/iteration, and fewer memory allocations.

The tree implementation was originally developed in
[ngtcp2](https://github.com/ngtcp2/ngtcp2) project.

See https://pkg.go.dev/github.com/ngtcp2/treemap-go/treemap for
documentation and how to use this package.

For benchmarks, see [bench](treemap/bench) directory.

## Example

```go
package main

import (
	"fmt"

	"github.com/ngtcp2/treemap-go/treemap"
)

func main() {
	m := treemap.New[int, string]()

	// Insertion
	m.Insert(1, "foo")
	m.Insert(2, "bar")

	// Lookup
	fmt.Println(m.Find(1))
	fmt.Println(m.Find(2))

	// Iteration
	it := m.Begin()

	for !it.End() {
		fmt.Println(it.Key(), it.Value())
		it = it.Next()
	}

	// Iteration with Go iterator
	for k, v := range m.Begin().Seq() {
		fmt.Println(k, v)
	}

	// Removal
	m.Remove(1)
	m.Remove(2)
}
```

## License

```
The MIT License

Copyright (c) 2026 treemap-go contributors

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
