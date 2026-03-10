package rules

import "strings"

// SpecialCharsRule rejects messages containing special characters or emoji.
type SpecialCharsRule struct{}

func (r *SpecialCharsRule) Name() string { return "special" }

func (r *SpecialCharsRule) Check(msg string) *Diagnostic {
	for _, ch := range msg {
		if isForbiddenPunct(ch) || isEmoji(ch) {
			return &Diagnostic{
				Message:  "log message must not contain special characters or emoji",
				FixedMsg: cleanMsg(msg),
			}
		}
	}
	if strings.HasSuffix(msg, ":") || strings.Contains(msg, "...") {
		return &Diagnostic{
			Message:  "log message must not contain special characters or emoji",
			FixedMsg: cleanMsg(msg),
		}
	}
	return nil
}

func isForbiddenPunct(r rune) bool {
	switch r {
	case '!', '?', ';':
		return true
	}
	return false
}

func isEmoji(r rune) bool {
	return (r >= 0x1F300 && r <= 0x1FAFF) ||
		(r >= 0x2600 && r <= 0x27BF)
}

func cleanMsg(msg string) string {
	var b strings.Builder
	for _, ch := range msg {
		if !isForbiddenPunct(ch) && !isEmoji(ch) {
			b.WriteRune(ch)
		}
	}
	result := strings.ReplaceAll(b.String(), "...", "")
	return strings.TrimRight(result, ": ")
}
