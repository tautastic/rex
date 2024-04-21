package utils

// An Error describes a failure to parse a regular expression
// and gives the offending expression.
type Error struct {
	Code ErrorCode
	Expr string
}

func (e *Error) Error() string {
	return "error parsing regexp: " + e.Code.String() + ": `" + e.Expr + "`"
}

// An ErrorCode describes a failure to parse a regular expression.
type ErrorCode string

const (
	ErrInvalidCharClass   ErrorCode = "invalid character class"
	ErrInvalidAssertion   ErrorCode = "invalid assertion"
	ErrRangeWithShorthand ErrorCode = "cannot create a range with shorthand escape sequences"
	ErrInvalidClassRange  ErrorCode = "invalid character class range"
	ErrInvalidRepeatOp    ErrorCode = "invalid repetition operator"
	ErrEmptyRegexPattern  ErrorCode = "regex pattern is empty"
	ErrInvalidRepeatSize  ErrorCode = "invalid repeat count"
	ErrUnexpectedSymbol   ErrorCode = "unexpected symbol"
)

func (e ErrorCode) String() string {
	return string(e)
}
