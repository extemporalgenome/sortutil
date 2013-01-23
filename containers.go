// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import "sort"

type ByteSlice []byte

func (b ByteSlice) Len() int           { return len(b) }
func (b ByteSlice) Less(i, j int) bool { return b[i] < b[j] }
func (b ByteSlice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByteSlice) String() string     { return string(b) }

func NewLetterSeq(n int) *Letters {
	b := make(ByteSlice, n)
	for i := range b {
		b[i] = 'a' + byte(i)%('z'-'a'+1)
	}
	return &Letters{b}
}

type Letters struct{ ByteSlice }

func (l Letters) Mark(i, j int) string {
	b := l.ByteSlice
	c := make([]byte, len(b))
	copy(c, b)
	const o = 'a' - 'A'
	c[i] -= o
	c[j] -= o
	return string(c)
}

func NewIntSeq(n int) sort.IntSlice {
	s := make(sort.IntSlice, n)
	for i := range s {
		s[i] = i
	}
	return s
}
