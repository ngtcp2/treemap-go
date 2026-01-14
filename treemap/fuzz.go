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

//go:build fuzz

package treemap

import (
	"bytes"
	"encoding/binary"
)

type FuzzerProvider struct {
	buf *bytes.Reader
}

func NewFuzzerProvider(buf []byte) *FuzzerProvider {
	return &FuzzerProvider{
		buf: bytes.NewReader(buf),
	}
}

func (fp *FuzzerProvider) ConsumeUint32() (uint32, bool) {
	var n uint32

	err := binary.Read(fp.buf, binary.BigEndian, &n)
	if err != nil {
		return 0, false
	}

	return n, true
}

func (fp *FuzzerProvider) ConsumeBool() (bool, bool) {
	var b bool

	err := binary.Read(fp.buf, binary.BigEndian, &b)
	if err != nil {
		return false, false
	}

	return b, true
}

func FuzzMap(input []byte) int {
	m := New[uint32, uint32]()
	fp := NewFuzzerProvider(input)

	for {
		key, ok := fp.ConsumeUint32()
		if !ok {
			break
		}

		insert, ok := fp.ConsumeBool()
		if !ok {
			break
		}

		remove, ok := fp.ConsumeBool()
		if !ok {
			break
		}

		if insert {
			m.Insert(key, key)
		}

		if remove {
			m.Remove(key)
		}

		m.LowerBound(key)
	}

	return 0
}
