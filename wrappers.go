// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import (
	"fmt"
	"io"
	"sort"
)

const panicmsg = "bounds out of range"

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
	fmt.Fprint(l.W, "(", l.I, ").Len() [", r, "]\n")
	return r
}

func (l *Log) Less(i, j int) bool {
	r := l.I.Less(i, j)
	fmt.Fprintf(l.W, "(%v).Less(%*d, %*d) [%v]\n", l.I, l.p, i, l.p, j, r)
	return r
}

func (l *Log) Swap(i, j int) {
	v := fmt.Sprint(l.I)
	l.I.Swap(i, j)
	fmt.Fprintf(l.W, "(%v).Swap(%*d, %*d) [%v]\n", v, l.p, i, l.p, j, l.I)
}

func (l *Log) String() string {
	return fmt.Sprint(l.I)
}

// NewSub opaquely wraps a sub-sequence of the provided sort.Interface.
// NewSub(s,i,j) is semantically equivalent to s[i:j], though the underlying
// implementation does not need to involve a slice. j may not exceed s.Len().
func NewSub(s sort.Interface, i, j int) sort.Interface {
	if i < 0 || j < i || j > s.Len() {
		panic(panicmsg)
	} else if v, ok := s.(sub); ok {
		// collapse subs of subs
		return sub{v.s, v.i + i, j - i}
	}
	return sub{s, i, j - i}
}

type sub struct {
	s    sort.Interface
	i, n int
}

func (s sub) Len() int           { return s.n }
func (s sub) Less(i, j int) bool { return s.s.Less(s.i+i, s.i+j) }
func (s sub) Swap(i, j int)      { s.s.Swap(s.i+i, s.i+j) }

// NewRev returns a reverse sorter for any sort.Interface.
// To quickly reverse a sort.Interface relative to its current order, see Reverse.
func NewRev(s sort.Interface) sort.Interface {
	if v, ok := s.(rev); ok {
		return v.Interface
	}
	return rev{s}
}

type rev struct{ sort.Interface }

func (r rev) Less(i, j int) bool { return !r.Interface.Less(i, j) }

// NewProxy sorts comp, duplicating all swaps on each item of data.
// data may be either be empty, or each item must have the same Len as comp.
func NewProxy(comp sort.Interface, data ...sort.Interface) sort.Interface {
	l := comp.Len()
	for _, d := range data {
		if l != d.Len() {
			panic(panicmsg)
		}
	}
	return proxy{comp, data}
}

type proxy struct {
	c sort.Interface
	d []sort.Interface
}

func (p proxy) Len() int           { return p.c.Len() }
func (p proxy) Less(i, j int) bool { return p.c.Less(i, j) }

func (p proxy) Swap(i, j int) {
	p.c.Swap(i, j)
	for _, d := range p.d {
		d.Swap(i, j)
	}
}
