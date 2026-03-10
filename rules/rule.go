package rules

// Diagnostic describes a rule violation found in a log message.
type Diagnostic struct {
	Message  string
	FixedMsg string
}

// Rule validates a log message string and returns a Diagnostic
// when a violation is detected, or nil if the message is valid.
type Rule interface {
	Name() string
	Check(msg string) *Diagnostic
}

// All returns the default ordered set of rules applied to every log message.
func All() []Rule {
	return []Rule{
		&LowercaseRule{},
		&EnglishRule{},
		&SpecialCharsRule{},
		&SensitiveRule{Keywords: DefaultSensitiveKeywords},
	}
}
