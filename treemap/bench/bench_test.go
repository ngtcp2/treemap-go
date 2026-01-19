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

package bench

import (
	"cmp"
	"math/rand/v2"
	"runtime"
	"testing"

	gods "github.com/emirpasic/gods/v2/maps/treemap"
	google "github.com/google/btree"
	"github.com/ngtcp2/treemap-go/treemap"
	"github.com/tidwall/btree"
	k8s "k8s.io/utils/third_party/forked/golang/btree"
)

const N = 10000

func makeArray(n int) ([]int, []int) {
	rnd := rand.New(rand.NewPCG(1000000007, 1000000009))

	return rnd.Perm(n), rnd.Perm(n)
}

var a, d = makeArray(N)

type Foo struct {
	X int
}

func compareFoo(x, y Foo) int {
	return cmp.Compare(x.X, y.X)
}

func BenchmarkInsertRand(b *testing.B) {
	for b.Loop() {
		m := treemap.New[int, int]()

		for _, k := range a {
			m.Insert(k, k)
		}
	}
}

func BenchmarkLookupRand(b *testing.B) {
	m := treemap.New[int, int]()

	for _, k := range a {
		m.Insert(k, k)
	}

	for b.Loop() {
		for _, k := range d {
			m.Find(k)
		}
	}
}

func BenchmarkIterateRand(b *testing.B) {
	m := treemap.New[int, int]()

	for _, k := range a {
		m.Insert(k, k)
	}

	f := func() {
		it := m.Begin()

		for !it.End() {
			runtime.KeepAlive(it.Key())
			runtime.KeepAlive(it.Value())
			it = it.Next()
		}
	}

	for b.Loop() {
		f()
	}
}

func BenchmarkIterateSeqRand(b *testing.B) {
	m := treemap.New[int, int]()

	for _, k := range a {
		m.Insert(k, k)
	}

	f := func() {
		for k, v := range m.Begin().Seq() {
			runtime.KeepAlive(k)
			runtime.KeepAlive(v)
		}
	}

	for b.Loop() {
		f()
	}
}

func BenchmarkRemoveRand(b *testing.B) {
	for b.Loop() {
		b.StopTimer()

		m := treemap.New[int, int]()

		for _, k := range a {
			m.Insert(k, k)
		}

		b.StartTimer()

		for _, k := range d {
			m.Remove(k)
		}
	}
}

func BenchmarkInsertComparableRand(b *testing.B) {
	for b.Loop() {
		m := treemap.NewAny[Foo, int](compareFoo)

		for _, k := range a {
			m.Insert(Foo{X: k}, k)
		}
	}
}

func BenchmarkLookupComparableRand(b *testing.B) {
	m := treemap.NewAny[Foo, int](compareFoo)

	for _, k := range a {
		m.Insert(Foo{X: k}, k)
	}

	for b.Loop() {
		for _, k := range d {
			m.Find(Foo{X: k})
		}
	}
}

func BenchmarkIterateComparableRand(b *testing.B) {
	m := treemap.NewAny[Foo, int](compareFoo)

	for _, k := range a {
		m.Insert(Foo{X: k}, k)
	}

	f := func() {
		it := m.Begin()

		for !it.End() {
			runtime.KeepAlive(it.Key())
			runtime.KeepAlive(it.Value())
			it = it.Next()
		}
	}

	for b.Loop() {
		f()
	}
}

func BenchmarkRemoveComparableRand(b *testing.B) {
	for b.Loop() {
		b.StopTimer()

		m := treemap.NewAny[Foo, int](compareFoo)

		for _, k := range a {
			m.Insert(Foo{X: k}, k)
		}

		b.StartTimer()

		for _, k := range d {
			m.Remove(Foo{X: k})
		}
	}
}

func BenchmarkGODSInsertRand(b *testing.B) {
	for b.Loop() {
		m := gods.New[int, int]()

		for _, k := range a {
			m.Put(k, k)
		}
	}
}

func BenchmarkGODSLookupRand(b *testing.B) {
	m := gods.New[int, int]()

	for _, k := range a {
		m.Put(k, k)
	}

	for b.Loop() {
		for _, k := range d {
			m.Get(k)
		}
	}
}

func BenchmarkGODSIterateRand(b *testing.B) {
	m := gods.New[int, int]()

	for _, k := range a {
		m.Put(k, k)
	}

	f := func() {
		it := m.Iterator()

		for it.Next() {
			runtime.KeepAlive(it.Key())
			runtime.KeepAlive(it.Value())
		}
	}

	for b.Loop() {
		f()
	}
}

func BenchmarkGODSRemoveRand(b *testing.B) {
	for b.Loop() {
		b.StopTimer()

		m := gods.New[int, int]()

		for _, k := range a {
			m.Put(k, k)
		}

		b.StartTimer()

		for _, k := range d {
			m.Remove(k)
		}
	}
}

