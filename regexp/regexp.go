package regexp

import "unicode"

// An Op is a single regular expression operator.
type Op uint8

const noMatch = -1

const (
	OpLiteral         Op = 1 + iota // matches a single rune
	OpCharClass                     // matches Runes interpreted as range pair list
	OpRepeat                        // matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
	OpConcat                        // matches concatenation of Subs
	OpAlternate                     // matches alternation of Subs
	OpLineStart                     // asserts position at the start of a line
	OpLineEnd                       // asserts position at the end of a line
	OpWordBoundary                  // asserts position at a word boundary
	OpNotWordBoundary               // asserts position where \b does not match
	OpAccept
)

// A Regexp is a node in a regular expression syntax tree.
type Regexp struct {
	Op  Op
	Min int
	Max int
	Sym RuneRange
	Sub []*Regexp
}

var iFlag bool

// matchRune checks whether the expression matches (and consumes) r.
func (re *Regexp) matchRune(r rune) bool {
	return re.matchRunePos(r) != noMatch
}

// isWordChar checks whether the rune r is a word character.
func isWordChar(r rune) bool {
	if r == endOfText {
		return false
	}
	re := &Regexp{Op: OpCharClass, Sym: PerlClass['w']}
	return re.matchRunePos(r) != noMatch
}

// matchRunePos checks whether the expression matches (and consumes) r.
// If so, matchRunePos returns the index of the matching rune pair.
// If not, matchRunePos returns -1.
func (re *Regexp) matchRunePos(ch rune) int {
	switch len(re.Sym) {
	case 0:
		return noMatch
	case 1:
		if ch == re.Sym[0] || iFlag &&
			unicode.SimpleFold(ch) == re.Sym[0] {
			return 0
		}
		return noMatch
	case 2:
		if (re.Sym[0] <= ch && ch <= re.Sym[1]) || iFlag &&
			re.Sym[0] <= unicode.SimpleFold(ch) &&
			unicode.SimpleFold(ch) <= re.Sym[1] {
			return 0
		}
		return noMatch
	case 4, 6, 8:
		// Linear search for a few pairs.
		for j := 0; j < len(re.Sym); j += 2 {
			if ch < re.Sym[j] || iFlag &&
				unicode.SimpleFold(ch) < re.Sym[j] {
				return noMatch
			}
			if ch <= re.Sym[j+1] || iFlag &&
				unicode.SimpleFold(ch) <= re.Sym[j+1] {
				return j / 2
			}
		}
		return noMatch
	}
	// Otherwise binary search.
	lo := 0
	hi := len(re.Sym) / 2
	for lo < hi {
		m := lo + (hi-lo)/2
		c := re.Sym[2*m]
		if c <= ch || iFlag && c <= unicode.SimpleFold(ch) {
			if c <= ch && ch <= re.Sym[2*m+1] ||
				c <= unicode.SimpleFold(ch) &&
					unicode.SimpleFold(ch) <= re.Sym[2*m+1] {
				return m
			}
			lo = m + 1
		} else {
			hi = m
		}
	}
	return noMatch
}

func (re *Regexp) allMatches(str string, n int, deliver func([]int)) {
	end := len(str)

	for pos, i := 0, 0; i < n && pos <= end; {
		matches := re.doOnePass(str, pos, nil)
		if len(matches) == 0 {
			// No match found, move on.
			pos++
			continue
		}
		if matches[0] == matches[1] {
			_, width := step(str, pos)
			if width > 0 {
				pos += width
			} else {
				pos = end + 1
			}
		} else {
			pos = matches[1]
		}

		deliver(matches)
		i++
	}
}

func (re *Regexp) MatchString(str string, i bool) bool {
	iFlag = i
	return re.doMatch(str)
}

func (re *Regexp) FindString(str string, i bool) string {
	iFlag = i
	var dstCap [2]int
	a := re.doOnePass(str, 0, dstCap[:0])
	if a == nil {
		return ""
	}
	return str[a[0]:a[1]]
}

func (re *Regexp) FindStringIndex(str string, i bool) []int {
	iFlag = i
	a := re.doOnePass(str, 0, nil)
	if a == nil {
		return nil
	}
	return a[0:2]
}

func (re *Regexp) FindAllString(str string, n int, i bool) []string {
	iFlag = i
	if n < 0 {
		n = len(str) + 1
	}
	var result []string
	re.allMatches(str, n, func(match []int) {
		if result == nil {
			result = make([]string, 0, 10)
		}
		result = append(result, str[match[0]:match[1]])
	})
	return result
}

func (re *Regexp) FindAllStringIndex(str string, n int, i bool) [][]int {
	iFlag = i
	if n < 0 {
		n = len(str) + 1
	}
	var result [][]int
	re.allMatches(str, n,
		func(match []int) {
			if result == nil {
				result = make([][]int, 0, 10)
			}
			result = append(result, match[0:2])
		})
	return result
}
