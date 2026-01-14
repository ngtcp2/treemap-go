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

const (
	keyDegr  = 16
	maxNodes = 2 * keyDegr
	minNodes = keyDegr
)

type node[Key comparable, Value any] interface {
	Size() int
	LastKey() Key
	IsFull() bool
	Split(*Map[Key, Value]) node[Key, Value]
	Merge(node[Key, Value], *Map[Key, Value])
	ShiftLeft(o node[Key, Value], n int)
	ShiftRight(o node[Key, Value], n int)
}

type internalNode[Key comparable, Value any] struct {
	nodes [maxNodes]node[Key, Value]
	keys  [maxNodes]Key
	n     int
}

func (inode *internalNode[Key, Value]) Keys() []Key {
	return inode.keys[:inode.n]
}

func (inode *internalNode[Key, Value]) KeysForFindAndRemove() []Key {
	return inode.keys[:inode.n-1]
}

func (inode *internalNode[Key, Value]) Size() int {
	return inode.n
}

func (inode *internalNode[Key, Value]) LastKey() Key {
	return inode.keys[inode.n-1]
}

func (inode *internalNode[Key, Value]) IsFull() bool {
	return inode.n == maxNodes
}

func (inode *internalNode[Key, Value]) Split(
	*Map[Key, Value],
) node[Key, Value] {
	rnode := &internalNode[Key, Value]{}

	n := inode.n
	rnode.n = inode.n / 2
	inode.n -= rnode.n

	copy(rnode.nodes[:], inode.nodes[inode.n:n])
	clear(inode.nodes[inode.n:n])
	copy(rnode.keys[:], inode.keys[inode.n:n])
	clear(inode.keys[inode.n:n])

	return rnode
}

func (inode *internalNode[Key, Value]) SplitAt(idx int, m *Map[Key, Value]) {
	lnode := inode.nodes[idx]
	rnode := lnode.Split(m)

	copy(inode.nodes[idx+2:], inode.nodes[idx+1:inode.n])
	copy(inode.keys[idx+1:], inode.keys[idx:inode.n])

	inode.nodes[idx+1] = rnode
	inode.n++

	inode.keys[idx] = lnode.LastKey()
}

func (inode *internalNode[Key, Value]) Merge(
	o node[Key, Value], _ *Map[Key, Value],
) {
	rnode := o.(*internalNode[Key, Value])

	copy(inode.nodes[inode.n:], rnode.nodes[:rnode.n])
	copy(inode.keys[inode.n:], rnode.keys[:rnode.n])

	inode.n += rnode.n
}

func (inode *internalNode[Key, Value]) RemoveAt(i int) {
	copy(inode.nodes[i:], inode.nodes[i+1:inode.n])
	clear(inode.nodes[inode.n-1 : inode.n])
	copy(inode.keys[i:], inode.keys[i+1:inode.n])
	clear(inode.keys[inode.n-1 : inode.n])

	inode.n--
}

func (inode *internalNode[Key, Value]) ShiftLeft(o node[Key, Value], n int) {
	rnode := o.(*internalNode[Key, Value])

	copy(inode.nodes[inode.n:], rnode.nodes[:n])
	copy(inode.keys[inode.n:], rnode.keys[:n])

	inode.n += n

	copy(rnode.nodes[:], rnode.nodes[n:rnode.n])
	clear(rnode.nodes[rnode.n-n : rnode.n])
	copy(rnode.keys[:], rnode.keys[n:rnode.n])
	clear(rnode.keys[rnode.n-n : rnode.n])

	rnode.n -= n
}

func (inode *internalNode[Key, Value]) ShiftRight(o node[Key, Value], n int) {
	rnode := o.(*internalNode[Key, Value])

	copy(rnode.nodes[n:], rnode.nodes[:rnode.n])
	copy(rnode.keys[n:], rnode.keys[:rnode.n])

	rnode.n += n

	copy(rnode.nodes[:], inode.nodes[inode.n-n:inode.n])
	clear(inode.nodes[inode.n-n : inode.n])
	copy(rnode.keys[:], inode.keys[inode.n-n:inode.n])
	clear(inode.keys[inode.n-n : inode.n])

	inode.n -= n
}

