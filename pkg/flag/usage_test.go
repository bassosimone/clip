// usage_test.go - Tests for usage formatting functionality
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"strings"
	"testing"

	"github.com/bassosimone/clip/pkg/parser"
)

func TestFlagSet_UsageOptions(t *testing.T) {
	tests := []struct {
		name          string
		shortPrefixes []string
		longPrefixes  []string
		options       []*Option
		expected      string
	}{
		{
			name:          "short option only with bool value",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					ShortName: 'h',
					Usage:     "show help",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
			},
			expected: "  -h\n    show help\n\n",
		},

		{
			name:          "long option only with bool value",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					LongName: "help",
					Usage:    "show help",
					Value:    &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
			},
			expected: "  --help\n    show help\n\n",
		},

		{
			name:          "short and long option with bool value",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					ShortName: 'h',
					LongName:  "help",
					Usage:     "show help",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
			},
			expected: "  -h, --help\n    show help\n\n",
		},

		{
			name:          "short option only with string value",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					ShortName: 'f',
					Usage:     "input file",
					Value:     &mockValue{optionType: parser.OptionTypeString, value: ""},
				},
			},
			expected: "  -f VALUE\n    input file\n\n",
		},

		{
			name:          "long option only with string value",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					LongName: "file",
					Usage:    "input file",
					Value:    &mockValue{optionType: parser.OptionTypeString, value: ""},
				},
			},
			expected: "  --file VALUE\n    input file\n\n",
		},

		{
			name:          "short and long option with string value",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					ShortName: 'f',
					LongName:  "file",
					Usage:     "input file",
					Value:     &mockValue{optionType: parser.OptionTypeString, value: ""},
				},
			},
			expected: "  -f, --file VALUE\n    input file\n\n",
		},

		{
			name:          "multiple options mixed types",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					ShortName: 'h',
					LongName:  "help",
					Usage:     "show help",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
				{
					ShortName: 'v',
					Usage:     "verbose output",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
				{
					LongName: "config",
					Usage:    "configuration file path",
					Value:    &mockValue{optionType: parser.OptionTypeString, value: ""},
				},
			},
			expected: "  --config VALUE\n    configuration file path\n\n  -h, --help\n    show help\n\n  -v\n    verbose output\n\n",
		},

		{
			name:          "custom prefixes",
			shortPrefixes: []string{"/"},
			longPrefixes:  []string{"+"},
			options: []*Option{
				{
					ShortName: 'h',
					LongName:  "help",
					Usage:     "show help",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
			},
			expected: "  /h, +help\n    show help\n\n",
		},

		{
			name:          "empty prefixes fall back gracefully",
			shortPrefixes: []string{},
			longPrefixes:  []string{},
			options: []*Option{
				{
					ShortName: 'h',
					LongName:  "help",
					Usage:     "show help",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
			},
			expected: "  h, help\n    show help\n\n",
		},

		{
			name:          "option with no short or long name should be skipped",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					Usage: "this should be skipped",
					Value: &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
				{
					ShortName: 'h',
					Usage:     "show help",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
			},
			expected: "  -h\n    show help\n\n",
		},

		{
			name:          "long text wrapping",
			shortPrefixes: []string{"-"},
			longPrefixes:  []string{"--"},
			options: []*Option{
				{
					ShortName: 'h',
					Usage:     "this is a very long usage description that should wrap to multiple lines when it exceeds the maximum width limit",
					Value:     &mockValue{optionType: parser.OptionTypeBool, value: "false"},
				},
			},
			expected: "  -h\n    this is a very long usage description that should wrap to multiple\n    lines when it exceeds the maximum width limit\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a FlagSet with custom parser configuration
			fx := NewFlagSet("test", ContinueOnError)
			fx.parser.ShortOptionPrefixes = tt.shortPrefixes
			fx.parser.LongOptionPrefixes = tt.longPrefixes

			// Add all options to the flag set
			for _, opt := range tt.options {
				fx.AddOption(opt)
			}

			// Get the usage options output
			result := fx.UsageOptions()

			// Compare the result
			if result != tt.expected {
				t.Fatalf("UsageOptions() mismatch:\nexpected:\n%q\ngot:\n%q", tt.expected, result)
			}
		})
	}
}

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
			result := wrapText(tt.text, tt.width, tt.indent)
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
