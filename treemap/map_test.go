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
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func verifyMapNode[Key, Value any](
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

func verifyMapLen[Key, Value any](t *testing.T, m *Map[Key, Value]) {
	t.Helper()

	sum := 0

	for tnode := m.front; tnode != nil; tnode = tnode.next {
		sum += tnode.n
	}

	assert.Equal(t, sum, m.Len())
}

func verifyMapKey[Key, Value any](
	t *testing.T, m *Map[Key, Value], minKey Key,
) {
	t.Helper()

	key := minKey

	for k := range m.Begin().Seq() {
		assert.LessOrEqual(t, key, k)
		key = k
	}
}

func verifyMap[Key, Value any](
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

func printMap[Key, Value any](m *Map[Key, Value]) { //nolint:unused
	printMapNode(m.root, 0)
}

func printMapNode[Key, Value any]( //nolint:unused
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

	it, oldValue, ok := m.Insert("foo", 1)

	require.False(t, it.End())
	assert.Zero(t, oldValue)
	assert.False(t, ok)
	assert.Equal(t, "foo", it.Key())
	assert.Equal(t, 1, it.Value())

	it = it.Next()

	require.True(t, it.End())

	it, oldValue, ok = m.Insert("bar", 2)

	require.False(t, it.End())
	assert.Zero(t, oldValue)
	assert.False(t, ok)
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

	it, oldValue, ok = m.Insert("foo", 100)

	require.False(t, it.End())
	assert.Equal(t, 1, oldValue)
	assert.True(t, ok)
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
	it, _, _ := m.Insert(32, 99)

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
		it, _, _ := m.Insert(i, i+1)

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

	oldValue, ok := m.Remove(0)

	assert.Zero(t, oldValue)
	assert.False(t, ok)

	m.Insert(912, 1)
	m.Insert(313, 2)
	m.Insert(78, 3)

	assert.Equal(t, 3, m.Len())

	oldValue, ok = m.Remove(912)

	assert.Equal(t, 1, oldValue)
	assert.True(t, ok)
	assert.Equal(t, 2, m.Len())

	oldValue, ok = m.Remove(912)

	assert.Zero(t, oldValue)
	assert.False(t, ok)
	assert.Equal(t, []int{78, 313}, slices.Collect(m.Keys()))
	assert.Equal(t, []int{3, 2}, slices.Collect(m.Values()))

	verifyMap(t, m, 78, 912)

	m = New[int, int]()

	for i := range 1000 {
		m.Insert(i, i)
	}

	oldValue, ok = m.Remove(1000)

	assert.Zero(t, oldValue)
	assert.False(t, ok)
}

func TestMapRemoveIter(t *testing.T) {
	m := New[int, int]()

	it := m.RemoveIter(m.End())

	assert.True(t, it.End())

	m.Insert(912, 1)
	m.Insert(78, 3)
	it, _, _ = m.Insert(313, 2)

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
		oldValue, ok := m.Remove(i)

		assert.Equal(t, i+1, oldValue)
		assert.True(t, ok)
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

func TestMapNewAny(t *testing.T) {
	m := NewAny[string, int](cmp.Compare[string])

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

func TestMapInsertRemoveSplitExtendKey(t *testing.T) {
	m := New[uint64, int]()

	m.Insert(9631484016779065335, 0)
	m.Insert(17868022691004923650, 0)
	m.Remove(39359)
	m.Insert(13833645161281420026, 0)
	m.Insert(15555861491690288297, 0)
	m.Insert(10430266092290031551, 0)
	m.Insert(13775596190496173567, 0)
	m.Insert(17632622606063210373, 0)
	m.Insert(9652578060282094057, 0)
	m.Insert(11068046731657850367, 0)
	m.Insert(17605959298697047211, 0)
	m.Insert(9645450310482657280, 0)
	m.Insert(12576213941979451588, 0)
	m.Insert(13262819427629400064, 0)
	m.Insert(17353126493894524672, 0)
	m.Insert(6968380983289572249, 0)
	m.Insert(11788714069324677567, 0)
	m.Insert(13816973012063010751, 0)
	m.Insert(13811842810106927258, 0)
	m.Insert(15555861465928021161, 0)
	m.Insert(10416739736715106780, 0)
	m.Insert(16406050011325329075, 0)
	m.Insert(15316432918403723250, 0)
	m.Insert(11166062425818695557, 0)
	m.Insert(15316179399658878106, 0)
	m.Insert(16035721616952779220, 0)
	m.Remove(34560)
	m.Insert(17498454949193438613, 0)
	m.Insert(16847020544331923391, 0)
	m.Insert(13816973012072595456, 0)
	m.Remove(0)
	m.Remove(0)
	m.Remove(0)
	m.Remove(0)
	m.Remove(49087)
	m.Insert(13810192568557345177, 0)
	m.Remove(11068046444225731031)
	m.Insert(11349633585725296300, 0)
	m.Insert(18415206453810212011, 0)
	m.Insert(9645450313596580311, 0)
	m.Insert(9625797500427571396, 0)
	m.Insert(13316529160701451649, 0)
	m.Insert(11312971619254663900, 0)
	m.Insert(11067989837730010939, 0)
	m.Insert(1206684267757682623, 0)
	m.Insert(10777272142887491981, 0)
	m.Insert(15551711123532140417, 0)
	m.Insert(10148182166526213961, 0)
	m.Insert(0, 0)
	m.Insert(11068046431340829081, 0)
	m.Insert(11068046444229754896, 0)
	m.Insert(4613099049867460543, 0)
	m.Insert(13810192565461621145, 0)
	m.Remove(11068046444225730969)
	m.Remove(1208329986916299417)
	m.Insert(11068046444225730969, 0)
	m.Insert(13810192568557345061, 0)
	m.Remove(11068088389516378047)
	m.Insert(12081179472156006709, 0)
	m.Remove(361736048658595775)
	m.Insert(12081179472156006706, 0)
	m.Remove(9365982719239559577)
	m.Insert(3646114258319153561, 0)
	m.Insert(4719896757377329, 0)
	m.Insert(11068046444225730969, 0)
	m.Insert(11078757952523901337, 0)
	m.Insert(3646114258319153561, 0)
	m.Insert(11068046444225730969, 0)
	m.Insert(4720039014813879, 0)
	m.Remove(11068046444225730969)
	m.Insert(11068046444225746688, 0)
	m.Remove(3646114258319153561)
	m.Insert(11068046444225730969, 0)
	m.Insert(4720009734379471, 0)
	m.Insert(16130445648890337586, 0)
	m.Insert(11068046444225730969, 0)
	m.Insert(11068046444225746688, 0)
	m.Remove(3646114173938499584)
	m.Remove(11362168099665543184)
	m.Remove(9647746139901155054)
	m.Remove(2173955599088251055)
	m.Remove(15536000149505360519)
	m.Insert(16272787484313703092, 0)
	m.Remove(18149095641086133729)
	m.Remove(10232178353385767047)
	m.Remove(12678432296282675669)
	m.Insert(9073007928742288575, 0)
	m.Insert(13816973012072644543, 0)
	m.Insert(13816946525965162905, 0)
	m.Insert(3837072402173725375, 0)
	m.Insert(13816946525965162905, 0)
	m.Insert(3646114258319153561, 0)
	m.Insert(4720039011391794, 0)
	m.Insert(11068046444225730969, 0)
	m.Insert(13816946525965162905, 0)
	m.Insert(3646114258319153561, 0)
	m.Insert(11068046444225730969, 0)
	m.Insert(4720039365759159, 0)
	m.Remove(17860931781328224050)
	m.Insert(11068046444225730969, 0)
	m.Insert(11068046444225746688, 0)
	m.Remove(14184571889434270105)
	m.Insert(3646114258319153561, 0)
	m.Insert(4720040018431873, 0)
	m.Insert(11068046444225730969, 0)
	m.Insert(11068046444225746688, 0)
	m.Remove(14166522987856632217)
	m.Insert(4719107240869823, 0)
	m.Insert(13806235017666992537, 0)
	m.Insert(13373889453439440640, 0)
	m.Remove(3646114258319153561)
	m.Insert(4720039014813743, 0)
	m.Insert(15492401155806907856, 0)
	m.Insert(13301831859401496985, 0)
	m.Insert(11068046444225731031, 0)
	m.Remove(1167164138428406169)
	m.Insert(11068046444225730969, 0)
	m.Insert(15492401155692531159, 0)
	m.Insert(11520167005346519705, 0)
	m.Remove(11068046444225730969)
	m.Insert(11068046444225731031, 0)
	m.Remove(1167164138098794240)
	m.Remove(44383469139318528)
	m.Remove(14160974939792864342)
	m.Remove(5484869160941461616)
	m.Insert(17498625257762565826, 0)
	m.Insert(9647225583155015890, 0)
	m.Remove(11167764386688894853)
	m.Insert(15316179382733635584, 0)
	m.Insert(49525126157354201, 0)
	m.Insert(10772022948872912028, 0)
	m.Insert(13816973012072644543, 0)
	m.Insert(13816972908611287193, 0)
	m.Insert(11039800451873965322, 0)
	m.Insert(13816972908611287193, 0)
	m.Insert(11039054521624533401, 0)
	m.Insert(11024830325455362457, 0)
	m.Remove(11068046444225730969)
	m.Insert(11078784335170025625, 0)
	m.Insert(2680373613224892825, 0)
	m.Insert(11068046444225730969, 0)

	verifyMap(t, m, 0, math.MaxUint64)
}
