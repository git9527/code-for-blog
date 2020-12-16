package main

import (
	"fmt"
	"math/big"
)

// Includes several implementations of CRT (Chinese Remainder Theorem).
// All functions are to be used as follows:
//
//   result := f(a, n)
//
// Where 'a' is the slice of residues and 'n' the slice of corresponding moduli.
// These represent the set of congruences:
//
//    x = a[0] (mod n[0])
//    x = a[1] (mod n[1])
//    ...
//
// Restrictions: len(a) == len(n). Also, all n have to be coprime; otherwise
// the result is not guaranteed to be correct.

func crtSearch(a, n []int64) int64 {
	var N int64 = 1
	for _, nk := range n {
		N *= nk
	}

search:
	for i := int64(0); i < N; i++ {
		for k := 0; k < len(n); k++ {
			if i%n[k] != a[k] {
				continue search
			}
		}

		return i
	}

	return -1
}

func crtSieve(a, n []int64) int64 {
	var N int64 = 1
	for _, nk := range n {
		N *= nk
	}

	base := a[0]
	incr := n[0]

nextBase:
	for i := 1; i < len(a); i++ {
		for candidate := base; candidate < N; candidate += incr {
			if candidate%n[i] == a[i] {
				base = candidate
				incr = incr * n[i]
				continue nextBase
			}
		}
		// Inner loop exited without finding candidate
		return -1
	}
	return base
}

func crtSieveBig(a, n []*big.Int) *big.Int {
	// Compute N: product(n[...])
	N := new(big.Int).Set(n[0])
	for _, nk := range n[1:] {
		N.Mul(N, nk)
	}

	base := a[0]
	incr := n[0]

nextBase:
	for i := 1; i < len(a); i++ {
		for candidate := base; candidate.Cmp(N) < 0; candidate.Add(candidate, incr) {
			mod := new(big.Int).Mod(candidate, n[i])
			if mod.Cmp(a[i]) == 0 {
				base = candidate
				incr.Mul(incr, n[i])
				continue nextBase
			}
		}
		// Inner loop exited without finding candidate
		return big.NewInt(-1)
	}
	return base
}

func crtConstructBig(a, n []*big.Int) *big.Int {
	// Compute N: product(n[...])
	N := new(big.Int).Set(n[0])
	for _, nk := range n[1:] {
		N.Mul(N, nk)
	}

	// x is the accumulated answer.
	x := new(big.Int)

	for k, nk := range n {
		// Nk = N/nk
		Nk := new(big.Int).Div(N, nk)

		// N'k (Nkp) is the multiplicative inverse of Nk modulo nk.
		Nkp := new(big.Int)
		if Nkp.ModInverse(Nk, nk) == nil {
			return big.NewInt(-1)
		}

		// x += ak*Nk*Nkp
		x.Add(x, Nkp.Mul(a[k], Nkp.Mul(Nkp, Nk)))
	}
	return x.Mod(x, N)
}

func main() {
	fmt.Println("run tests")
}
