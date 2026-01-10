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
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifyMapNode[Key comparable, Value any](
	t *testing.T, node node[Key, Value],
	maxKey Key, m *Map[Key, Value],
) {
	t.Helper()

	switch node := node.(type) {
	case *internalNode[Key, Value]:
		verifyClear(t, node.nodes[node.n:])
		verifyClear(t, node.keys[node.n:])
		assert.True(t, slices.IsSortedFunc(node.keys[:node.n],
			m.compare))
		assert.LessOrEqual(t, node.keys[node.n-1], maxKey)

		for i, tnode := range node.nodes {
			verifyMapNode(t, tnode, node.keys[i], m)
		}
	case *leafNode[Key, Value]:
		verifyClear(t, node.values[node.n:])
		verifyClear(t, node.keys[node.n:])
		assert.True(t, slices.IsSortedFunc(node.keys[:node.n],
			m.compare))
		assert.LessOrEqual(t, node.keys[node.n-1], maxKey)
	}
}

func verifyClear[T any](t *testing.T, v []T) {
	t.Helper()

	for _, a := range v {
		assert.Zero(t, a)
	}
}

func verifyMapLen[Key comparable, Value any](t *testing.T, m *Map[Key, Value]) {
	t.Helper()

	sum := 0

	for tnode := m.front; tnode != nil; tnode = tnode.next {
		sum += tnode.n
	}

	assert.Equal(t, sum, m.Len())
}

func verifyMapKey[Key comparable, Value any](
	t *testing.T, m *Map[Key, Value], minKey Key,
) {
	t.Helper()

	key := minKey

	for k := range m.Begin().Seq() {
		assert.LessOrEqual(t, key, k)
		key = k
	}
}

func verifyMap[Key comparable, Value any](
	t *testing.T, m *Map[Key, Value], minKey, maxKey Key,
) {
	t.Helper()

	if m.root.Size() == 0 {
		assert.Equal(t, 0, m.Len())

		return
	}

	verifyMapNode(t, m.root, maxKey, m)
	verifyMapLen(t, m)
	verifyMapKey(t, m, minKey)
}

func printMap[Key comparable, Value any](m *Map[Key, Value]) { //nolint:unused
	printMapNode(m.root, 0)
}

func printMapNode[Key comparable, Value any]( //nolint:unused
	node node[Key, Value], level int,
) {
	fmt.Printf("lv=%d\n", level)

	switch n := node.(type) {
	case *internalNode[Key, Value]:
		fmt.Printf("len=%d %v\n", n.n, n.keys[:n.n])

		for _, n := range n.nodes[:n.n] {
			printMapNode(n, level+1)
		}
	case *leafNode[Key, Value]:
		fmt.Printf("len=%d %v\n", n.n, n.keys[:n.n])
	}
}

func TestMapInsert(t *testing.T) {
	m := New[string, int]()

	it := m.Insert("foo", 1)

	require.False(t, it.End())
	assert.Equal(t, "foo", it.Key())
	assert.Equal(t, 1, it.Value())

	it = it.Next()

	require.True(t, it.End())

	it = m.Insert("bar", 2)

	require.False(t, it.End())
	assert.Equal(t, "bar", it.Key())
	assert.Equal(t, 2, it.Value())

	it = it.Next()

	require.False(t, it.End())
	assert.Equal(t, "foo", it.Key())
	require.True(t, it.Next().End())

	verifyMap(t, m, "bar", "foo")

	assert.Equal(t, []string{"bar", "foo"},
		slices.Collect(m.Keys()))
	assert.Equal(t, []int{2, 1},
		slices.Collect(m.Values()))

	it = m.Insert("foo", 100)

	require.False(t, it.End())
	assert.Equal(t, "foo", it.Key())
	assert.Equal(t, 100, it.Value())
	assert.Equal(t, []int{2, 100},
		slices.Collect(m.Values()))
}

func TestMapInsertSplitNode(t *testing.T) {
	m := New[int, int]()

	for i := range 48 {
		m.Insert(i, i)
	}

	// Select right node after split
	it := m.Insert(32, 99)

	assert.Equal(t, 99, it.Value())
}

