# rex 🦖

[![GoDoc](https://godoc.org/github.com/tautastic/rex?status.svg)](https://godoc.org/github.com/tautastic/rex)

>Zero dependency, lightweight regex engine written in Go.

## Supported syntax:

### Single characters:
```text
.              any character
[xyz]          character class
[^xyz]         negated character class
\d             Perl character class
\D             negated Perl character class
\p{XX}         Unicode character category XX
\P{XX}         negated Unicode character category XX
```

### Composites:
```text
xy             x followed by y
x|y            x or y (prefer x)
```

### Repetitions:
```text
x*             zero or more x, prefer more
x+             one or more x, prefer more
x?             zero or one x, prefer one
x{n,m}         n or n+1 or ... or m x, prefer more
x{n,}          n or more x, prefer more
x{n}           exactly n x
```

### Character class elements:
```text
x              single character
A-Z            character range (inclusive)
\d             Perl character class
\p{XX}         Unicode character category XX
```

### Named character classes as character class elements:
```text
[\d]           digits (== \d)
[^\d]          not digits (== \D)
[\D]           not digits (== \D)
[^\D]          not not digits (== \d)
[\p{XX}]       Unicode category inside character class (== \p{XX})
[^\p{XX}]      Unicode category inside negated character class (== \P{XX})
```

### Perl character classes (all in Unicode):
```text
\d             digits (== [0-9])
\D             not digits (== [^0-9])
\s             whitespace (== [\t\n\f\r ])
\S             not whitespace (== [^\t\n\f\r ])
\w             word characters (== [0-9A-Za-z_])
\W             not word characters (== [^0-9A-Za-z_])
```

### Unicode character categories:
```text
\p{C}          Other, Any
\p{Cc}         Other, Control
\p{Cf}         Other, Format
\p{Co}         Other, Private Use
\p{Cs}         Other, Surrogate

\p{L}          Letter, Any
\p{Ll}         Letter, Lowercase
\p{Lm}         Letter, Modifier
\p{Lo}         Letter, Other
\p{Lt}         Letter, Titlecase
\p{Lu}         Letter, Uppercase

\p{M}          Mark, Any
\p{Mc}         Mark, Spacing Combining
\p{Me}         Mark, Enclosing
\p{Mn}         Mark, Nonspacing

\p{N}          Number, Any
\p{Nd}         Number, Decimal Digit
\p{Nl}         Number, Letter
\p{No}         Number, Other

\p{P}          Punctuation, Any
\p{Pc}         Punctuation, Connector
\p{Pd}         Punctuation, Dash
\p{Pe}         Punctuation, Close
\p{Po}         Punctuation, Other
\p{Ps}         Punctuation, Open

\p{S}          Symbol, Any
\p{Sc}         Symbol, Currency
\p{Sk}         Symbol, Modifier
\p{Sm}         Symbol, Math
\p{So}         Symbol, Other

\p{Z}          Separator, Any
\p{Zl}         Separator, Line
\p{Zp}         Separator, Paragraph
\p{Zs}         Separator, Space
```
