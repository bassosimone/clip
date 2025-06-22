// long_test.go - getopt_long tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"testing"

	"github.com/bassosimone/clip/pkg/parser"
	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/google/go-cmp/cmp"
)

func TestLong(t *testing.T) {
	type config struct {
		name      string
		argv      []string
		optstring string
		options   []Option
		env       map[string]string // environment variables to set
		want      []parser.CommandLineItem
	}

	tests := []config{
		{
			name: "with POSIXLY_CORRECT stops at first non-option",
			argv: []string{
				"program",
				"--file=input.txt", // this remains first because it's an option
				"subcommand",       // this stops parsing because of POSIXLY_CORRECT
				"-v",               // this remains an argument
				"--verbose",        // this remains an argument
			},
			optstring: "v",
			options: []Option{
				{Name: "file", HasArg: true},
				{Name: "verbose", HasArg: false},
			},
			env: map[string]string{
				"POSIXLY_CORRECT": "1",
			},
			want: []parser.CommandLineItem{
				parser.ProgramNameItem{
					Name:  "program",
					Token: scanner.ProgramNameToken{Index: 0, Name: "program"},
				},
				parser.OptionItem{
					Name:    "file",
					Token:   scanner.OptionToken{Index: 1, Name: "file=input.txt", Prefix: "--"},
					Value:   "input.txt",
					IsShort: false,
					Type:    parser.OptionTypeString,
					Prefix:  "--",
				},
				parser.ArgumentItem{
					Token: scanner.ArgumentToken{Index: 2, Value: "subcommand"},
					Value: "subcommand",
				},
				parser.ArgumentItem{
					Token: scanner.OptionToken{Index: 3, Prefix: "-", Name: "v"},
					Value: "-v",
				},
				parser.ArgumentItem{
					Token: scanner.OptionToken{Index: 4, Prefix: "--", Name: "verbose"},
					Value: "--verbose",
				},
			},
		},

		{
			name: "with dash prefix stops at first non-option",
			argv: []string{
				"program",
				"--file=input.txt", // this remains first because it's an option
				"subcommand",       // this stops parsing because of -v in optstring
				"-v",               // this remains an argument
				"--verbose",        // this remains an argument
			},
			optstring: "-v", // leading dash disables permutation
			options: []Option{
				{Name: "file", HasArg: true},
				{Name: "verbose", HasArg: false},
			},
			env: map[string]string{}, // no environment variables
			want: []parser.CommandLineItem{
				parser.ProgramNameItem{
					Name:  "program",
					Token: scanner.ProgramNameToken{Index: 0, Name: "program"},
				},
				parser.OptionItem{
					Name:    "file",
					Token:   scanner.OptionToken{Index: 1, Name: "file=input.txt", Prefix: "--"},
					Value:   "input.txt",
					IsShort: false,
					Type:    parser.OptionTypeString,
					Prefix:  "--",
				},
				parser.ArgumentItem{
					Token: scanner.ArgumentToken{Index: 2, Value: "subcommand"},
					Value: "subcommand",
				},
				parser.ArgumentItem{
					Token: scanner.OptionToken{Index: 3, Prefix: "-", Name: "v"},
					Value: "-v",
				},
				parser.ArgumentItem{
					Token: scanner.OptionToken{Index: 4, Prefix: "--", Name: "verbose"},
					Value: "--verbose",
				},
			},
		},

		{
			name: "default behavior reorders options before arguments",
			argv: []string{
				"program",
				"subcommand",       // this gets moved after options in default mode
				"-v",               // this gets reordered before subcommand
				"--file=input.txt", // this gets reordered before subcommand
				"--verbose",        // this gets reordered before subcommand
			},
			optstring: "v",
			options: []Option{
				{Name: "file", HasArg: true},
				{Name: "verbose", HasArg: false},
			},
			env: map[string]string{}, // no environment variables
			want: []parser.CommandLineItem{
				parser.ProgramNameItem{
					Name:  "program",
					Token: scanner.ProgramNameToken{Index: 0, Name: "program"},
				},
				parser.OptionItem{
					Name:    "v",
					Token:   scanner.OptionToken{Index: 2, Name: "v", Prefix: "-"},
					Value:   "true",
					IsShort: true,
					Type:    parser.OptionTypeBool,
					Prefix:  "-",
				},
				parser.OptionItem{
					Name:    "file",
					Token:   scanner.OptionToken{Index: 3, Name: "file=input.txt", Prefix: "--"},
					Value:   "input.txt",
					IsShort: false,
					Type:    parser.OptionTypeString,
					Prefix:  "--",
				},
				parser.OptionItem{
					Name:    "verbose",
					Token:   scanner.OptionToken{Index: 4, Name: "verbose", Prefix: "--"},
					Value:   "true",
					IsShort: false,
					Type:    parser.OptionTypeBool,
					Prefix:  "--",
				},
				parser.ArgumentItem{
					Token: scanner.ArgumentToken{Index: 1, Value: "subcommand"},
					Value: "subcommand",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Save the original lookupEnv and restore it after the test
			originalLookupEnv := lookupEnv
			defer func() { lookupEnv = originalLookupEnv }()

			// Mock lookupEnv to return our test environment
			lookupEnv = func(key string) (string, bool) {
				v, ok := tc.env[key]
				return v, ok
			}

			// Call Long and check the results
			got, err := Long(tc.argv, tc.optstring, tc.options)
			if err != nil {
				t.Fatalf("Long() error = %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Long() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
