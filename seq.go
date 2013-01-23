// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil

import (
	"math/rand"
	"sort"
)

// Reverse inverts the current order of the provided data.
func Reverse(data sort.Interface) {
	n := data.Len()
	for i := 0; i < n/2; i++ {
		data.Swap(i, n-i-1)
	}
}

// Rotate cycles data by d moves to the right.
// The d rightmost block of items will be shifted to the front.
// If d is negative, the shift will be leftward.
func Rotate(data sort.Interface, d int) {
	k := data.Len()
	d = (k + d) % k
	Skew(data, 0, d, k-d)
}

// Skew slides a group of k consecutive elements from index i to index j.
// i and j respectively represent the source and destination indices of the
// group's minimum-indexed edge. If j > i, the group will slide toward larger
// indices, while if j < i, the group will slide toward smaller indices.
//
// i, j, and k should all be non-negative integers within the range of
// [0,n), where n == data.Len().
func Skew(data sort.Interface, i, j, k int) {
	if k == 0 || i == j {
		return
	} else if j < i {
		i, j, k = j, j+k, i-j
	}
	if j-i < k {
		// if the block size is larger than the delta...
		p := k / 2
		q := k - p
		Skew(data, i+p, j+p, q)
		Skew(data, i, j, p)
	} else if p := (j - i) % k; p != 0 {
		// if the delta is not evenly divisible by the block size...
		Skew(data, i, j-p, k)
		Skew(data, j-p, j, k)
	} else {
		for ; i < j; i++ {
			data.Swap(i, i+k)
		}
	}
}

// Shuffle sorts data randomly.
func Shuffle(data sort.Interface) {
	n := data.Len()
	// this does not account for second order swapping, so entropy may vary
	for i, j := range rand.Perm(n) {
		data.Swap(i, j)
	}
}
