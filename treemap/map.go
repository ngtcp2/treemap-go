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
	"strings"
)

// Compare is the function to compare x and y.  If x is less than y,
// it must return -1.  If y is less than x, it must return 1.
// Otherwise, x and y are considered equal, and this function must
// return 0.
type Compare[Key any] func(x, y Key) int

type search[Key any] func([]Key, Key) (int, bool)

// Map is the sorted, key-value storage.
type Map[Key, Value any] struct {
	root    node[Key, Value]
	front   *leafNode[Key, Value]
	back    *leafNode[Key, Value]
	n       int
	compare func(lhs, rhs Key) int
	search  search[Key]
}

// New returns new Map for the ordered keys.
func New[Key cmp.Ordered, Value any]() *Map[Key, Value] {
	node := &leafNode[Key, Value]{}

	return &Map[Key, Value]{
		root:    node,
		front:   node,
		back:    node,
		compare: cmp.Compare[Key],
		search:  linearSearchOrdered[Key],
	}
}

// linearSearchOrdered searches target in keys in O(n).  With keyDegr
// = 16, linear search is faster than binary search for cmp.Ordered
// types.
func linearSearchOrdered[Key cmp.Ordered](keys []Key, target Key) (int, bool) {
	for i, key := range keys {
		switch {
		case key == target:
			return i, true
		case key > target:
			return i, false
		}
	}

	return len(keys), false
}

// NewAny returns new Map with custom [Compare] function.  [New]
// should be used for Key that is of type [cmp.Ordered] because it is
// much more efficient.
func NewAny[Key, Value any](
	compare Compare[Key],
) *Map[Key, Value] {
	node := &leafNode[Key, Value]{}

	return &Map[Key, Value]{
		root:    node,
		front:   node,
		back:    node,
		compare: compare,
		search: func(keys []Key, target Key) (int, bool) {
			return slices.BinarySearchFunc(keys, target, compare)
		},
	}
}

func (m *Map[Key, Value]) splitRoot() {
	rnode := m.root.Split(m)
	lnode := m.root

	root := &internalNode[Key, Value]{
		n: 2,
	}

	root.keys[0] = lnode.LastKey()
	root.nodes[0] = lnode

	root.keys[1] = rnode.LastKey()
	root.nodes[1] = rnode

	m.root = root
}

// Insert inserts the given key-value pair.  If the key already
// exists, its value is replaced with the given value.  It returns the
// Iterator that points to the inserted or updated item.  If the
// existing value is replaced with new value, this function returns
// the old value and true.  Otherwise, zero value and false.
func (m *Map[Key, Value]) Insert(
	key Key, value Value,
) (Iterator[Key, Value], Value, bool) {
	if m.root.IsFull() {
		m.splitRoot()
	}

	node := m.root

	for {
		if tnode, ok := node.(*leafNode[Key, Value]); ok {
			var oldValue Value

			i, ok := m.search(tnode.Keys(), key)
			if ok {
				oldValue = tnode.values[i]
				tnode.values[i] = value
			} else {
				tnode.InsertAt(i, key, value)

				m.n++
			}

			return Iterator[Key, Value]{
				node: tnode,
				idx:  i,
			}, oldValue, ok
		}

		inode := node.(*internalNode[Key, Value])

		i, _ := m.search(inode.Keys(), key)
		if i == inode.n {
			for {
				node = inode.nodes[inode.n-1]
				if node.IsFull() {
					inode.SplitAt(inode.n-1, m)
					node = inode.nodes[inode.n-1]
				}

				inode.keys[inode.n-1] = key

				var ok bool

				inode, ok = node.(*internalNode[Key, Value])
				if !ok {
					break
				}
			}

			tnode := node.(*leafNode[Key, Value])
			idx := tnode.n
			tnode.InsertAt(idx, key, value)

			m.n++

			var oldValue Value

			return Iterator[Key, Value]{
				node: tnode,
				idx:  idx,
			}, oldValue, false
		}

		descNode := inode.nodes[i]

		if descNode.IsFull() {
			inode.SplitAt(i, m)

			if m.compare(inode.keys[i], key) < 0 {
				descNode = inode.nodes[i+1]
			}
		}

		node = descNode
	}
}

// Find returns value associated by key.  If such value exists, the
// value and true are returned.  Otherwise, zero value and false are
// returned.
func (m *Map[Key, Value]) Find(key Key) (Value, bool) {
	var z Value

	node := m.root

	for {
		if tnode, ok := node.(*leafNode[Key, Value]); ok {
			i, ok := m.search(tnode.Keys(), key)
			if !ok {
				return z, false
			}

			return tnode.values[i], true
		}

		inode := node.(*internalNode[Key, Value])

		i, _ := m.search(inode.KeysForFindAndRemove(), key)
		node = inode.nodes[i]
	}
}

// LowerBound returns the Iterator that points to the item whose key
// is the smallest key that is greater than or equal to key.  If all
// stored keys are smaller than key, it returns the Iterator whose
// [Iterator.End] returns true.
func (m *Map[Key, Value]) LowerBound(key Key) Iterator[Key, Value] {
	node := m.root

	for {
		if tnode, ok := node.(*leafNode[Key, Value]); ok {
			i, _ := m.search(tnode.Keys(), key)
			if i == tnode.n && tnode.next != nil {
				tnode = tnode.next
				i = 0
			}

			return Iterator[Key, Value]{
				node: tnode,
				idx:  i,
			}
		}

		inode := node.(*internalNode[Key, Value])

		i, _ := m.search(inode.KeysForFindAndRemove(), key)
		node = inode.nodes[i]
	}
}

