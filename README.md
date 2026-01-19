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

## Benchmark: vs github.com/google/btree

This is the benchmark result against popular but now archived
[github.com/google/btree](https://github.com/google/btree).  We
adopted some of its benchmark tests for `treemap-go` to compare the
two implementations.  The `github.com/google/btree` package was
configured with degree=16 to match the configuration of `treemap-go`.

```
                 │ btree-deg16.txt │             treemap.txt             │
                 │     sec/op      │   sec/op     vs base                │
InsertG-24             91.52n ± 8%   88.81n ± 3%        ~ (p=0.143 n=10)
SeekG-24               55.95n ± 0%   44.00n ± 0%  -21.35% (p=0.000 n=10)
DeleteInsertG-24       166.2n ± 1%   146.8n ± 0%  -11.65% (p=0.000 n=10)
DeleteG-24             94.52n ± 2%   86.41n ± 2%   -8.58% (p=0.000 n=10)
GetG-24                77.25n ± 1%   70.91n ± 1%   -8.21% (p=0.000 n=10)
AscendG-24            24.283µ ± 1%   9.942µ ± 1%  -59.06% (p=0.000 n=10)
DescendG-24            24.33µ ± 2%   11.36µ ± 9%  -53.31% (p=0.000 n=10)
geomean                448.8n        326.7n       -27.22%
```

The degree in treemap-go is fixed to 16, and it does not support
copy-on-write feature, so it is not a drop-in replacement for btree.
For use cases that do not require those features, treemap-go might
provide an advantage, especially for iteration-heavy workloads.

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
