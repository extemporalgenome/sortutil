// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import "sort"

// ByteSlice attaches the methods of sort.Interface to []byte,
// sorting in increasing order.
type ByteSlice []byte

func (b ByteSlice) Len() int           { return len(b) }
func (b ByteSlice) Less(i, j int) bool { return b[i] < b[j] }
func (b ByteSlice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByteSlice) String() string     { return string(b) }

// NewLetterSeq returns an ascending Letters sequence of length n.
// If n is greater than 26, the sequence will be duplicated, starting
// again with 'a'.
func NewLetterSeq(n int) *Letters {
	b := make(ByteSlice, n)
	for i := range b {
		b[i] = 'a' + byte(i)%('z'-'a'+1)
	}
	return &Letters{b}
}

// Letters is designed for developing and debugging sorting algorithms,
// and should contain only bytes in the ASCII lowercase letter range.
type Letters struct{ ByteSlice }

// Mark behaves like String, except the specified indices will be uppercased.
func (l Letters) Mark(i, j int) string {
	b := l.ByteSlice
	c := make([]byte, len(b))
	copy(c, b)
	const o = 'a' - 'A'
	c[i] -= o
	c[j] -= o
	return string(c)
}

// NewIntSeq returns an ascending int sequence, starting with zero.
func NewIntSeq(n int) sort.IntSlice {
	s := make(sort.IntSlice, n)
	for i := range s {
		s[i] = i
	}
	return s
}
