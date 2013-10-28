// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

const panicmsg = "bounds out of range"

// NewLogStat wraps data with Log and Stat.
// When Log and Stat are used together, Stat should wrap Log. NewLogStat
// is a convenience function for correctly composing the wrappers.
func NewLogStat(w io.Writer, data sort.Interface) *Stat {
	return NewStat(&Log{I: data, W: w})
}

// NewStat initializes a *Stat for recording per-element call counts.
// NewStat makes a call to data.Len.
func NewStat(data sort.Interface) *Stat {
	l := 0
	if log, ok := data.(*Log); ok {
		l = log.I.Len()
	} else {
		l = data.Len()
	}
	return &Stat{
		I: data,
		O: make([]struct{ Less, Swap int }, l),
	}
}

// Stat wraps sort.Interface, counting the number of Len, Less, and Swap calls.
// Initialize with `&Stat{I: data}`, or use NewStat to initialize for more
// comprehensive statistics.
type Stat struct {
	I sort.Interface
	N struct{ Len, Less, Swap int }
	O []struct{ Less, Swap int }
}

func (s *Stat) Len() int { s.N.Len++; return s.I.Len() }

func (s *Stat) Less(i, j int) bool {
	s.N.Less++
	if s.O != nil {
		s.O[i].Less++
		s.O[j].Less++
	}
	return s.I.Less(i, j)
}

func (s *Stat) Swap(i, j int) {
	s.N.Swap++
	if s.O != nil {
		s.O[i].Swap++
		s.O[j].Swap++
	}
	s.I.Swap(i, j)
}

// StatAggregate contains a summary of element-wise call statistics.
// Index zero represents Less, while index one represents Swap.
type StatAggregate [2]struct {
	Min, Max  int
	Mean, Std float32
}

// Aggregate will return aggregate statistics.
// If the *Stat was not initialized via NewStat, a zero-valued StatAggregate
// will be returned.
func (s *Stat) Aggregate() StatAggregate {
	var a StatAggregate
	n := len(s.O)
	if n == 0 {
		return a
	}
	lMean := float32(s.N.Less) / float32(n)
	sMean := float32(s.N.Swap) / float32(n)
	lMin, sMin := n, n
	lMax, sMax := 0, 0
	for _, v := range s.O {
		l, s := v.Less, v.Swap
		std := float32(l) - lMean
		a[0].Std += std * std
		std = float32(s) - sMean
		a[1].Std += std * std
		if l < lMin {
			lMin = l
		}
		if l > lMax {
			lMax = l
		}
		if s < sMin {
			sMin = s
		}
		if s > sMax {
			sMax = s
		}
	}
	a[0].Std = float32(math.Sqrt(float64(a[0].Std) / float64(n)))
	a[1].Std = float32(math.Sqrt(float64(a[1].Std) / float64(n)))
	a[0].Mean, a[1].Mean = lMean, sMean
	a[0].Min, a[0].Max, a[1].Min, a[1].Max = lMin, lMax, sMin, sMax
	return a
}

// String summarizes the statistical results and, if possible, aggregated results.
func (s *Stat) String() string {
	if s.O == nil {
		return fmt.Sprintf("Calls: %+v", s.N)
	}
	a := s.Aggregate()
	return fmt.Sprintf("Calls: %+v\nLess:  %+v\nSwap:  %+v", s.N, a[0], a[1])
}

// Mark should produce output with the same visible length that fmt.Sprint
// would produce when passed the receiver. Within the same alignment
// constraints, the returned string should emphasize the indices i and j.
// This is an optional method for use with Log.
type Marker interface {
	Mark(i, j int) string
}

// Log wraps sort.Interface, sending debug messages to the supplied Writer.
// Less and Swap parameters will be space-padded based on the most recent Len
// call. Since writes are not synchronized, a serializing writer should be
// provided when used with concurrent sorting algorithms. A single *Log can be
// reused between separate sorts as long as they do not coincide.
//
// If the sort.Interface value implements Marker, Mark will be called for Less
// and Swap if Len() returned a small enough value.
type Log struct {
	I    sort.Interface
	W    io.Writer
	n, p int
}

// LOG_ITEM_THRESH is the maximum Len of a sort.Interface that will be printed
// inline with log messages. If negative, inline display of the data will be
// disabled.
var LOG_ITEM_THRESH = 26

