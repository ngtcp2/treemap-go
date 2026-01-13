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
)

// Iterator points to the specific item and can iterate items in both
// direction.  Iterator is invalidated when there is a change in the
// underlying [Map].  In general, insertion and removal make all
// existing Iterators invalidated.
type Iterator[Key comparable, Value any] struct {
	node *leafNode[Key, Value]
	idx  int
}

// Key returns the key pointed by it.  This function must not be
// called if [Iterator.End] returns true.
func (it Iterator[Key, Value]) Key() Key {
	return it.node.keys[it.idx]
}

// Value returns the value pointed by it.  This function must not be
// called if [Iterator.End] returns true.
func (it Iterator[Key, Value]) Value() Value {
	return it.node.values[it.idx]
}

// SetValue sets value to the current position.  This function must
// not be called if [Iterator.End] returns true.
func (it Iterator[Key, Value]) SetValue(value Value) {
	it.node.values[it.idx] = value
}

// Begin returns true if it points to the first item.
func (it Iterator[Key, Value]) Begin() bool {
	return it.idx == 0 && it.node.prev == nil
}

// End returns true if it points to the one beyond the last item.
func (it Iterator[Key, Value]) End() bool {
	return it.node.n == it.idx && it.node.next == nil
}

// Next returns the Iterator that points to the next item.  This
// function must not be called if [Iterator.End] returns true.
func (it Iterator[Key, Value]) Next() Iterator[Key, Value] {
	it.idx++

	if it.idx == it.node.n && it.node.next != nil {
		it.node = it.node.next
		it.idx = 0
	}

	return it
}

// Prev returns the Iterator that points to the previous item.  This
// function must not be called if [Iterator.Begin] returns true.
func (it Iterator[Key, Value]) Prev() Iterator[Key, Value] {
	if it.idx == 0 {
		it.node = it.node.prev
		it.idx = it.node.n - 1
	} else {
		it.idx--
	}

	return it
}

// Seq returns Go iterator.
func (it Iterator[Key, Value]) Seq() iter.Seq2[Key, Value] {
	return func(yield func(Key, Value) bool) {
		for ; !it.End(); it = it.Next() {
			if !yield(it.Key(), it.Value()) {
				return
			}
		}
	}
}