type leafNode[Key comparable, Value any] struct {
	next   *leafNode[Key, Value]
	prev   *leafNode[Key, Value]
	values [maxNodes]Value
	keys   [maxNodes]Key
	n      int
}

func (tnode *leafNode[Key, Value]) Keys() []Key {
	return tnode.keys[:tnode.n]
}

func (tnode *leafNode[Key, Value]) Size() int {
	return tnode.n
}

func (tnode *leafNode[Key, Value]) LastKey() Key {
	return tnode.keys[tnode.n-1]
}

func (tnode *leafNode[Key, Value]) IsFull() bool {
	return tnode.n == maxNodes
}

func (tnode *leafNode[Key, Value]) Split(m *Map[Key, Value]) node[Key, Value] {
	rnode := &leafNode[Key, Value]{
		next: tnode.next,
	}

	tnode.next = rnode

	if rnode.next != nil {
		rnode.next.prev = rnode
	} else if m.back == tnode {
		m.back = rnode
	}

	rnode.prev = tnode

	n := tnode.n
	rnode.n = tnode.n / 2
	tnode.n -= rnode.n

	copy(rnode.values[:], tnode.values[tnode.n:n])
	clear(tnode.values[tnode.n:n])
	copy(rnode.keys[:], tnode.keys[tnode.n:n])
	clear(tnode.keys[tnode.n:n])

	return rnode
}

func (tnode *leafNode[Key, Value]) Merge(
	o node[Key, Value], m *Map[Key, Value],
) {
	rnode := o.(*leafNode[Key, Value])

	copy(tnode.values[tnode.n:], rnode.values[:rnode.n])
	copy(tnode.keys[tnode.n:], rnode.keys[:rnode.n])

	tnode.n += rnode.n
	tnode.next = rnode.next

	if tnode.next != nil {
		tnode.next.prev = tnode
	} else {
		m.back = tnode
	}
}

func (tnode *leafNode[Key, Value]) RemoveAt(i int) {
	copy(tnode.values[i:], tnode.values[i+1:tnode.n])
	clear(tnode.values[tnode.n-1 : tnode.n])
	copy(tnode.keys[i:], tnode.keys[i+1:tnode.n])
	clear(tnode.keys[tnode.n-1 : tnode.n])

	tnode.n--
}

func (tnode *leafNode[Key, Value]) ShiftLeft(o node[Key, Value], n int) {
	rnode := o.(*leafNode[Key, Value])

	copy(tnode.values[tnode.n:], rnode.values[:n])
	copy(tnode.keys[tnode.n:], rnode.keys[:n])

	tnode.n += n

	copy(rnode.values[:], rnode.values[n:rnode.n])
	clear(rnode.values[rnode.n-n : rnode.n])
	copy(rnode.keys[:], rnode.keys[n:rnode.n])
	clear(rnode.keys[rnode.n-n : rnode.n])

	rnode.n -= n
}

func (tnode *leafNode[Key, Value]) ShiftRight(o node[Key, Value], n int) {
	rnode := o.(*leafNode[Key, Value])

	copy(rnode.values[n:], rnode.values[:rnode.n])
	copy(rnode.keys[n:], rnode.keys[:rnode.n])

	rnode.n += n

	copy(rnode.values[:], tnode.values[tnode.n-n:tnode.n])
	clear(tnode.values[tnode.n-n : tnode.n])
	copy(rnode.keys[:], tnode.keys[tnode.n-n:tnode.n])
	clear(tnode.keys[tnode.n-n : tnode.n])

	tnode.n -= n
}

func (tnode *leafNode[Key, Value]) InsertAt(idx int, key Key, value Value) {
	copy(tnode.values[idx+1:], tnode.values[idx:tnode.n])
	copy(tnode.keys[idx+1:], tnode.keys[idx:tnode.n])

	tnode.keys[idx] = key
	tnode.values[idx] = value
	tnode.n++
}
