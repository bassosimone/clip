// long_test.go - getopt_long tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"testing"

	"github.com/bassosimone/clip/pkg/nparser"
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
		want      []nparser.Value
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
			want: []nparser.Value{
				nparser.ValueProgramName{
					Name: "program",
					Tok:  scanner.ProgramNameToken{Idx: 0, Name: "program"},
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						Name:   "file",
						Prefix: "--",
						Type:   nparser.OptionTypeStandaloneArgumentRequired,
					},
					Tok:   scanner.OptionToken{Idx: 1, Name: "file=input.txt", Prefix: "--"},
					Value: "input.txt",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 2, Value: "subcommand"},
					Value: "subcommand",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 3, Prefix: "-", Name: "v"},
					Value: "-v",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 4, Prefix: "--", Name: "verbose"},
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
			want: []nparser.Value{
				nparser.ValueProgramName{
					Name: "program",
					Tok:  scanner.ProgramNameToken{Idx: 0, Name: "program"},
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						Name:   "file",
						Prefix: "--",
						Type:   nparser.OptionTypeStandaloneArgumentRequired,
					},
					Tok:   scanner.OptionToken{Idx: 1, Name: "file=input.txt", Prefix: "--"},
					Value: "input.txt",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 2, Value: "subcommand"},
					Value: "subcommand",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 3, Prefix: "-", Name: "v"},
					Value: "-v",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 4, Prefix: "--", Name: "verbose"},
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
			want: []nparser.Value{
				nparser.ValueProgramName{
					Name: "program",
					Tok:  scanner.ProgramNameToken{Idx: 0, Name: "program"},
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						Name:   "v",
						Prefix: "-",
						Type:   nparser.OptionTypeGroupableArgumentNone,
					},
					Tok:   scanner.OptionToken{Idx: 2, Name: "v", Prefix: "-"},
					Value: "",
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						Name:   "file",
						Prefix: "--",
						Type:   nparser.OptionTypeStandaloneArgumentRequired,
					},
					Tok:   scanner.OptionToken{Idx: 3, Name: "file=input.txt", Prefix: "--"},
					Value: "input.txt",
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						Prefix: "--",
						Name:   "verbose",
						Type:   nparser.OptionTypeStandaloneArgumentNone,
					},
					Tok:   scanner.OptionToken{Idx: 4, Name: "verbose", Prefix: "--"},
					Value: "",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 1, Value: "subcommand"},
					Value: "subcommand",
				},
			},
		},

		{
			name: "options with optional arguments are possible",
			argv: []string{
				"program",
				"subcommand",       // this gets moved after options in default mode
				"-v",               // this gets reordered before subcommand
				"--file=input.txt", // this gets reordered before subcommand
				"--verbose=true",   // this gets reordered before subcommand
			},
			optstring: "v",
			options: []Option{
				{Name: "file", HasArg: true},
				{Name: "verbose", HasArg: true, IsArgOptional: true, DefaultValue: "false"},
			},
			env: map[string]string{}, // no environment variables
			want: []nparser.Value{
				nparser.ValueProgramName{
					Name: "program",
					Tok:  scanner.ProgramNameToken{Idx: 0, Name: "program"},
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						Name:   "v",
						Prefix: "-",
						Type:   nparser.OptionTypeGroupableArgumentNone,
					},
					Tok:   scanner.OptionToken{Idx: 2, Name: "v", Prefix: "-"},
					Value: "",
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						Name:   "file",
						Prefix: "--",
						Type:   nparser.OptionTypeStandaloneArgumentRequired,
					},
					Tok:   scanner.OptionToken{Idx: 3, Name: "file=input.txt", Prefix: "--"},
					Value: "input.txt",
				},
				nparser.ValueOption{
					Option: &nparser.Option{
						DefaultValue: "false",
						Prefix:       "--",
						Name:         "verbose",
						Type:         nparser.OptionTypeStandaloneArgumentOptional,
					},
					Tok:   scanner.OptionToken{Idx: 4, Name: "verbose=true", Prefix: "--"},
					Value: "true",
				},
				nparser.ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 1, Value: "subcommand"},
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
