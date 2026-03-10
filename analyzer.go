package loglinter

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/100bench/loglinter/rules"
)

// Analyzer is the loglinter analysis pass exported for standalone and plugin use.
var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "checks log messages for style and security violations",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// --- configuration flags ---

var (
	flagDisabledRules     string
	flagSensitiveKeywords string
)

func init() {
	Analyzer.Flags.StringVar(&flagDisabledRules, "disabled-rules", "",
		"comma-separated list of rules to disable (lowercase,english,special,sensitive)")
	Analyzer.Flags.StringVar(&flagSensitiveKeywords, "extra-sensitive-keywords", "",
		"comma-separated additional sensitive data keywords")
}

// logMethods maps package import paths to the set of function/method names
// whose first argument is the log message.
var logMethods = map[string]map[string]bool{
	"log": {
		"Print": true, "Printf": true, "Println": true,
		"Fatal": true, "Fatalf": true, "Fatalln": true,
		"Panic": true, "Panicf": true, "Panicln": true,
	},
	"log/slog": {
		"Debug": true, "Info": true, "Warn": true, "Error": true,
	},
	"go.uber.org/zap": {
		"Debug": true, "Info": true, "Warn": true, "Error": true,
		"Fatal": true, "Panic": true, "DPanic": true,
	},
}

func run(pass *analysis.Pass) (any, error) {
	allRules, sensitiveKW := buildRules()
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		pkgPath := resolvePackage(pass, sel)
		if pkgPath == "" {
			return
		}

		methods := logMethods[pkgPath]
		if !methods[sel.Sel.Name] {
			return
		}

		if len(call.Args) == 0 {
			return
		}

		checkArg(pass, call.Args[0], allRules, sensitiveKW)
	})

	return nil, nil
}

// buildRules constructs the active rule set based on flag configuration.
func buildRules() ([]rules.Rule, []string) {
	disabled := toSet(flagDisabledRules)

	var (
		result      []rules.Rule
		sensitiveKW []string
	)

	if !disabled["lowercase"] {
		result = append(result, &rules.LowercaseRule{})
	}
	if !disabled["english"] {
		result = append(result, &rules.EnglishRule{})
	}
	if !disabled["special"] {
		result = append(result, &rules.SpecialCharsRule{})
	}
	if !disabled["sensitive"] {
		sensitiveKW = append([]string{}, rules.DefaultSensitiveKeywords...)
		sensitiveKW = append(sensitiveKW, splitCSV(flagSensitiveKeywords)...)
		result = append(result, &rules.SensitiveRule{Keywords: sensitiveKW})
	}

	return result, sensitiveKW
}

// resolvePackage determines the package import path for a selector expression,
// handling package-level calls (slog.Info), method calls (logger.Info),
// and chained calls (zap.L().Info).
func resolvePackage(pass *analysis.Pass, sel *ast.SelectorExpr) string {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj := pass.TypesInfo.Uses[ident]; obj != nil {
			if pkgName, ok := obj.(*types.PkgName); ok {
				return pkgName.Imported().Path()
			}
		}
	}

	if selection, ok := pass.TypesInfo.Selections[sel]; ok {
		return typePackagePath(selection.Recv())
	}

	return ""
}

func typePackagePath(t types.Type) string {
	switch v := t.(type) {
	case *types.Named:
		if pkg := v.Obj().Pkg(); pkg != nil {
			return pkg.Path()
		}
	case *types.Pointer:
		return typePackagePath(v.Elem())
	}
	return ""
}

// checkArg inspects the first argument of a log call: either a plain string
// literal (all rules) or a concatenation expression (sensitive-data rule only).
func checkArg(pass *analysis.Pass, arg ast.Expr, allRules []rules.Rule, sensitiveKW []string) {
	if lit, ok := asStringLit(arg); ok {
		msg, err := strconv.Unquote(lit.Value)
		if err != nil {
			return
		}
		reportLiteral(pass, lit, msg, allRules)
		return
	}

	if len(sensitiveKW) > 0 {
		checkConcatSensitive(pass, arg, sensitiveKW)
	}
}

func reportLiteral(pass *analysis.Pass, lit *ast.BasicLit, msg string, allRules []rules.Rule) {
	for _, r := range allRules {
		d := r.Check(msg)
		if d == nil {
			continue
		}
		diag := analysis.Diagnostic{
			Pos:     lit.Pos(),
			Message: d.Message,
		}
		if d.FixedMsg != "" {
			diag.SuggestedFixes = []analysis.SuggestedFix{{
				Message: d.Message,
				TextEdits: []analysis.TextEdit{{
					Pos:     lit.Pos(),
					End:     lit.End(),
					NewText: []byte(strconv.Quote(d.FixedMsg)),
				}},
			}}
		}
		pass.Report(diag)
	}
}

// checkConcatSensitive inspects a string concatenation expression for
// sensitive data: string literals containing keyword+separator patterns
// and identifiers whose names match sensitive keywords.
func checkConcatSensitive(pass *analysis.Pass, expr ast.Expr, keywords []string) {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return
	}

	sensitive := &rules.SensitiveRule{Keywords: keywords}
	var found bool

	ast.Inspect(expr, func(n ast.Node) bool {
		if n == nil || found {
			return false
		}
		switch v := n.(type) {
		case *ast.BasicLit:
			if v.Kind == token.STRING {
				if msg, err := strconv.Unquote(v.Value); err == nil && sensitive.Check(msg) != nil {
					found = true
				}
			}
		case *ast.Ident:
			if containsSensitiveKeyword(v.Name, keywords) {
				found = true
			}
		}
		return !found
	})

	if found {
		pass.Reportf(bin.Pos(), "log message must not contain potentially sensitive data")
	}
}

func containsSensitiveKeyword(name string, keywords []string) bool {
	lower := strings.ToLower(name)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func asStringLit(expr ast.Expr) (*ast.BasicLit, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if ok && lit.Kind == token.STRING {
		return lit, true
	}
	return nil, false
}

// --- helpers ---

func toSet(csv string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range splitCSV(csv) {
		m[s] = true
	}
	return m
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			result = append(result, v)
		}
	}
	return result
}
