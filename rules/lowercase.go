package rules

import (
	"unicode"
	"unicode/utf8"
)

// LowercaseRule ensures the first letter of a log message is lowercase.
type LowercaseRule struct{}

func (r *LowercaseRule) Name() string { return "lowercase" }

func (r *LowercaseRule) Check(msg string) *Diagnostic {
	if msg == "" {
		return nil
	}
	first, size := utf8.DecodeRuneInString(msg)
	if !unicode.IsUpper(first) {
		return nil
	}
	return &Diagnostic{
		Message:  "log message must start with a lowercase letter",
		FixedMsg: string(unicode.ToLower(first)) + msg[size:],
	}
}
