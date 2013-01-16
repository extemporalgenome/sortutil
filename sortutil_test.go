// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import (
	"testing"
)

func TestNewLetters(t *testing.T) {
	s := NewLetters(27).String()
	if s[:2] != "ab" || s[25:] != "za" {
		t.Fail()
	}
}

func TestLettersMark(t *testing.T) {
	s := NewLetters(10).Mark(2, 4)
	if s[2] != 'C' || s[4] != 'E' || s[5] != 'f' {
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
	s := NewLetters(26).ByteSlice
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
	b := NewLetters(26).ByteSlice
	s := b.String()
	Shuffle(b)
	if s == b.String() {
		t.Fail()
	}
}

/*
func TestSkew(t *testing.T) {
	b := NewLetters(26).ByteSlice
}
*/
