// permute_test.go - permutation tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"testing"

	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/google/go-cmp/cmp"
)

func Test_maybePermute(t *testing.T) {

	// The command line we are dealing with is the following:
	//
	//	rbmk-run -v testlist7.txt testlist111.txt --logs logs.jsonl -o output.txt testlist444.txt -- curl -o /dev/null

	programName := ValueProgramName{Name: "rbmk-run", Tok: scanner.ProgramNameToken{Idx: 0}}

	options := []Value{
		ValueOption{Option: &Option{Prefix: "-", Name: "v"}, Tok: scanner.OptionToken{Idx: 1}},
		ValueOption{Option: &Option{Prefix: "--", Name: "logs"}, Tok: scanner.OptionToken{Idx: 4}, Value: "logs.jsonl"},
		ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Tok: scanner.PositionalArgumentToken{Idx: 6}, Value: "output.txt"},
	}

	positionals := []Value{
		ValuePositionalArgument{Value: "testlist7.txt", Tok: scanner.PositionalArgumentToken{Idx: 2}},
		ValuePositionalArgument{Value: "testlist111.txt", Tok: scanner.PositionalArgumentToken{Idx: 3}},
		ValuePositionalArgument{Value: "testlist444.txt", Tok: scanner.PositionalArgumentToken{Idx: 8}},
		ValueOptionsArgumentsSeparator{Separator: "--", Tok: scanner.OptionsArgumentsSeparatorToken{Idx: 9}},
		ValuePositionalArgument{Value: "curl", Tok: scanner.PositionalArgumentToken{Idx: 10}},
		ValuePositionalArgument{Value: "-o", Tok: scanner.PositionalArgumentToken{Idx: 11}},
		ValuePositionalArgument{Value: "/dev/null", Tok: scanner.PositionalArgumentToken{Idx: 12}},
	}

	// Define the test cases
	type testcase struct {
		name    string
		disable bool
		expect  []Value
	}
	cases := []testcase{
		{
			name:    "with permutation",
			disable: false,
			expect: []Value{
				// program name
				ValueProgramName{Name: "rbmk-run", Tok: scanner.ProgramNameToken{Idx: 0}},

				// sorted options
				ValueOption{Option: &Option{Prefix: "-", Name: "v"}, Tok: scanner.OptionToken{Idx: 1}},
				ValueOption{Option: &Option{Prefix: "--", Name: "logs"}, Tok: scanner.OptionToken{Idx: 4}, Value: "logs.jsonl"},
				ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Tok: scanner.PositionalArgumentToken{Idx: 6}, Value: "output.txt"},

				// sorted positional arguments
				ValuePositionalArgument{Value: "testlist7.txt", Tok: scanner.PositionalArgumentToken{Idx: 2}},
				ValuePositionalArgument{Value: "testlist111.txt", Tok: scanner.PositionalArgumentToken{Idx: 3}},
				ValuePositionalArgument{Value: "testlist444.txt", Tok: scanner.PositionalArgumentToken{Idx: 8}},
				ValueOptionsArgumentsSeparator{Separator: "--", Tok: scanner.OptionsArgumentsSeparatorToken{Idx: 9}},
				ValuePositionalArgument{Value: "curl", Tok: scanner.PositionalArgumentToken{Idx: 10}},
				ValuePositionalArgument{Value: "-o", Tok: scanner.PositionalArgumentToken{Idx: 11}},
				ValuePositionalArgument{Value: "/dev/null", Tok: scanner.PositionalArgumentToken{Idx: 12}},
			},
		},

		{
			name:    "without permutation",
			disable: true,
			expect: []Value{
				ValueProgramName{Name: "rbmk-run", Tok: scanner.ProgramNameToken{Idx: 0}},
				ValueOption{Option: &Option{Prefix: "-", Name: "v"}, Tok: scanner.OptionToken{Idx: 1}},
				ValuePositionalArgument{Value: "testlist7.txt", Tok: scanner.PositionalArgumentToken{Idx: 2}},
				ValuePositionalArgument{Value: "testlist111.txt", Tok: scanner.PositionalArgumentToken{Idx: 3}},
				ValueOption{Option: &Option{Prefix: "--", Name: "logs"}, Tok: scanner.OptionToken{Idx: 4}, Value: "logs.jsonl"},
				ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Tok: scanner.PositionalArgumentToken{Idx: 6}, Value: "output.txt"},
				ValuePositionalArgument{Value: "testlist444.txt", Tok: scanner.PositionalArgumentToken{Idx: 8}},
				ValueOptionsArgumentsSeparator{Separator: "--", Tok: scanner.OptionsArgumentsSeparatorToken{Idx: 9}},
				ValuePositionalArgument{Value: "curl", Tok: scanner.PositionalArgumentToken{Idx: 10}},
				ValuePositionalArgument{Value: "-o", Tok: scanner.PositionalArgumentToken{Idx: 11}},
				ValuePositionalArgument{Value: "/dev/null", Tok: scanner.PositionalArgumentToken{Idx: 12}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := permute(&config{parser: &Parser{DisablePermute: tc.disable}}, programName, options, positionals)
			if diff := cmp.Diff(tc.expect, err); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