func BenchmarkGODSInsertComparableRand(b *testing.B) {
	for b.Loop() {
		m := gods.NewWith[Foo, int](compareFoo)

		for _, k := range a {
			m.Put(Foo{X: k}, k)
		}
	}
}

func BenchmarkGODSLookupComparableRand(b *testing.B) {
	m := gods.NewWith[Foo, int](compareFoo)

	for _, k := range a {
		m.Put(Foo{X: k}, k)
	}

	for b.Loop() {
		for _, k := range d {
			m.Get(Foo{X: k})
		}
	}
}

func BenchmarkGODSIterateComparableRand(b *testing.B) {
	m := gods.NewWith[Foo, int](compareFoo)

	for _, k := range a {
		m.Put(Foo{X: k}, k)
	}

	f := func() {
		it := m.Iterator()

		for it.Next() {
			runtime.KeepAlive(it.Key())
			runtime.KeepAlive(it.Value())
		}
	}

	for b.Loop() {
		f()
	}
}

func BenchmarkGODSRemoveComparableRand(b *testing.B) {
	for b.Loop() {
		b.StopTimer()

		m := gods.NewWith[Foo, int](compareFoo)

		for _, k := range a {
			m.Put(Foo{X: k}, k)
		}

		b.StartTimer()

		for _, k := range d {
			m.Remove(Foo{X: k})
		}
	}
}

const btDegree = 16

func BenchmarkBTInsertRand(b *testing.B) {
	for b.Loop() {
		m := btree.NewMap[int, int](btDegree)

		for _, k := range a {
			m.Set(k, k)
		}
	}
}

func BenchmarkBTLookupRand(b *testing.B) {
	m := btree.NewMap[int, int](btDegree)

	for _, k := range a {
		m.Set(k, k)
	}

	for b.Loop() {
		for _, k := range d {
			m.Get(k)
		}
	}
}

func BenchmarkBTIterateRand(b *testing.B) {
	m := btree.NewMap[int, int](btDegree)

	for _, k := range a {
		m.Set(k, k)
	}

	f := func() {
		it := m.Iter()

		for it.Next() {
			runtime.KeepAlive(it.Key())
			runtime.KeepAlive(it.Value())
		}
	}

	for b.Loop() {
		f()
	}
}

func BenchmarkBTRemoveRand(b *testing.B) {
	for b.Loop() {
		b.StopTimer()

		m := btree.NewMap[int, int](btDegree)

		for _, k := range a {
			m.Set(k, k)
		}

		b.StartTimer()

		for _, k := range d {
			m.Delete(k)
		}
	}
}

func BenchmarkGoogleInsertRand(b *testing.B) {
	for b.Loop() {
		m := google.NewOrderedG[int](btDegree)

		for _, k := range a {
			m.ReplaceOrInsert(k)
		}
	}
}

func BenchmarkGoogleLookupRand(b *testing.B) {
	m := google.NewOrderedG[int](btDegree)

	for _, k := range a {
		m.ReplaceOrInsert(k)
	}

	for b.Loop() {
		for _, k := range d {
			m.Get(k)
		}
	}
}

func BenchmarkGoogleIterateRand(b *testing.B) {
	m := google.NewOrderedG[int](btDegree)

	for _, k := range a {
		m.ReplaceOrInsert(k)
	}

	for b.Loop() {
		m.Ascend(func(int) bool { return true })
	}
}

func BenchmarkGoogleRemoveRand(b *testing.B) {
	for b.Loop() {
		b.StopTimer()

		m := google.NewOrderedG[int](btDegree)

		for _, k := range a {
			m.ReplaceOrInsert(k)
		}

		b.StartTimer()

		for _, k := range d {
			m.Delete(k)
		}
	}
}

func BenchmarkK8sInsertRand(b *testing.B) {
	for b.Loop() {
		m := k8s.NewOrdered[int](btDegree)

		for _, k := range a {
			m.ReplaceOrInsert(k)
		}
	}
}

func BenchmarkK8sLookupRand(b *testing.B) {
	m := k8s.NewOrdered[int](btDegree)

	for _, k := range a {
		m.ReplaceOrInsert(k)
	}

	for b.Loop() {
		for _, k := range d {
			m.Get(k)
		}
	}
}

func BenchmarkK8sIterateRand(b *testing.B) {
	m := k8s.NewOrdered[int](btDegree)

	for _, k := range a {
		m.ReplaceOrInsert(k)
	}

	for b.Loop() {
		m.Ascend(func(int) bool { return true })
	}
}

func BenchmarkK8sRemoveRand(b *testing.B) {
	for b.Loop() {
		b.StopTimer()

		m := k8s.NewOrdered[int](btDegree)

		for _, k := range a {
			m.ReplaceOrInsert(k)
		}

		b.StartTimer()

		for _, k := range d {
			m.Delete(k)
		}
	}
}
