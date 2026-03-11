package rules_test

import (
	"testing"

	"github.com/100bench/loglinter/rules"
)

func TestLowercaseRule(t *testing.T) {
	r := &rules.LowercaseRule{}
	cases := []struct {
		msg     string
		wantErr bool
		wantFix string
	}{
		{"starting server", false, ""},
		{"Starting server", true, "starting server"},
		{"", false, ""},
		{"123 numeric start", false, ""},
		{"ALLCAPS", true, "aLLCAPS"},
	}
	for _, tc := range cases {
		d := r.Check(tc.msg)
		if (d != nil) != tc.wantErr {
			t.Errorf("Check(%q): got err=%v, want err=%v", tc.msg, d != nil, tc.wantErr)
		}
		if d != nil && d.FixedMsg != tc.wantFix {
			t.Errorf("Check(%q): got fix=%q, want fix=%q", tc.msg, d.FixedMsg, tc.wantFix)
		}
	}
}

func TestEnglishRule(t *testing.T) {
	r := &rules.EnglishRule{}
	cases := []struct {
		msg     string
		wantErr bool
	}{
		{"starting server", false},
		{"запуск сервера", true},
		{"mixed english и русский", true},
		{"numbers 123 and symbols +-=", false},
		{"", false},
	}
	for _, tc := range cases {
		d := r.Check(tc.msg)
		if (d != nil) != tc.wantErr {
			t.Errorf("Check(%q): got err=%v, want err=%v", tc.msg, d != nil, tc.wantErr)
		}
	}
}

func TestSpecialCharsRule(t *testing.T) {
	r := &rules.SpecialCharsRule{}
	cases := []struct {
		msg     string
		wantErr bool
		wantFix string
	}{
		{"server started", false, ""},
		{"server started!", true, "server started"},
		{"is this working?", true, "is this working"},
		{"wait for it...", true, "wait for it"},
		{"done;", true, "done"},
		{"starting:", true, "starting"},
		{"hello 🎉", true, "hello"},
		{"clean message", false, ""},
		{"", false, ""},
	}
	for _, tc := range cases {
		d := r.Check(tc.msg)
		if (d != nil) != tc.wantErr {
			t.Errorf("Check(%q): got err=%v, want err=%v", tc.msg, d != nil, tc.wantErr)
		}
		if d != nil && d.FixedMsg != tc.wantFix {
			t.Errorf("Check(%q): got fix=%q, want fix=%q", tc.msg, d.FixedMsg, tc.wantFix)
		}
	}
}

func TestSensitiveRule(t *testing.T) {
	r := &rules.SensitiveRule{Keywords: rules.DefaultSensitiveKeywords}
	cases := []struct {
		msg     string
		wantErr bool
	}{
		{"password: admin", true},
		{"token: abc", true},
		{"api_key=xxx", true},
		{"secret=val", true},
		{"credential: test", true},
		{"token validated", false},
		{"user authenticated", false},
		{"password reset successful", false},
		{"starting server", false},
		{"", false},
	}
	for _, tc := range cases {
		d := r.Check(tc.msg)
		if (d != nil) != tc.wantErr {
			t.Errorf("Check(%q): got err=%v, want err=%v", tc.msg, d != nil, tc.wantErr)
		}
	}
}
