package rules

import "unicode"

// EnglishRule rejects log messages containing non-ASCII letters.
type EnglishRule struct{}

func (r *EnglishRule) Name() string { return "english" }

func (r *EnglishRule) Check(msg string) *Diagnostic {
	for _, ch := range msg {
		if unicode.IsLetter(ch) && ch > unicode.MaxASCII {
			return &Diagnostic{Message: "log message must be in English"}
		}
	}
	return nil
}
