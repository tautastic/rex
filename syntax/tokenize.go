package syntax

import (
	"unicode/utf8"

	"github.com/tautastic/rex/utils"
)

var pattern string
var pos int

const endOfText rune = -1

func peek(n int) rune {
	if pos < len(pattern)+n {
		c := pattern[pos+n]
		if c < utf8.RuneSelf {
			return rune(c)
		}
		r, _ := utf8.DecodeRuneInString(pattern[pos+n:])
		return r
	}
	return endOfText
}

func match(ch rune) {
	if peek(0) != ch {
		panic(utils.ErrUnexpectedSymbol)
	}
	pos += utf8.RuneLen(ch)
}

func next(n int) rune {
	var ch rune
	for i := 0; i <= n; i++ {
		ch = peek(i)
		match(ch)
	}
	return ch
}

func stripSpace() {
	for peek(0) == ' ' {
		next(0)
	}
}

func disjunction() (node *Node) {
	node = &Node{Label: "Disjunction", Sub: nil}
	trm := term()
	if peek(0) == '|' {
		match('|')
		dis := disjunction()
		node.Sub = []*Node{trm, dis}
	} else {
		node.Sub = []*Node{trm}
	}
	return node
}

func term() (node *Node) {
	node = &Node{Label: "Term", Sub: nil}
	factr := factor()
	if peek(0) != endOfText &&
		!utils.IsAnyOf(peek(0), []rune{'|', ')', ']', '}'}) {
		trm := term()
		node.Sub = []*Node{factr, trm}
	} else {
		node.Sub = []*Node{factr}
	}
	return node
}

func factor() (node *Node) {
	node = &Node{Label: "Factor", Sub: nil}
	if utils.IsAnyOf(peek(0), []rune{'^', '$'}) {
		asr := assertion()
		node.Sub = []*Node{asr}
	} else if peek(0) == '\\' && (peek(1) == 'b' || peek(1) == 'B') {
		match('\\')
		asr := assertion()
		node.Sub = []*Node{asr}
	} else {
		atm := atom()
		if utils.IsAnyOf(peek(0), []rune{'*', '+', '?', '{'}) {
			qnt := quantifier()
			node.Sub = []*Node{atm, qnt}
		} else {
			node.Sub = []*Node{atm}
		}
	}
	return node
}

func assertion() (node *Node) {
	node = &Node{Label: "Assertion",
		Sub: []*Node{{Label: string(next(0))}}}
	return node
}

func quantifier() (node *Node) {
	node = &Node{Label: "Quantifier", Sub: nil}
	switch next(0) {
	default:
		panic(utils.ErrInvalidRepeatOp)
	case '*':
		// Zero or more
		node.Sub = []*Node{{Label: "0"}, {Label: "-1"}}
	case '+':
		// One or more
		node.Sub = []*Node{{Label: "1"}, {Label: "-1"}}
	case '?':
		// Zero or one
		node.Sub = []*Node{{Label: "0"}, {Label: "1"}}
	case '{':
		stripSpace()
		lower := decimalDigits()
		stripSpace()
		if peek(0) == ',' {
			match(',')
			stripSpace()
			if peek(0) == '}' {
				node.Sub = []*Node{lower, {Label: "-1"}}
			} else {
				upper := decimalDigits()
				node.Sub = []*Node{lower, upper}
			}
		} else {
			node.Sub = []*Node{lower, lower}
		}
		stripSpace()
		match('}')
	}
	return node
}

func atom() (node *Node) {
	node = &Node{Label: "Atom", Sub: nil}
	switch peek(0) {
	default:
		lit := anyLiteralExcept([]rune{
			'^', '$', '\\', '.', '*', '+', '?',
			'(', ')', '[', ']', '{', '}', '|'}, utils.ErrUnexpectedSymbol)
		node.Sub = []*Node{lit}
	case '.':
		match('.')
		node.Sub = []*Node{{Label: "."}}

	case '\\':
		match('\\')
		esc := atomEscape()
		node.Sub = []*Node{esc}

	case '[':
		cls := characterClass()
		node.Sub = []*Node{cls}

	case '(':
		match('(')
		dis := disjunction()
		match(')')
		node.Sub = []*Node{dis}

	}
	return node
}

func atomEscape() (node *Node) {
	switch peek(0) {
	case 'f', 'n', 'r', 't', 'v':
		node = &Node{Label: "Control",
			Sub: []*Node{{Label: string(next(0))}}}
	case 'd', 'D', 's', 'S', 'w', 'W':
		node = &Node{Label: "Perl",
			Sub: []*Node{{Label: string(next(0))}}}
	case 'x':
		match('x')
		hex := hexSequence()
		node = &Node{Label: "HexSeq", Sub: []*Node{hex}}
	case 'p', 'P':
		uni := unicodeSequence()
		node = &Node{Label: "UniSeq", Sub: []*Node{uni}}
	}
	return node
}

func hexSequence() (node *Node) {
	match('{')
	node = &Node{Label: ""}
	for utils.IsAnyOf(peek(0), []rune{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	}) {
		node.Label += string(next(0))
	}
	match('}')
	return node
}

func unicodeSequence() (node *Node) {
	if next(0) == 'p' {
		node = &Node{Label: ""}
	} else {
		node = &Node{Label: "^"}
	}
	match('{')
	if utils.IsAnyOf(peek(0), []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}) {
		node.Label += string(next(0))
	}
	if utils.IsAnyOf(peek(0), []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}) {
		node.Label += string(next(0))
	}
	match('}')
	return node
}

func characterClass() (node *Node) {
	match('[')
	node = &Node{Label: "Class", Sub: nil}
	for peek(0) != ']' {
		clr := classRange()
		if len(clr.Sub) == 1 {
			clr = clr.Sub[0]
		}
		node.Sub = append(node.Sub, clr)
	}
	match(']')
	return node
}

func classRange() (node *Node) {
	node = &Node{Label: "ClassRange", Sub: nil}
	cla0 := classAtom()
	node.Sub = append(node.Sub, cla0)
	if peek(0) == '-' {
		match('-')
		cla1 := classAtom()
		if cla0.Label == "Perl" || cla0.Label == "UniSeq" ||
			cla1.Label == "Perl" || cla1.Label == "UniSeq" {
			panic(utils.ErrRangeWithShorthand)
		}
		node.Sub = append(node.Sub, cla1)
	}
	return node
}

func classAtom() (node *Node) {
	if peek(0) == '\\' {
		match('\\')
		node = atomEscape()
	} else {
		node = anyLiteralExcept([]rune{'\\', ']', '-'}, utils.ErrInvalidCharClass)
	}
	return node
}

func decimalDigits() (node *Node) {
	if !utils.IsAnyOf(peek(0),
		[]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}) {
		panic(utils.ErrInvalidRepeatSize)
	}
	node = &Node{Label: ""}
	for utils.IsAnyOf(peek(0),
		[]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}) {
		node.Label += string(next(0))
	}
	return node
}

func anyLiteralExcept(rs []rune, err utils.ErrorCode) (node *Node) {
	if utils.IsAnyOf(peek(0), rs) {
		panic(err)
	}
	node = &Node{Label: "Literal",
		Sub: []*Node{{Label: string(next(0))}}}
	return node
}

func ToSyntaxTree(regex string) *Node {
	pattern = regex
	pos = 0

	return disjunction()
}
