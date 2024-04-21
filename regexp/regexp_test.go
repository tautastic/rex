package regexp

import "testing"

func benchmarkRegexpMatch(b *testing.B, re *Regexp, ch rune, times int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < times; n++ {
			re.matchRune(ch)
		}
	}
}

func BenchmarkMatchLiteral100(b *testing.B) {
	b.StopTimer()
	re := &Regexp{Op: OpLiteral, Sym: RuneRange{'a', 'a'}}
	b.StartTimer()
	benchmarkRegexpMatch(b, re, 'a', 100)
}

func BenchmarkMatchLiteral1000(b *testing.B) {
	b.StopTimer()
	re := &Regexp{Op: OpLiteral, Sym: RuneRange{'a', 'a'}}
	b.StartTimer()
	benchmarkRegexpMatch(b, re, 'a', 1000)
}

func BenchmarkMatchLiteral10000(b *testing.B) {
	b.StopTimer()
	re := &Regexp{Op: OpLiteral, Sym: RuneRange{'a', 'a'}}
	b.StartTimer()
	benchmarkRegexpMatch(b, re, 'a', 10000)
}

func BenchmarkMatchClass100(b *testing.B) {
	b.StopTimer()
	re := &Regexp{Op: OpCharClass, Sym: PerlClass['w']}
	b.StartTimer()
	benchmarkRegexpMatch(b, re, 'a', 100)
}

func BenchmarkMatchClass1000(b *testing.B) {
	b.StopTimer()
	re := &Regexp{Op: OpCharClass, Sym: PerlClass['w']}
	b.StartTimer()
	benchmarkRegexpMatch(b, re, 'a', 1000)
}

func BenchmarkMatchClass10000(b *testing.B) {
	b.StopTimer()
	re := &Regexp{Op: OpCharClass, Sym: PerlClass['w']}
	b.StartTimer()
	benchmarkRegexpMatch(b, re, 'a', 10000)
}

func TestMatchLiteral(t *testing.T) {
	for _, test := range []struct {
		lit, ch rune
		want    bool
	}{
		{'a', 'a', true},
		{'a', 'b', false},
		{'a', ' ', false},
		{'a', '?', false},
		{'a', '∑', false},

		{'∑', '∑', true},
		{'∑', 'a', false},
		{'∑', 'b', false},
		{'∑', ' ', false},
		{'∑', '?', false},

		{'ﺙ', 'ﺙ', true},
		{'ﺙ', 'a', false},
		{'ﺙ', 'b', false},
		{'ﺙ', ' ', false},
		{'ﺙ', '?', false},
	} {
		re := &Regexp{Op: OpLiteral, Sym: RuneRange{test.lit, test.lit}}
		got := re.matchRune(test.ch)
		if got != test.want {
			t.Errorf("error:\ngot: %v\nwant: %v", got, test.want)
		}
	}
}

func TestMatchClass(t *testing.T) {
	for _, test := range []struct {
		cls  uint8
		ch   rune
		want bool
	}{
		// Digits
		{'d', '0', true},
		{'d', '1', true},
		{'d', '2', true},
		{'d', '3', true},
		{'d', '4', true},
		{'d', '5', true},
		{'d', '6', true},
		{'d', '7', true},
		{'d', '8', true},
		{'d', '9', true},

		{'d', ' ', false},
		{'d', 'a', false},
		{'d', 'b', false},
		{'d', 'c', false},
		{'d', 'd', false},
		{'d', 'e', false},
		{'d', 'f', false},
		{'d', 'g', false},
		{'d', 'h', false},
		{'d', 'i', false},

		// Non-digits
		{'D', '0', false},
		{'D', '1', false},
		{'D', '2', false},
		{'D', '3', false},
		{'D', '4', false},
		{'D', '5', false},
		{'D', '6', false},
		{'D', '7', false},
		{'D', '8', false},
		{'D', '9', false},

		{'D', ' ', true},
		{'D', 'a', true},
		{'D', 'b', true},
		{'D', 'c', true},
		{'D', 'd', true},
		{'D', 'e', true},
		{'D', 'f', true},
		{'D', 'g', true},
		{'D', 'h', true},
		{'D', 'i', true},

		// Whitespace
		{'s', ' ', true},
		{'s', '\t', true},
		{'s', '\n', true},
		{'s', '\v', true},
		{'s', '\f', true},
		{'s', '\r', true},

		{'s', '0', false},
		{'s', '1', false},
		{'s', '2', false},
		{'s', 'a', false},
		{'s', 'b', false},
		{'s', 'c', false},

		// Non-whitespace
		{'S', ' ', false},
		{'S', '\t', false},
		{'S', '\n', false},
		{'S', '\v', false},
		{'S', '\f', false},
		{'S', '\r', false},

		{'S', '0', true},
		{'S', '1', true},
		{'S', '2', true},
		{'S', 'a', true},
		{'S', 'b', true},
		{'S', 'c', true},

		// Word
		{'w', '_', true},
		{'w', '0', true},
		{'w', 'a', true},
		{'w', 'A', true},
		{'w', 'ﺙ', true},
		{'w', '串', true},

		{'w', ' ', false},
		{'w', '\t', false},
		{'w', '\n', false},
		{'w', '\v', false},
		{'w', '\f', false},
		{'w', '\r', false},

		// Non-Word
		{'W', '_', false},
		{'W', '0', false},
		{'W', 'a', false},
		{'W', 'A', false},
		{'W', 'ﺙ', false},
		{'W', '串', false},

		{'W', ' ', true},
		{'W', '\t', true},
		{'W', '\n', true},
		{'W', '\v', true},
		{'W', '\f', true},
		{'W', '\r', true},
	} {
		re := &Regexp{Op: OpCharClass, Sym: PerlClass[test.cls]}
		got := re.matchRune(test.ch)
		if got != test.want {
			t.Errorf("error: %v in \\%v\ngot: %v\nwant: %v",
				string(test.ch), test.cls, got, test.want)
		}
	}
}