func (l *Log) Mark(i, j int) string {
	if m, ok := l.I.(Marker); ok {
		return m.Mark(i, j)
	}
	return fmt.Sprint(l.I)
}

func (l *Log) Len() int {
	r := l.I.Len()
	l.n, l.p = r, 0
	if r > 0 {
		l.p = len(fmt.Sprint(r - 1))
	}
	if r <= LOG_ITEM_THRESH {
		fmt.Fprint(l.W, "(", l.I, ").Len() [", r, "]\n")
	} else {
		fmt.Fprint(l.W, "Len() [", r, "]\n")
	}
	return r
}

func (l *Log) Less(i, j int) bool {
	r := l.I.Less(i, j)
	if l.n <= LOG_ITEM_THRESH && l.n > 0 {
		fmt.Fprintf(l.W, "(%v).Less(%*d, %*d) [%v]\n", l.Mark(i, j), l.p, i, l.p, j, r)
	} else {
		fmt.Fprintf(l.W, "Less(%*d, %*d) [%v]\n", l.p, i, l.p, j, r)
	}
	return r
}

func (l *Log) Swap(i, j int) {
	if l.n > LOG_ITEM_THRESH || l.n <= 0 {
		l.I.Swap(i, j)
		fmt.Fprintf(l.W, "Swap(%*d, %*d)\n", l.p, i, l.p, j)
		return
	}
	v := l.Mark(i, j)
	l.I.Swap(i, j)
	fmt.Fprintf(l.W, "(%v).Swap(%*d, %*d) [%v]\n", v, l.p, i, l.p, j, l.Mark(i, j))
}

func (l *Log) String() string {
	return fmt.Sprint(l.I)
}

// NewSub opaquely wraps a sub-sequence of the provided sort.Interface.
// NewSub(s,i,j) is semantically equivalent to s[i:j], though the underlying
// implementation does not need to use a slice.
// NewSub will panic unless 0 <= i <= j <= s.Len().
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

func (r rev) Less(i, j int) bool { return r.Interface.Less(j, i) }

// NewProxy sorts comp, duplicating all swaps on each item of data.
// NewProxy will panic if any item in data has a different Len() than comp.
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

// Analyze runs preselected datasets through the sorting function f.
// Any runs that fail to be correctly sorted will be listed first. For each
// run, if verbose is true or a run fails its Len, Less, and Swap calls will
// be logged to the provided Writer. In all cases, a summary of call count
// statistics will be written to the Writer.
func Analyze(w io.Writer, verbose bool, f func(sort.Interface)) {
	tests := [][2]string{
		{"qozxgwajmcnisphfldterkvbuy", "Shuffle"},
		{"abcdefghijklmnopqrstuvwxyz", "Ascending"},
		{"zyxwvutsrqponmlkjihgfedcba", "Descending"},
		{"badcfehgjilknmporqtsvuxwzy", "Pair-Transposition"},
		{"azcxevgtirkpmnolqjshufwdyb", "Zig-Zag"},
		{"zaxcvetgripknmlojqhsfudwby", "Desc-Zag-Trans"},
		{"qogwajmcnisphfldterkvbu", "Shuffle Prime"},
	}
	n := len(tests)
	succ := make([]int, 0, n*2)
	succ, fail := succ[:0], succ[n:n]
	var data Letters
	tlen := 0
	// Sort failures first
	for i, v := range tests {
		data = append(data[:0], v[0]...)
		title := v[1]
		if len(title) > tlen {
			tlen = len(title)
		}
		f(data)
		if sort.IsSorted(data) {
			succ = append(succ, i)
		} else {
			fail = append(fail, i)
		}
	}
	n = len(fail)
	pad := 4 + 7 + 4
	banner := strings.Repeat("#", tlen+pad)
	for i, j := range append(fail, succ...) {
		v := tests[j]
		data = append(data[:0], v[0]...)
		title := v[1]
		status := "[ OK ]"
		stat := NewStat(data)
		switch {
		case i < n:
			status = "[FAIL]"
			fallthrough
		case verbose:
			stat.I = &Log{I: data, W: w}
		}
		fmt.Fprintf(w, "%s\n### %s %-*s ###\n%s\n", banner, status, tlen, title, banner)
		f(stat)
		fmt.Fprint(w, "\n", stat, "\n\n")
	}
}
