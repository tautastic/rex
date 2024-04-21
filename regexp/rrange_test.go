package regexp

import (
	"testing"

	"github.com/tautastic/rex/utils"
)

func benchmarkAppendLiteral(b *testing.B, rr *RuneRange, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			appendLiteral(*rr, rune(i))
		}
	}
}

func BenchmarkAppendLiteral100(b *testing.B) {
	b.StopTimer()
	size := 100
	rr := RuneRange{}
	b.StartTimer()
	benchmarkAppendLiteral(b, &rr, size)
}

func BenchmarkAppendLiteral1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	rr := RuneRange{}
	b.StartTimer()
	benchmarkAppendLiteral(b, &rr, size)
}

func BenchmarkAppendLiteral10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	rr := RuneRange{}
	b.StartTimer()
	benchmarkAppendLiteral(b, &rr, size)
}

func TestAppendRange(t *testing.T) {
	var got RuneRange
	for _, test := range []struct {
		lo, hi rune
		want   RuneRange
	}{
		{'0', '9', RuneRange{48, 57}},
		{'1', '3', RuneRange{48, 57}},
		{'4', '5', RuneRange{48, 57}},
		{'0', '9', RuneRange{48, 57}},
		{'6', '9', RuneRange{48, 57}},

		{'A', 'Z', RuneRange{48, 57, 65, 90}},
		{'G', 'O', RuneRange{48, 57, 65, 90}},
		{'B', 'F', RuneRange{48, 57, 65, 90}},
		{'X', 'Y', RuneRange{48, 57, 65, 90}},
		{'L', 'M', RuneRange{48, 57, 65, 90}},

		{'a', 'z', RuneRange{48, 57, 65, 90, 97, 122}},
		{'g', 'o', RuneRange{48, 57, 65, 90, 97, 122}},
		{'b', 'f', RuneRange{48, 57, 65, 90, 97, 122}},
		{'x', 'y', RuneRange{48, 57, 65, 90, 97, 122}},
		{'l', 'm', RuneRange{48, 57, 65, 90, 97, 122}},
	} {
		got = appendRange(got, test.lo, test.hi)
		if !utils.Equal(got, test.want) {
			t.Errorf("error:\ngot: %v\nwant: %v", got, test.want)
		}
	}
}

func TestAppendLiteral(t *testing.T) {
	for _, test := range []struct {
		literal  rune
		rr, want RuneRange
	}{
		{'0', RuneRange{}, RuneRange{48, 48}},
		{'0', RuneRange{48, 48}, RuneRange{48, 48}},
		{'2', RuneRange{48, 48}, RuneRange{48, 48, 50, 50}},
		{'4', RuneRange{48, 48, 50, 50}, RuneRange{48, 48, 50, 50, 52, 52}},
		{'1', RuneRange{48, 48, 50, 50, 52, 52}, RuneRange{48, 50, 52, 52}},
		{'3', RuneRange{48, 50, 52, 52}, RuneRange{48, 52}},
	} {
		got := appendLiteral(test.rr, test.literal)
		got = cleanClass(&got)
		if !utils.Equal(got, test.want) {
			t.Errorf("error:\ngot:  %v\nwant: %v", got, test.want)
		}
	}
}

func TestCharClasses(t *testing.T) {
	for _, test := range []struct {
		c0, c1 uint8
	}{
		{'s', 'S'},
		{'d', 'D'},
		{'w', 'W'},
	} {
		r0, r1 := PerlClass[test.c0], PerlClass[test.c1]
		got := appendClass(r0, r1)
		got = cleanClass(&got)
		if got[0] != 0x0 || got[1] != 0x10ffff {
			t.Errorf("error:\ngot:  %v\nwant: [0x0, 0x10ffff]", got)
		}
	}
}
