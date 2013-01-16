// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import (
	"math/rand"
	"sort"
)

func Reverse(data sort.Interface) {
	n := data.Len()
	for i := 0; i < n/2; i++ {
		data.Swap(i, n-i-1)
	}
}

func Skew(data sort.Interface, k int) {
	n := data.Len()
	k %= n
	if k == 0 || n < 2 {
		return
	}

	left := func(u, v, k int) {
		for i := u; i < v-k; i++ {
			data.Swap(i, i+k)
		}
	}
	right := func(u, v, k int) {
		for i := v - 1; i >= u+k; i-- {
			data.Swap(i, i-k)
		}
	}

	if k > n/2 {
		k -= n
	} else if k < -n/2 {
		k += n
	}
	p := n % k
	if k > 0 {
		right(p, n, k)
		left(0, k+p, p)
	} else {
		k = -k
		left(0, n-p, k)
		right(n-k-p, n, p)
	}
}

func Shuffle(data sort.Interface) {
	n := data.Len()
	// this does not account for second order swapping, so entropy may vary
	for i, j := range rand.Perm(n) {
		data.Swap(i, j)
	}
}
