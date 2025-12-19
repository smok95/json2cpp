package nameutil

// package nameutil

import (
	"strings"
	"unicode"
)

// SanitizeToCppIdentifier converts arbitrary JSON keys into safe C++ identifiers.
// - camelCase: when true and not forType, produces lowerCamelCase for fields.
// - forType: when true produces PascalCase suitable for type/struct names.
// This treats any non-alphanumeric rune as a word boundary (e.g. '-', ':', '.', ' ').
func SanitizeToCppIdentifier(key string, camelCase bool, forType bool) string {
	words := make([]string, 0)
	var cur []rune

	// split into words on non-alnum (treat '_' also as separator)
	for _, r := range key {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cur = append(cur, r)
		} else {
			if len(cur) > 0 {
				words = append(words, string(cur))
				cur = cur[:0]
			}
		}
	}
	if len(cur) > 0 {
		words = append(words, string(cur))
	}

	// fallback: if no word (e.g. key consisted only of symbols), use underscore
	if len(words) == 0 {
		return "_"
	}

	// Helper to title-case a word (preserve digits)
	title := func(s string) string {
		if s == "" {
			return s
		}
		// Use ToLower then Uppercase first rune for deterministic behavior
		lower := strings.ToLower(s)
		rs := []rune(lower)
		rs[0] = unicode.ToUpper(rs[0])
		return string(rs)
	}

	var result string
	if forType {
		// PascalCase for types
		out := make([]string, 0, len(words))
		for _, w := range words {
			out = append(out, title(w))
		}
		result = strings.Join(out, "")
	} else if camelCase {
		// lowerCamelCase for fields
		out := make([]string, 0, len(words))
		for i, w := range words {
			if i == 0 {
				out = append(out, strings.ToLower(w))
			} else {
				out = append(out, title(w))
			}
		}
		result = strings.Join(out, "")
	} else {
		// snake_case
		out := make([]string, 0, len(words))
		for _, w := range words {
			out = append(out, strings.ToLower(w))
		}
		result = strings.Join(out, "_")
	}

	// If starts with digit, prefix underscore
	if len(result) > 0 {
		r0 := []rune(result)[0]
		if unicode.IsDigit(r0) {
			result = "_" + result
		}
	}

	// if result is a C++ keyword, append underscore
	if isCppKeyword(result) {
		result = result + "_"
	}

	// final safety: ensure not empty
	if result == "" {
		return "_"
	}

	return result
}

func isCppKeyword(s string) bool {
	keywords := []string{
		"alignas", "alignof", "and", "and_eq", "asm", "auto", "bitand", "bitor", "bool", "break",
		"case", "catch", "char", "char16_t", "char32_t", "class", "compl", "const", "constexpr",
		"const_cast", "continue", "decltype", "default", "delete", "do", "double", "dynamic_cast",
		"else", "enum", "explicit", "export", "extern", "false", "float", "for", "friend", "goto",
		"if", "inline", "int", "long", "mutable", "namespace", "new", "noexcept", "not", "not_eq",
		"nullptr", "operator", "or", "or_eq", "private", "protected", "public", "register",
		"reinterpret_cast", "return", "short", "signed", "sizeof", "static", "static_assert",
		"static_cast", "struct", "switch", "template", "this", "thread_local", "throw", "true",
		"try", "typedef", "typeid", "typename", "union", "unsigned", "using", "virtual", "void",
		"volatile", "wchar_t", "while", "xor", "xor_eq",
	}
	lo := strings.ToLower(s)
	for _, kw := range keywords {
		if lo == kw {
			return true
		}
	}
	return false
}
