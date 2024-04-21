package regexp

import (
	"strconv"
	"unicode/utf8"

	"github.com/tautastic/rex/syntax"
	"github.com/tautastic/rex/utils"
)

func union(first, second *Regexp) *Regexp {
	return &Regexp{Op: OpAlternate, Sub: []*Regexp{first, second}}
}

func concat(first, second *Regexp) *Regexp {
	re := &Regexp{Op: OpConcat, Sub: []*Regexp{}}
	if first.Op == OpConcat {
		re.Sub = append(re.Sub, second)
		re.Sub = append(re.Sub, first.Sub...)
	} else if second.Op == OpConcat {
		re.Sub = append(re.Sub, first)
		re.Sub = append(re.Sub, second.Sub...)
	} else {
		re.Sub = []*Regexp{first, second}
	}
	return re
}

func repeat(sub0 *Regexp, quant *syntax.Node) *Regexp {
	if len(quant.Sub) != 2 {
		panic(utils.ErrInvalidRepeatSize)
	}
	lower, err := strconv.Atoi(quant.Sub[0].Label)
	upper, err := strconv.Atoi(quant.Sub[1].Label)
	if err != nil || (upper < lower && upper != -1) {
		panic(utils.ErrInvalidRepeatSize)
	}
	return &Regexp{Op: OpRepeat, Min: lower, Max: upper, Sub: []*Regexp{sub0}}
}

func ctrlToRune(ch uint8) rune {
	switch ch {
	case 't':
		return 9
	case 'n':
		return 10
	case 'v':
		return 11
	case 'f':
		return 12
	case 'r':
		return 13
	}
	return endOfText
}

func hexSeqToRune(str string) rune {
	dec, err := strconv.ParseInt(str, 16, 32)
	if err != nil {
		panic("error: invalid hexadecimal number")
	}
	return rune(dec)
}

func classRange(rr RuneRange, node0, node1 *syntax.Node) RuneRange {
	switch node0.Label {
	default:
		panic(utils.ErrUnexpectedSymbol)

	case "Literal":
		lo, _ := utf8.DecodeRuneInString(node0.Sub[0].Label)
		hi, _ := utf8.DecodeRuneInString(node1.Sub[0].Label)
		rr = appendRange(rr, lo, hi)
	case "Control":
		lo := ctrlToRune(node0.Sub[0].Label[0])
		hi := ctrlToRune(node1.Sub[0].Label[0])
		rr = appendRange(rr, lo, hi)
	case "HexSeq":
		lo := hexSeqToRune(node0.Sub[0].Label)
		hi := hexSeqToRune(node1.Sub[0].Label)
		rr = appendRange(rr, lo, hi)
	}
	return rr
}

func fromPerl(ch uint8) *Regexp {
	return &Regexp{Op: OpCharClass, Sym: PerlClass[ch]}
}

func fromControl(ch uint8) *Regexp {
	re := &Regexp{Op: OpLiteral, Sym: RuneRange{}}
	re.Sym = appendLiteral(re.Sym, ctrlToRune(ch))

	return re
}

func fromHexSeq(str string) *Regexp {
	return fromRune(hexSeqToRune(str))
}

func fromUniSeq(str string) *Regexp {
	re := &Regexp{Op: OpCharClass}
	re.Sym = UniClass[str]
	return re
}

func fromClass(children []*syntax.Node) *Regexp {
	re := &Regexp{Op: OpCharClass}
	for i, child := range children {
		switch child.Label {

		default:
			panic(utils.ErrUnexpectedSymbol)

		case "Literal", "^":
			if i != 0 || child.Sub[0].Label != "^" {
				lo, _ := utf8.DecodeRuneInString(child.Sub[0].Label)
				re.Sym = appendLiteral(re.Sym, lo)
			}

		case "Control":
			ctrl := ctrlToRune(child.Sub[0].Label[0])
			re.Sym = appendLiteral(re.Sym, ctrl)

		case "Perl":
			re.Sym = appendClass(re.Sym, PerlClass[child.Sub[0].Label[0]])

		case "HexSeq":
			hseq := hexSeqToRune(child.Sub[0].Label)
			re.Sym = appendLiteral(re.Sym, hseq)

		case "UniSeq":
			useq := UniClass[child.Sub[0].Label]
			re.Sym = appendClass(re.Sym, useq)

		case "ClassRange":
			if len(child.Sub) == 2 {
				re.Sym = classRange(re.Sym, child.Sub[0], child.Sub[1])
			}

		}
	}
	re.Sym = cleanClass(&re.Sym)
	if children[0].Label == "^" {
		re.Sym = negateClass(re.Sym)
	}
	return re
}

func fromRune(r rune) *Regexp {
	re := &Regexp{Op: OpLiteral, Sym: RuneRange{}}
	re.Sym = append(re.Sym, r)
	return re
}

func fromLiteral(str string) *Regexp {
	r, _ := utf8.DecodeRuneInString(str)
	return fromRune(r)
}

func fromAssertion(ch uint8) *Regexp {
	switch ch {
	default:
		panic(utils.ErrInvalidAssertion)
	case '^':
		return &Regexp{Op: OpLineStart}
	case '$':
		return &Regexp{Op: OpLineEnd}
	case 'b':
		return &Regexp{Op: OpWordBoundary}
	case 'B':
		return &Regexp{Op: OpNotWordBoundary}
	}
}

func fromSyntaxTree(root *syntax.Node) *Regexp {
	if root.Sub != nil {
		switch root.Label {
		case "Disjunction":
			term := fromSyntaxTree(root.Sub[0])
			if len(root.Sub) == 2 {
				return union(term, fromSyntaxTree(root.Sub[1]))
			}
			return term

		case "Term":
			factor := fromSyntaxTree(root.Sub[0])
			if len(root.Sub) == 2 {
				return concat(factor, fromSyntaxTree(root.Sub[1]))
			}
			return factor

		case "Factor":
			// Check for assertion missing.
			if root.Sub[0].Label == "Assertion" {
				return fromAssertion(root.Sub[0].Sub[0].Label[0])
			}
			atom := fromSyntaxTree(root.Sub[0])
			if len(root.Sub) == 2 {
				quant := root.Sub[1]
				return repeat(atom, quant)
			}
			return atom

		case "Atom":
			if root.Sub[0].Label == "." {
				return fromPerl('.')
			}
			return fromSyntaxTree(root.Sub[0])

		case "Perl":
			return fromPerl(root.Sub[0].Label[0])

		case "Control":
			return fromControl(root.Sub[0].Label[0])

		case "HexSeq":
			return fromHexSeq(root.Sub[0].Label)

		case "UniSeq":
			return fromUniSeq(root.Sub[0].Label)

		case "Class":
			return fromClass(root.Sub)

		case "Literal":
			return fromLiteral(root.Sub[0].Label)

		}
	}
	panic(utils.ErrUnexpectedSymbol)
}

func FromInfixExp(infixExp string) Regexp {
	if infixExp == "" {
		panic(utils.ErrEmptyRegexPattern)
	}

	re := fromSyntaxTree(syntax.ToSyntaxTree(infixExp))
	return *simplifyRegexp(re)
}
