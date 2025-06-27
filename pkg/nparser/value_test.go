// value_test.go - parsed value tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"testing"

	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/google/go-cmp/cmp"
)

func TestValue(t *testing.T) {
	// Just a random token for testing the Token method
	testtoken := scanner.PositionalArgumentToken{
		Idx:   155,
		Value: "antani",
	}

	// Test case definition
	type testcase struct {
		name    string
		input   Value
		strings []string
		panics  bool
	}

	cases := []testcase{
		{
			name: "ValueProgramName",
			input: ValueProgramName{
				Name: "curl",
				Tok:  testtoken,
			},
			strings: []string{"curl"},
			panics:  false,
		},

		{
			name: "OptionTypeEarlyArgumentNone",
			input: ValueOption{
				Tok: testtoken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "help",
					Type:         OptionTypeEarlyArgumentNone,
				},
				Value: "xx",
			},
			strings: []string{"-help"},
			panics:  false,
		},

		{
			name: "OptionTypeGroupableArgumentNone",
			input: ValueOption{
				Tok: testtoken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "z",
					Type:         OptionTypeGroupableArgumentNone,
				},
				Value: "xx",
			},
			strings: []string{"-z"},
			panics:  false,
		},

		{
			name: "OptionTypeStandaloneArgumentNone",
			input: ValueOption{
				Tok: testtoken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "--",
					Name:         "verbose",
					Type:         OptionTypeStandaloneArgumentNone,
				},
				Value: "xx",
			},
			strings: []string{"--verbose"},
			panics:  false,
		},

		{
			name: "OptionTypeStandaloneArgumentOptional",
			input: ValueOption{
				Tok: testtoken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "--",
					Name:         "verbose",
					Type:         OptionTypeStandaloneArgumentOptional,
				},
				Value: "false",
			},
			strings: []string{"--verbose=false"},
			panics:  false,
		},

		{
			name: "OptionTypeStandaloneArgumentRequired",
			input: ValueOption{
				Tok: testtoken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "--",
					Name:         "file",
					Type:         OptionTypeStandaloneArgumentRequired,
				},
				Value: "/dev/null",
			},
			strings: []string{"--file", "/dev/null"},
			panics:  false,
		},

		{
			name: "OptionTypeGroupableArgumentRequired",
			input: ValueOption{
				Tok: testtoken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "o",
					Type:         OptionTypeGroupableArgumentRequired,
				},
				Value: "/dev/null",
			},
			strings: []string{"-o", "/dev/null"},
			panics:  false,
		},

		{
			name: "OptionType_invalid",
			input: ValueOption{
				Tok: testtoken,
				Option: &Option{
					DefaultValue: "antani",
					Prefix:       "-",
					Name:         "o",
					Type:         0, // invalid
				},
				Value: "/dev/null",
			},
			strings: []string{},
			panics:  true,
		},

		{
			name: "ValuePositionalArgument",
			input: ValuePositionalArgument{
				Tok:   testtoken,
				Value: "/dev/null",
			},
			strings: []string{"/dev/null"},
			panics:  false,
		},

		{
			name: "ValueOptionsArgumentsSeparator",
			input: ValueOptionsArgumentsSeparator{
				Separator: "--",
				Tok:       testtoken,
			},
			strings: []string{"--"},
			panics:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// make sure the code only panics when expected
			defer func() {
				r := recover()
				if got := r != nil; got != tc.panics {
					t.Fatal("expected", tc.panics, "got", got)
				}
			}()

			gotStrings := tc.input.Strings()
			if diff := cmp.Diff(tc.strings, gotStrings); diff != "" {
				t.Fatal(diff)
			}

			gotToken := tc.input.Token()
			if diff := cmp.Diff(testtoken, gotToken); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func Test_sortValues(t *testing.T) {
	input := []Value{
		ValuePositionalArgument{Tok: scanner.PositionalArgumentToken{Idx: 2}, Value: "b"},
		ValuePositionalArgument{Tok: scanner.PositionalArgumentToken{Idx: 1}, Value: "a"},
		ValuePositionalArgument{Tok: scanner.PositionalArgumentToken{Idx: 3}, Value: "c"},
	}

	expected := []Value{
		ValuePositionalArgument{Tok: scanner.PositionalArgumentToken{Idx: 1}, Value: "a"},
		ValuePositionalArgument{Tok: scanner.PositionalArgumentToken{Idx: 2}, Value: "b"},
		ValuePositionalArgument{Tok: scanner.PositionalArgumentToken{Idx: 3}, Value: "c"},
	}

	sortValues(input)

	if diff := cmp.Diff(expected, input); diff != "" {
		t.Fatal(diff)
	}
}