func TestMapInsertSplitMiddle(t *testing.T) {
	m := New[int, int]()

	for i := range 16 {
		m.Insert(i, i)
	}

	for i := range 32 {
		m.Insert(i+32, i+32)
	}

	m.Remove(15)

	for i := range 10 {
		m.Insert(i+15, i+15)
	}

	var keys []int //nolint:prealloc

	for i := range 15 {
		keys = append(keys, i)
	}

	for i := range 10 {
		keys = append(keys, i+15)
	}

	for i := range 32 {
		keys = append(keys, i+32)
	}

	assert.Equal(t, keys, slices.Collect(m.Keys()))
}

func genIntSeq(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				return
			}
		}
	}
}

func genIntSeqStep(start, n, step int) iter.Seq[int] { //nolint:unparam
	return func(yield func(int) bool) {
		for i := start; i < n; i += step {
			if !yield(i) {
				return
			}
		}
	}
}

func TestMapInsert1000(t *testing.T) {
	m := New[int, int]()

	for i := range 1000 {
		it := m.Insert(i, i+1)

		require.False(t, it.End())
		assert.Equal(t, i, it.Key())
		assert.Equal(t, i+1, it.Value())

		verifyMap(t, m, 0, i)
	}

	assert.Equal(t, slices.Collect(genIntSeq(1000)),
		slices.Collect(m.Keys()))
	assert.Equal(t, slices.Collect(genIntSeq(1001))[1:],
		slices.Collect(m.Values()))
}

func TestMapRemove(t *testing.T) {
	m := New[int, int]()

	assert.False(t, m.Remove(0))

	m.Insert(912, 1)
	m.Insert(313, 2)
	m.Insert(78, 3)

	assert.Equal(t, 3, m.Len())
	assert.True(t, m.Remove(912))
	assert.Equal(t, 2, m.Len())
	assert.False(t, m.Remove(912))
	assert.Equal(t, []int{78, 313}, slices.Collect(m.Keys()))
	assert.Equal(t, []int{3, 2}, slices.Collect(m.Values()))

	verifyMap(t, m, 78, 912)

	m = New[int, int]()

	for i := range 1000 {
		m.Insert(i, i)
	}

	assert.False(t, m.Remove(1000))
}

func TestMapRemoveIter(t *testing.T) {
	m := New[int, int]()

	it := m.RemoveIter(m.End())

	assert.True(t, it.End())

	m.Insert(912, 1)
	m.Insert(78, 3)
	it = m.Insert(313, 2)

	assert.Equal(t, 3, m.Len())

	it = m.RemoveIter(it)

	assert.Equal(t, 2, m.Len())
	require.False(t, it.End())
	assert.Equal(t, 912, it.Key())
	assert.Equal(t, 1, it.Value())
	assert.Equal(t, []int{78, 912}, slices.Collect(m.Keys()))
	assert.Equal(t, []int{3, 1}, slices.Collect(m.Values()))

	verifyMap(t, m, 78, 912)
}

func TestMapRemove1000(t *testing.T) {
	m := New[int, int]()

	for i := range 1000 {
		m.Insert(i, i+1)
	}

	for i := 0; i < 1000; i += 2 {
		assert.True(t, m.Remove(i))
		assert.Equal(t, 999-i/2, m.Len())
	}

	assert.Equal(t, slices.Collect(genIntSeqStep(1, 1000, 2)),
		slices.Collect(m.Keys()))
	assert.Equal(t, slices.Collect(genIntSeqStep(2, 1001, 2)),
		slices.Collect(m.Values()))

	verifyMap(t, m, 1, 999)
}

func TestMapRemoveIter1000(t *testing.T) {
	m := New[int, int]()

	for i := range 1000 {
		m.Insert(i, i+1)
	}

	it := m.Begin()
	for !it.End() {
		it = m.RemoveIter(it)

		require.False(t, it.End())

		it = it.Next()
	}

	assert.Equal(t, 500, m.Len())
	assert.Equal(t, slices.Collect(genIntSeqStep(1, 1000, 2)),
		slices.Collect(m.Keys()))
	assert.Equal(t, slices.Collect(genIntSeqStep(2, 1001, 2)),
		slices.Collect(m.Values()))

	verifyMap(t, m, 1, 999)
}