func (m *Map[Key, Value]) mergeNode(
	node *internalNode[Key, Value], i int,
) node[Key, Value] {
	lnode := node.nodes[i]
	rnode := node.nodes[i+1]

	lnode.Merge(rnode, m)

	if m.root == node && node.n == 2 {
		m.root = lnode
	} else {
		node.RemoveAt(i + 1)
		node.keys[i] = lnode.LastKey()
	}

	return lnode
}

func (m *Map[Key, Value]) shiftLeft(node *internalNode[Key, Value], i int) {
	lnode := node.nodes[i-1]
	rnode := node.nodes[i]

	n := (lnode.Size()+rnode.Size()+1)/2 - lnode.Size()

	lnode.ShiftLeft(rnode, n)
	node.keys[i-1] = lnode.LastKey()
}

func (m *Map[Key, Value]) shiftRight(node *internalNode[Key, Value], i int) {
	lnode := node.nodes[i]
	rnode := node.nodes[i+1]

	n := (lnode.Size()+rnode.Size()+1)/2 - rnode.Size()

	lnode.ShiftRight(rnode, n)
	node.keys[i] = lnode.LastKey()
}

// Remove removes the item identified by key.  If an item is removed,
// it returns the removed value and true.  Otherwise, returns zero
// value and false.
func (m *Map[Key, Value]) Remove(key Key) (Value, bool) {
	_, oldValue, ok := m.remove(key)

	return oldValue, ok
}

// RemoveIter removes the item pointed by it.  It returns the Iterator
// that points to the item that follows the removed item.  The
// provided it must not be invalidated, that means this function
// always successfully remove the item.  The one exception is the case
// where [Iterator.End] returns true.  In this case, this function
// returns it without doing anything.
func (m *Map[Key, Value]) RemoveIter(
	it Iterator[Key, Value],
) Iterator[Key, Value] {
	if it.End() {
		return it
	}

	tnode := it.node

	if tnode != m.root && tnode.n == minNodes {
		it, _, _ := m.remove(it.Key())
		return it
	}

	tnode.RemoveAt(it.idx)

	m.n--

	if tnode.n == it.idx && tnode.next != nil {
		return Iterator[Key, Value]{
			node: tnode.next,
		}
	}

	return Iterator[Key, Value]{
		node: tnode,
		idx:  it.idx,
	}
}

func (m *Map[Key, Value]) remove(key Key) (Iterator[Key, Value], Value, bool) {
	node := m.root

	if inode, ok := node.(*internalNode[Key, Value]); ok {
		if inode.n == 2 &&
			inode.nodes[0].Size() == minNodes &&
			inode.nodes[1].Size() == minNodes {
			node = m.mergeNode(inode, 0)
		}
	}

	for {
		if tnode, ok := node.(*leafNode[Key, Value]); ok {
			var oldValue Value

			i, _ := m.search(tnode.Keys(), key)
			if i == tnode.n || m.compare(key, tnode.keys[i]) != 0 {
				return m.End(), oldValue, false
			}

			oldValue = tnode.values[i]
			tnode.RemoveAt(i)

			m.n--

			if tnode.n == i && tnode.next != nil {
				return Iterator[Key, Value]{
					node: tnode.next,
				}, oldValue, true
			}

			return Iterator[Key, Value]{
				node: tnode,
				idx:  i,
			}, oldValue, true
		}

		inode := node.(*internalNode[Key, Value])

		i, _ := m.search(inode.KeysForFindAndRemove(), key)
		descNode := inode.nodes[i]

		if descNode.Size() > minNodes {
			node = descNode
			continue
		}

		if i+1 < inode.n && inode.nodes[i+1].Size() > minNodes {
			m.shiftLeft(inode, i+1)

			node = descNode

			continue
		}

		if i > 0 && inode.nodes[i-1].Size() > minNodes {
			m.shiftRight(inode, i-1)

			node = descNode

			continue
		}

		if i+1 < inode.n {
			node = m.mergeNode(inode, i)
			continue
		}

		node = m.mergeNode(inode, i-1)
	}
}

// Begin returns the Iterator that points to the first item.
func (m *Map[Key, Value]) Begin() Iterator[Key, Value] {
	return Iterator[Key, Value]{
		node: m.front,
	}
}

// End returns the Iterator that points to the one beyond the last
// item.
func (m *Map[Key, Value]) End() Iterator[Key, Value] {
	return Iterator[Key, Value]{
		node: m.back,
		idx:  m.back.n,
	}
}

// Len returns the number of items that m contains.
func (m *Map[Key, Value]) Len() int {
	return m.n
}

// Keys returns an iterator over keys in m in the sorted order.
func (m *Map[Key, Value]) Keys() iter.Seq[Key] {
	return func(yield func(Key) bool) {
		for it := m.Begin(); !it.End(); it = it.Next() {
			if !yield(it.Key()) {
				return
			}
		}
	}
}

// Values returns an iterator over values in m in the sorted order of
// the corresponding keys.
func (m *Map[Key, Value]) Values() iter.Seq[Value] {
	return func(yield func(Value) bool) {
		for it := m.Begin(); !it.End(); it = it.Next() {
			if !yield(it.Value()) {
				return
			}
		}
	}
}

// String returns the string representation of m.
func (m *Map[Key, Value]) String() string {
	var b strings.Builder

	b.WriteString("Map[")

	it := m.Begin()
	if !it.End() {
		fmt.Fprintf(&b, "%v:%v", it.Key(), it.Value())
		it = it.Next()

		for k, v := range it.Seq() {
			fmt.Fprintf(&b, " %v:%v", k, v)
		}
	}

	b.WriteString("]")

	return b.String()
}

// Clear removes all items from m.
func (m *Map[Key, Value]) Clear() {
	if m.n == 0 {
		return
	}

	node := &leafNode[Key, Value]{}
	m.root = node
	m.front = node
	m.back = node
	m.n = 0
}
