package regexp

import (
	"unicode/utf8"
)

const endOfText rune = -1

func step(str string, pos int) (rune, int) {
	if -1 < pos && pos < len(str) {
		c := str[pos]
		if c < utf8.RuneSelf {
			return rune(c), 1
		}
		return utf8.DecodeRuneInString(str[pos:])
	}
	return endOfText, 0
}

// doMatch reports whether str matches the regexp.
func (re *Regexp) doMatch(str string) bool {
	return re.doOnePass(str, 0, nil) != nil
}

// doOnePass finds the leftmost match in the input, appends its position
// to matches and returns matches.
//
// nil is returned if no matches are found and non-nil if matches are found.
func (re *Regexp) doOnePass(str string, pos0 int, matches []int) []int {

	didMatch := false
	pos1 := pos0
	rePos := 0

	r0, r1, w0, w1 := endOfText, endOfText, 0, 0
	r0, w0 = step(str, pos1)

	for {
		if len(re.Sub) <= rePos {
			goto Return
		}
		re1 := re.Sub[rePos]

		switch re1.Op {
		case OpAccept:
			didMatch = true
			goto Return

		case OpWordBoundary:
			prevRune, _ := step(str, pos1-1)
			if isWordChar(prevRune) != isWordChar(r0) {
				rePos++
				continue
			} else {
				goto Return
			}

		case OpNotWordBoundary:
			prevRune, _ := step(str, pos1-1)
			if isWordChar(prevRune) == isWordChar(r0) {
				rePos++
				continue
			} else {
				goto Return
			}

		case OpLiteral, OpCharClass:
			if !re1.matchRune(r0) {
				goto Return
			}
			rePos++

		case OpRepeat:
			rMin, rMax, lastPos, mCount := re1.Min, re1.Max, pos1, 0
			if rMax < 0 {
				rMax = len(str)
			}
			re1.Sub = append(re1.Sub, &Regexp{Op: OpAccept})
			for {
				tmp := re1.doOnePass(str, lastPos, matches)
				if rMax <= mCount || len(tmp) == 0 {
					break
				}
				mCount++
				lastPos = tmp[len(tmp)-1]
			}
			re1.Sub = re1.Sub[:len(re1.Sub)-1]
			if rMin <= mCount && mCount <= rMax {
				if lastPos != pos1 {
					w0 = lastPos - pos1
				}
				rePos++
			} else {
				goto Return
			}

		case OpConcat:
			if re1.Sub[len(re1.Sub)-1].Op != OpAccept {
				re1.Sub = append(re1.Sub, &Regexp{Op: OpAccept})
			}
			tmp := re1.doOnePass(str, pos1, matches)
			if tmp != nil {
				w0 = tmp[1] - pos1
				rePos++
			} else {
				goto Return
			}

		case OpAlternate:
			alt0, alt1 := re1.Sub[0], re1.Sub[1]
			if alt0.Sub[len(alt0.Sub)-1].Op != OpAccept {
				alt0.Sub = append(alt0.Sub, &Regexp{Op: OpAccept})
			}
			tmp0 := alt0.doOnePass(str, pos1, matches)
			if tmp0 != nil {
				w0 = tmp0[1] - pos1
				rePos++
			} else {
				if alt1.Sub[len(alt1.Sub)-1].Op != OpAccept {
					alt1.Sub = append(alt1.Sub, &Regexp{Op: OpAccept})
				}
				tmp1 := alt1.doOnePass(str, pos1, matches)
				if tmp1 != nil {
					w0 = tmp1[1] - pos1
					rePos++
				} else {
					goto Return
				}
			}
		}

		if w0 == 0 {
			break
		}
		if r0 != endOfText {
			r1, w1 = step(str, pos1+w0)
		}
		pos1 += w0
		r0, w0 = r1, w1
	}

Return:
	if didMatch {
		matches = append(matches, pos0, pos1)
		return matches
	}
	return nil
}