func TestMapRemoveIterNextNode(t *testing.T) {
	m := New[int, int]()

	for i := range 48 {
		m.Insert(i, i+1)
	}

	m.Remove(15)

	it := m.LowerBound(23)

	require.False(t, it.End())
	assert.Equal(t, 23, it.Key())

	it = m.RemoveIter(it)

	require.False(t, it.End())
	assert.Equal(t, 24, it.Key())
}

func TestMapFind(t *testing.T) {
	m := New[int, int]()

	m.Insert(98, 1)
	m.Insert(96, 2)
	m.Insert(99, 3)

	v, ok := m.Find(96)

	require.True(t, ok)
	assert.Equal(t, 2, v)

	v, ok = m.Find(98)

	require.True(t, ok)
	assert.Equal(t, 1, v)

	v, ok = m.Find(99)

	require.True(t, ok)
	assert.Equal(t, 3, v)

	_, ok = m.Find(100)

	assert.False(t, ok)

	_, ok = m.Find(97)

	assert.False(t, ok)

	_, ok = m.Find(95)

	assert.False(t, ok)
}

func TestMapFind1000(t *testing.T) {
	m := New[int, int]()

	for i := range 1000 {
		m.Insert(i, i+1)
	}

	v, ok := m.Find(111)

	assert.True(t, ok)
	assert.Equal(t, 112, v)
}

func TestMapFindAfterRemove(t *testing.T) {
	m := New[int, int]()

	for i := range 1000 {
		m.Insert(i, i+1)
	}

	m.Remove(15)
	m.Remove(511)

	for i := 495; i <= 510; i++ {
		m.Remove(i)
	}

	_, ok := m.Find(511)

	assert.False(t, ok)

	it := m.LowerBound(511)

	require.False(t, it.End())
	assert.Equal(t, 512, it.Key())
}

func TestMapLowerBound(t *testing.T) {
	m := New[int, int]()

	m.Insert(98, 1)
	m.Insert(96, 2)
	m.Insert(99, 3)

	it := m.LowerBound(95)

	require.False(t, it.End())
	assert.Equal(t, 96, it.Key())

	it = m.LowerBound(96)

	require.False(t, it.End())
	assert.Equal(t, 96, it.Key())

	it = m.LowerBound(97)

	require.False(t, it.End())
	assert.Equal(t, 98, it.Key())

	it = m.LowerBound(99)

	require.False(t, it.End())
	assert.Equal(t, 99, it.Key())

	it = m.LowerBound(100)

	require.True(t, it.End())
}

func TestLowerBoundNextNode(t *testing.T) {
	m := New[int, int]()

	for i := range 48 {
		m.Insert(i, i)
	}

	m.Remove(15)
	m.Remove(23)

	it := m.LowerBound(23)

	require.False(t, it.End())
	assert.Equal(t, 24, it.Key())

	m = New[int, int]()

	for i := range 1000 {
		m.Insert(i, i)
	}

	it = m.LowerBound(1000)

	require.True(t, it.End())
}

func TestMapKeys(t *testing.T) {
	m := New[int, int]()

	keys := []int{3, 7, 9}

	for _, k := range keys {
		m.Insert(k, k)
	}

	assert.Equal(t, keys, slices.Collect(m.Keys()))

	for range m.Keys() {
		break
	}
}

func TestMapValues(t *testing.T) {
	m := New[int, int]()

	keys := []int{3, 7, 9}
	values := []int{8, 1, 5}

	for i, k := range keys {
		m.Insert(k, values[i])
	}

	assert.Equal(t, values, slices.Collect(m.Values()))

	for range m.Values() {
		break
	}
}

func TestMapComparable(t *testing.T) {
	m := NewComparable[string, int](cmp.Compare[string])

	m.Insert("foo", 1)
	m.Insert("bar", 2)

	assert.Equal(t, []string{"bar", "foo"}, slices.Collect(m.Keys()))
}

func TestMapString(t *testing.T) {
	m := New[int, string]()

	assert.Equal(t, "Map[]", m.String())

	m.Insert(1, "foo")
	m.Insert(2, "bar")

	assert.Equal(t, "Map[1:foo 2:bar]", m.String())
}
