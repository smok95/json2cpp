package nameutil

import "testing"

func TestSanitizeToCppIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		camelCase bool
		forType   bool
		want      string
	}{
		// Basic camelCase to snake_case
		{"camelCase basic", "userName", false, false, "user_name"},
		{"camelCase with digit", "user2Name", false, false, "user_2_name"},
		{"PascalCase", "UserName", false, false, "user_name"},

		// Consecutive uppercase (acronyms)
		{"HTTP prefix", "HTTPStatus", false, false, "http_status"},
		{"API prefix", "APIKey", false, false, "api_key"},
		{"JSON suffix", "parseJSON", false, false, "parse_json"},
		{"All caps", "URL", false, false, "url"},

		// Mixed cases
		{"isCompress", "isCompress", false, false, "is_compress"},
		{"getHTTPResponse", "getHTTPResponse", false, false, "get_http_response"},
		{"HTTPSConnection", "HTTPSConnection", false, false, "https_connection"},
		{"base64Encode", "base64Encode", false, false, "base_64_encode"},

		// Special characters
		{"dash separator", "user-name", false, false, "user_name"},
		{"dot separator", "user.name", false, false, "user_name"},
		{"space separator", "user name", false, false, "user_name"},
		{"underscore separator", "user_name", false, false, "user_name"},
		{"mixed separators", "api-key.v2_test", false, false, "api_key_v_2_test"},

		// camelCase output (--camelcase flag)
		{"to camelCase basic", "user-name", true, false, "userName"},
		{"to camelCase from snake", "user_name", true, false, "userName"},
		{"preserve camelCase", "userName", true, false, "userName"},
		{"HTTP to camelCase", "HTTPStatus", true, false, "httpStatus"},

		// Type names (PascalCase)
		{"type from snake", "user_config", false, true, "UserConfig"},
		{"type from dash", "user-config", false, true, "UserConfig"},
		{"type from camel", "userConfig", false, true, "UserConfig"},
		{"type HTTP", "HTTPStatus", false, true, "HttpStatus"},

		// Edge cases
		{"starts with digit", "123data", false, false, "_123_data"},
		{"only digits", "123", false, false, "_123"},
		{"only symbols", "---", false, false, "_"},
		{"empty string", "", false, false, "_"},

		// C++ keywords
		{"keyword class", "class", false, false, "class_"},
		{"keyword int", "int", false, false, "int_"},
		{"keyword return", "return", false, false, "return_"},
		{"keyword new", "new", false, false, "new_"},

		// Real-world examples
		{"config key", "max-retry-count", false, false, "max_retry_count"},
		{"API endpoint", "api.v2.users", false, false, "api_v_2_users"},
		{"isEnabled", "isEnabled", false, false, "is_enabled"},
		{"hasPermission", "hasPermission", false, false, "has_permission"},
		{"OAuth2Token", "OAuth2Token", false, false, "o_auth_2_token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeToCppIdentifier(tt.input, tt.camelCase, tt.forType)
			if got != tt.want {
				t.Errorf("SanitizeToCppIdentifier(%q, camelCase=%v, forType=%v) = %q, want %q",
					tt.input, tt.camelCase, tt.forType, got, tt.want)
			}
		})
	}
}

func TestIsCppKeyword(t *testing.T) {
	keywords := []string{"class", "int", "return", "new", "delete", "void", "struct"}
	for _, kw := range keywords {
		t.Run(kw, func(t *testing.T) {
			if !isCppKeyword(kw) {
				t.Errorf("isCppKeyword(%q) = false, want true", kw)
			}
		})
	}

	nonKeywords := []string{"user", "name", "data", "value", "myclass"}
	for _, word := range nonKeywords {
		t.Run("not_"+word, func(t *testing.T) {
			if isCppKeyword(word) {
				t.Errorf("isCppKeyword(%q) = true, want false", word)
			}
		})
	}
}
