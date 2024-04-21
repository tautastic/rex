package regexp

import (
	"fmt"
	"sort"
	"unicode"

	"github.com/tautastic/rex/utils"
)

type RuneRange []rune

// appendRange returns the result of appending the range lo-hi to the RuneRange rr.
func appendRange(rr RuneRange, lo, hi rune) RuneRange {
	if hi < lo {
		panic(utils.ErrInvalidClassRange)
	}
	length := len(rr)
	for i := 2; i <= 4; i += 2 {
		if length >= i {
			rlo, rhi := rr[length-i], rr[length-i+1]
			if lo <= rhi+1 && rlo <= hi+1 {
				if lo < rlo {
					rr[length-i] = lo
				}
				if hi > rhi {
					rr[length-i+1] = hi
				}
				return rr
			}
		}
	}
	return append(rr, lo, hi)
}

// appendLiteral returns the result of appending the rune ch to the RuneRange rr.
func appendLiteral(rr RuneRange, ch rune) RuneRange {
	return appendRange(rr, ch, ch)
}

// appendClass returns the result of appending the RuneRange rr1 to the RuneRange rr0.
// It assumes rr1 is clean.
func appendClass(rr0 RuneRange, rr1 RuneRange) RuneRange {
	for i := 0; i < len(rr1); i += 2 {
		rr0 = appendRange(rr0, rr1[i], rr1[i+1])
	}
	return rr0
}

// negateClass overwrites rr and returns rr's negation.
// It assumes rr is already clean.
func negateClass(rr RuneRange) RuneRange {
	nextLo := '\u0000'
	w := 0
	for i := 0; i < len(rr); i += 2 {
		lo, hi := rr[i], rr[i+1]
		if nextLo <= lo-1 {
			rr[w] = nextLo
			rr[w+1] = lo - 1
			w += 2
		}
		nextLo = hi + 1
	}
	rr = rr[:w]
	if nextLo <= unicode.MaxRune {
		rr = append(rr, nextLo, unicode.MaxRune)
	}
	return rr
}

// cleanClass sorts the ranges (pairs of elements of r),
// merges them, and eliminates duplicates.
func cleanClass(rrp *RuneRange) RuneRange {
	sort.Sort(ranges{rrp})

	rr := *rrp
	if len(rr) < 2 {
		return rr
	}

	w := 2 // write index
	for i := 2; i < len(rr); i += 2 {
		lo, hi := rr[i], rr[i+1]
		if lo <= rr[w-1]+1 {
			if hi > rr[w-1] {
				rr[w-1] = hi
			}
			continue
		}
		rr[w] = lo
		rr[w+1] = hi
		w += 2
	}

	return rr[:w]
}

type ranges struct {
	rr_ *RuneRange
}

func (ra ranges) Less(i, j int) bool {
	p := *ra.rr_
	i *= 2
	j *= 2
	return p[i] < p[j] || p[i] == p[j] && p[i+1] > p[j+1]
}

func (ra ranges) Len() int {
	return len(*ra.rr_) / 2
}

func (ra ranges) Swap(i, j int) {
	p := *ra.rr_
	i *= 2
	j *= 2
	p[i], p[i+1], p[j], p[j+1] = p[j], p[j+1], p[i], p[i+1]
}

func (rr RuneRange) String() (str string) {
	for i := 0; i < len(rr); i += 2 {
		str += fmt.Sprintf("{%v, %v}: ", rr[i], rr[i+1])
		for num := rr[i]; num < rr[i+1]; num++ {
			str += fmt.Sprintf("%s, ", string(num))
		}
		str += fmt.Sprintf("%s\n", string(rr[i+1]))
	}
	return str
}

func (rr RuneRange) String2() (str string) {
	for i := 0; i < len(rr); i += 2 {
		str += fmt.Sprintf("\t\t0x%x, 0x%x,\n", rr[i], rr[i+1])
	}
	return str
}
