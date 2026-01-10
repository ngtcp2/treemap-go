// treemap-go
//
// Copyright (c) 2026 treemap-go contributors
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package treemap

import (
	"cmp"
	"fmt"
	"slices"
)

func ExampleMap() {
	m := New[string, int]()

	m.Insert("foo", 100)
	m.Insert("bar", 250)

	fmt.Println(m.Find("foo"))
	fmt.Println(m.Find("bar"))
	fmt.Println(m.Find("baz"))
	// Output:
	// 100 true
	// 250 true
	// 0 false
}

func ExampleNewComparable() {
	type Pair struct {
		First, Second string
	}

	m := NewComparable[Pair, int](func(x, y Pair) int {
		return cmp.Or(
			cmp.Compare(x.First, y.First),
			cmp.Compare(x.Second, y.Second),
		)
	})

	m.Insert(Pair{First: "foo", Second: "alpha"}, 1)
	m.Insert(Pair{First: "bar", Second: "bravo"}, 2)

	for k, v := range m.Begin().Seq() {
		fmt.Println(k, v)
	}
	// Output:
	// {bar bravo} 2
	// {foo alpha} 1
}

func ExampleMap_End() {
	m := New[int, int]()

	m.Insert(1, 100)
	m.Insert(2, 200)
	m.Insert(3, 300)

	// Iterate backwards
	it := m.End()

	for !it.Begin() {
		it = it.Prev()
		fmt.Println(it.Key(), it.Value())
	}
	// Output:
	// 3 300
	// 2 200
	// 1 100
}

func ExampleMap_Find() {
	m := New[string, int]()

	m.Insert("alpha", 1)
	m.Insert("bravo", 2)
	m.Insert("charlie", 3)

	fmt.Println(m.Find("bravo"))
	fmt.Println(m.Find("echo"))
	// Output:
	// 2 true
	// 0 false
}

func ExampleMap_LowerBound() {
	m := New[string, int]()

	m.Insert("alpha", 1)
	m.Insert("bravo", 2)
	m.Insert("delta", 4)
	m.Insert("echo", 5)

	it := m.LowerBound("bravo")
	fmt.Println(it.Key(), it.Value())
	it = m.LowerBound("charlie")
	fmt.Println(it.Key(), it.Value())
	// Output:
	// bravo 2
	// delta 4
}

func ExampleMap_Insert() {
	m := New[int, int]()

	m.Insert(0, 100)
	fmt.Println(m.Find(0))
	m.Insert(0, 200)
	fmt.Println(m.Find(0))
	// Output:
	// 100 true
	// 200 true
}

func ExampleMap_Remove() {
	m := New[int, string]()

	m.Insert(1, "alpha")
	m.Insert(2, "bravo")
	m.Insert(3, "charlie")

	m.Remove(2)

	fmt.Println(m.Find(2))
	// Output:
	//  false
}

func ExampleMap_RemoveIter() {
	m := New[int, string]()

	m.Insert(1, "alpha")
	m.Insert(2, "bravo")
	m.Insert(3, "charlie")

	it := m.LowerBound(2)
	it = m.RemoveIter(it)

	fmt.Println(it.Key(), it.Value())
	// Output:
	// 3 charlie
}

func ExampleMap_Keys() {
	m := New[int, int]()

	m.Insert(1, 100)
	m.Insert(2, 200)
	m.Insert(3, 300)

	fmt.Println(slices.Collect(m.Keys()))
	// Output:
	// [1 2 3]
}

func ExampleMap_Values() {
	m := New[int, int]()

	m.Insert(1, 100)
	m.Insert(2, 200)
	m.Insert(3, 300)

	fmt.Println(slices.Collect(m.Values()))
	// Output:
	// [100 200 300]
}
