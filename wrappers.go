// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import (
	"fmt"
	"io"
	"sort"
)

// Stat wraps sort.Interface, counting the number of Len, Less, and Swap calls.
// Initialize with `&Stat{I: data}`.
type Stat struct {
	I sort.Interface
	N struct{ Len, Less, Swap int }
}

func (s *Stat) Len() int           { s.N.Len++; return s.I.Len() }
func (s *Stat) Less(i, j int) bool { s.N.Less++; return s.I.Less(i, j) }
func (s *Stat) Swap(i, j int)      { s.N.Swap++; s.I.Swap(i, j) }
func (s *Stat) String() string     { return fmt.Sprintf("%+v", s.N) }

// Log wraps sort.Interface, sending debug messages to the supplied Writer.
// Less and Swap parameters will be space-padded based on the most recent Len
// call. Since writes are not synchronized, a serializing writer should be
// provided when used with concurrent sorting algorithms. A single *Log can be
// reused between separate sorts as long as they do not coincide.
type Log struct {
	I sort.Interface
	W io.Writer
	p int
}

func (l *Log) Len() int {
	r := l.I.Len()
	l.p = 0
	if r > 0 {
		l.p = len(fmt.Sprint(r - 1))
	}
	fmt.Fprint(l.W, "Len() [", r, "]\n")
	return r
}

func (l *Log) Less(i, j int) bool {
	r := l.I.Less(i, j)
	fmt.Fprintf(l.W, "Less(%*d, %*d) [%v]\n", l.p, i, l.p, j, r)
	return r
}

func (l *Log) Swap(i, j int) {
	l.I.Swap(i, j)
	fmt.Fprintf(l.W, "Swap(%*d, %*d)\n", l.p, i, l.p, j)
}
