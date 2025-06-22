// string_test.go - Tests for string flag value implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"testing"
)

func TestStringValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "empty string",
			value:    "",
			expected: "",
		},
		{
			name:     "simple string",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "string with spaces",
			value:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with special characters",
			value:    "hello@world.com",
			expected: "hello@world.com",
		},
		{
			name:     "string with unicode",
			value:    "héllo wörld",
			expected: "héllo wörld",
		},
		{
			name:     "string with newlines",
			value:    "line1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "string with tabs",
			value:    "col1\tcol2",
			expected: "col1\tcol2",
		},
		{
			name:     "numeric string",
			value:    "12345",
			expected: "12345",
		},
		{
			name:     "path string",
			value:    "/path/to/file.txt",
			expected: "/path/to/file.txt",
		},
		{
			name:     "json-like string",
			value:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := stringValue{
				value: &tt.value,
			}

			result := v.String()
			if result != tt.expected {
				t.Errorf("stringValue.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStringValue_String_PointerModification(t *testing.T) {
	// Test that String() reflects changes to the underlying value
	value := "initial"
	v := stringValue{
		value: &value,
	}

	// Check initial value
	if got := v.String(); got != "initial" {
		t.Errorf("stringValue.String() = %q, want %q", got, "initial")
	}

	// Modify the underlying value
	value = "modified"

	// Check that String() reflects the change
	if got := v.String(); got != "modified" {
		t.Errorf("stringValue.String() after modification = %q, want %q", got, "modified")
	}
}
