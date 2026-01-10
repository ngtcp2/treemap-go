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
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorNext(t *testing.T) {
	m := New[int, int]()

	for i := range 100 {
		m.Insert(i, i+1)
	}

	it := m.Begin()

	key := 0

	for !it.End() {
		assert.Equal(t, key, it.Key())
		assert.Equal(t, key+1, it.Value())

		key++
		it = it.Next()
	}

	assert.Equal(t, 100, key)
}

func TestIteratorPrev(t *testing.T) {
	m := New[int, int]()

	for i := range 100 {
		m.Insert(i, i+1)
	}

	it := m.End()

	key := 99

	for !it.Begin() {
		it = it.Prev()

		assert.Equal(t, key, it.Key())
		assert.Equal(t, key+1, it.Value())

		key--
	}

	assert.Equal(t, -1, key)
}

type Item[Key comparable, Value any] struct {
	Key   Key
	Value Value
}

func Collect[Key comparable, Value any](
	seq iter.Seq2[Key, Value],
) []Item[Key, Value] {
	var s []Item[Key, Value]

	for k, v := range seq {
		s = append(s, Item[Key, Value]{Key: k, Value: v})
	}

	return s
}

func TestIteratorSeq(t *testing.T) {
	m := New[int, string]()

	items := []Item[int, string]{
		{
			Key:   11,
			Value: "foo",
		},
		{
			Key:   98,
			Value: "bar",
		},
		{
			Key:   129,
			Value: "baz",
		},
	}

	for _, item := range items {
		m.Insert(item.Key, item.Value)
	}

	assert.Equal(t, items, Collect(m.Begin().Seq()))

	for range m.Begin().Seq() {
		break
	}
}

func TestIteratorSetValue(t *testing.T) {
	m := New[int, string]()

	m.Insert(0, "foo")
	m.Insert(1, "bar")
	m.Insert(2, "baz")

	it := m.Begin()
	for !it.End() {
		if it.Key() == 1 {
			it.SetValue("BAR")
		}

		it = it.Next()
	}

	assert.Equal(t, []string{"foo", "BAR", "baz"},
		slices.Collect(m.Values()))
}
