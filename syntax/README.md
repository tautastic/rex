# rex - syntax

The rex parser implements the following grammar:
```text
<Pattern> ::=
	<Disjunction>

<Disjunction> ::=
	<Term>
	<Term> | <Disjunction>

<Term> ::=
	<Factor>
	<Factor> <Term>

<Factor> ::=
	<Assertion>
	<Atom>
	<Atom> <Quantifier>

<Assertion> ::=
	^
	$
	\ b
	\ B

<Quantifier> ::=
	*
	+
	?
	{ <DecimalDigits> }
	{ <DecimalDigits> , }
	{ <DecimalDigits> , <DecimalDigits> }

<Atom> ::=
	.
	\ <AtomEscape>
	<Class>
	( <Disjunction> )
	any character but not one of ^ $ \ . * + ? ( ) [ ] { } |

<AtomEscape> ::=
	<Control>
	<Perl>
	<HexSeq>
	<UniSeq>

<Control> ::= one of
	f n r t v

<Perl> ::= one of
	d D s S w W

<Class> ::=
	[ <ClassRange> ]
	[ ^ <ClassRange> ]

<ClassRange> ::=
	[empty]
	<ClassAtom>
	<ClassAtom> - <ClassAtom> <ClassRange>

<ClassAtom> ::=
	\ <AtomEscape>
	any character but not one of \ or ] or -

```