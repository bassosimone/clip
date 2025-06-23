// textwrap_test.go - Tests for text wrapping functionality
// SPDX-License-Identifier: GPL-3.0-or-later

package textwrap

import (
	"strings"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		indent   string
		expected string
	}{
		{
			name:     "empty text",
			text:     "",
			width:    20,
			indent:   "  ",
			expected: "",
		},

		{
			name:     "single word",
			text:     "hello",
			width:    20,
			indent:   "  ",
			expected: "  hello",
		},

		{
			name:     "multiple words fitting on one line",
			text:     "hello world",
			width:    20,
			indent:   "  ",
			expected: "  hello world",
		},

		{
			name:     "text requiring wrap",
			text:     "this is a long sentence that needs wrapping",
			width:    20,
			indent:   "  ",
			expected: "  this is a long\n  sentence that\n  needs wrapping",
		},

		{
			name:     "single very long word",
			text:     "supercalifragilisticexpialidocious",
			width:    20,
			indent:   "  ",
			expected: "  supercalifragilisticexpialidocious",
		},

		{
			name:     "exact width boundary",
			text:     "twelve chars",
			width:    14, // "  " + "twelve chars" = 14 chars exactly
			indent:   "  ",
			expected: "  twelve chars",
		},

		{
			name:     "width boundary exceeded by one char",
			text:     "thirteen chars",
			width:    14, // "  " + "thirteen" = 10, but "thirteen chars" = 15 > 14
			indent:   "  ",
			expected: "  thirteen\n  chars",
		},

		{
			name:     "no indent",
			text:     "hello world test",
			width:    12,
			indent:   "",
			expected: "hello world\ntest",
		},

		{
			name:     "custom indent",
			text:     "hello world",
			width:    20,
			indent:   "    ",
			expected: "    hello world",
		},

		{
			name:     "multiple spaces between words",
			text:     "hello    world    test",
			width:    20,
			indent:   "  ",
			expected: "  hello world test",
		},

		{
			name:     "leading and trailing spaces",
			text:     "  hello world  ",
			width:    20,
			indent:   "  ",
			expected: "  hello world",
		},

		{
			name:     "newlines in text are treated as spaces",
			text:     "hello\nworld\ntest",
			width:    20,
			indent:   "  ",
			expected: "  hello world test",
		},

		{
			name:     "tabs in text are treated as spaces",
			text:     "hello\tworld\ttest",
			width:    20,
			indent:   "  ",
			expected: "  hello world test",
		},

		{
			name:     "very small width",
			text:     "hello world",
			width:    8, // smaller than "  hello"
			indent:   "  ",
			expected: "  hello\n  world",
		},

		{
			name:     "width smaller than indent",
			text:     "hello",
			width:    4,
			indent:   "      ", // 6 spaces, wider than width
			expected: "      hello",
		},

		{
			name:     "multiple consecutive wraps",
			text:     "a b c d e f g h i j k l m n o p q r s t u v w x y z",
			width:    10,
			indent:   "  ",
			expected: "  a b c d\n  e f g h\n  i j k l\n  m n o p\n  q r s t\n  u v w x\n  y z",
		},

		{
			name:     "real usage example",
			text:     "Print verbose output including debugging information and detailed error messages",
			width:    72,
			indent:   "    ",
			expected: "    Print verbose output including debugging information and detailed\n    error messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Do(tt.text, tt.width, tt.indent)
			if result != tt.expected {
				t.Errorf("wrapText() mismatch:\nexpected:\n%q\ngot:\n%q", tt.expected, result)

				// Also show with visible whitespace for debugging
				expectedLines := strings.Split(tt.expected, "\n")
				resultLines := strings.Split(result, "\n")
				t.Errorf("Expected lines:")
				for i, line := range expectedLines {
					t.Errorf("  [%d]: %q (len=%d)", i, line, len(line))
				}
				t.Errorf("Actual lines:")
				for i, line := range resultLines {
					t.Errorf("  [%d]: %q (len=%d)", i, line, len(line))
				}
			}
		})
	}
}
