// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import (
	"bytes"
	"sort"
	"testing"
)

func TestNewIntSeq(t *testing.T) {
	for i, v := range NewIntSeq(8) {
		if i != v {
			t.FailNow()
		}
	}
}

func TestNewLetterSeq(t *testing.T) {
	s := NewLetterSeq(27).String()
	if s[:2] != "ab" || s[25:] != "za" {
		t.Fail()
	}
}

func TestLettersMark(t *testing.T) {
	s := NewLetterSeq(10).Mark(2, 4)
	if s[2] != 'C' || s[4] != 'E' || s[5] != 'f' {
		t.Fail()
	}
}

func TestNewSub(t *testing.T) {
	data := [8]int{7, 6, 5, 4, 3, 2, 1, 0}
	sort.Sort(NewSub(sort.IntSlice(data[:]), 4, 8))
	if data != [8]int{7, 6, 5, 4, 0, 1, 2, 3} {
		t.Fail()
	}
}

func TestStat(t *testing.T) {
	var (
		s    = &Stat{I: ByteSlice{0}}
		Len  = 3
		Less = 14
		Swap = 7
	)
	for i := 0; i < Len; i++ {
		s.Len()
	}
	for i := 0; i < Less; i++ {
		s.Less(0, 0)
	}
	for i := 0; i < Swap; i++ {
		s.Swap(0, 0)
	}
	if s.N.Len != Len || s.N.Less != Less || s.N.Swap != Swap {
		t.Fail()
	}
}

func TestReverse(t *testing.T) {
	s := NewLetterSeq(26).ByteSlice
	Reverse(s)
	l := byte(len(s))
	for i := range s {
		i := byte(i)
		if s[l-i-1] != 'a'+i {
			t.Fail()
		}
	}
}

func TestShuffle(t *testing.T) {
	b := NewLetterSeq(26).ByteSlice
	s := b.String()
	Shuffle(b)
	if s == b.String() {
		t.Fail()
	}
}

func TestRotate(t *testing.T) {
	const n = 29
	b := NewLetterSeq(n).ByteSlice
	c := make(ByteSlice, n)
	d := make(ByteSlice, n)
	for i := n; i > 0; i-- {
		b, c, d = b[:i], c[:i], d[:i]
		for j := -i - 1; j < i+1; j++ {
			copy(c, b)
			copy(d, b)
			j := j % i
			k := -j
			if k < 0 {
				k += i
			}
			copy(c, c[k:])
			copy(c[i-k:], d[:k])
			Rotate(d, j)
			same := bytes.Equal([]byte(c), []byte(d))
			if !same {
				t.Errorf("%3d %s", j, c)
				t.Errorf("    %s", d)
			}
		}
	}
}

var skewTests = []struct {
	r       string
	i, j, k int
}{
	{"bcdefghijklma", 0, 12, 1},
	{"fghijklmabcde", 0, 8, 5},
	{"abcdeijklfghm", 5, 9, 3},
	{"abjcdefghik", 2, 3, 7},
	{"defabcghij", 0, 3, 3},
	{"hijabcdefg", 7, 0, 3},
	{"abcdehijfg", 7, 5, 3},
	{"afgbcde", 1, 3, 4},
}

func TestSkew(t *testing.T) {
	for i, v := range skewTests {
		try := func(p, q int) {
			b := NewLetterSeq(len(v.r)).ByteSlice
			Skew(b, p, q, v.k)
			if string(b) != v.r {
				t.Errorf("#%2d [%2d %2d %2d] %s", i, p, q, v.k, v.r)
				t.Errorf("              %s", b)
			}
		}
		try(v.i, v.j)
	}
}
