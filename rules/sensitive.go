package rules

import "strings"

// DefaultSensitiveKeywords is the default set of keywords that indicate
// potentially sensitive data in log messages.
var DefaultSensitiveKeywords = []string{
	"password", "passwd", "token", "secret",
	"api_key", "apikey", "api key",
	"credential", "private_key", "privatekey",
}

// SensitiveRule rejects messages that appear to log sensitive values.
// A keyword followed by ':' or '=' suggests the actual value is being logged.
type SensitiveRule struct {
	Keywords []string
}

func (r *SensitiveRule) Name() string { return "sensitive" }

func (r *SensitiveRule) Check(msg string) *Diagnostic {
	lower := strings.ToLower(msg)
	for _, kw := range r.Keywords {
		if strings.Contains(lower, kw+":") || strings.Contains(lower, kw+"=") {
			return &Diagnostic{
				Message: "log message must not contain potentially sensitive data",
			}
		}
	}
	return nil
}
