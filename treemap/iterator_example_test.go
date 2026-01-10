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
	"fmt"
)

func ExampleIterator() {
	m := New[int, string]()

	m.Insert(0, "foo")
	m.Insert(1, "bar")

	it := m.Begin()
	for !it.End() {
		fmt.Println(it.Key(), it.Value())
		it = it.Next()
	}
	// Output:
	// 0 foo
	// 1 bar
}

func ExampleIterator_Seq() {
	m := New[int, string]()

	m.Insert(0, "foo")
	m.Insert(1, "bar")

	for key, value := range m.Begin().Seq() {
		fmt.Println(key, value)
	}
	// Output:
	// 0 foo
	// 1 bar
}

func ExampleIterator_Next() {
	m := New[int, string]()

	m.Insert(0, "foo")
	m.Insert(1, "bar")

	it := m.Begin()
	for !it.End() {
		fmt.Println(it.Key(), it.Value())
		it = it.Next()
	}
	// Output:
	// 0 foo
	// 1 bar
}

func ExampleIterator_Prev() {
	m := New[int, string]()

	m.Insert(0, "foo")
	m.Insert(1, "bar")

	it := m.End()
	for !it.Begin() {
		it = it.Prev()
		fmt.Println(it.Key(), it.Value())
	}
	// Output:
	// 1 bar
	// 0 foo
}
