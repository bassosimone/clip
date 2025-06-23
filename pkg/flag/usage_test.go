// usage_test.go - Tests for usage formatting functionality
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
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

func TestFlagSet_firstSeparator(t *testing.T) {
	fx := NewFlagSet("test", ContinueOnError)
	fx.parser.Separators = []string{"--"}

	if fx.firstSeparator() != " [--] " {
		t.Errorf("Expected separator to be ' [--] ', got '%s'", fx.firstSeparator())
	}

	fx.parser.Separators = []string{}
	if fx.firstSeparator() != "" {
		t.Errorf("Expected no separator, got '%s'", fx.firstSeparator())
	}
}
