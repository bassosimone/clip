// early_test.go - early options tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"testing"

	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/google/go-cmp/cmp"
)

func Test_parseEarly(t *testing.T) {

	// Define the configuration to use
	px := &Parser{
		Options: []*Option{
			{
				DefaultValue: "",
				Prefix:       "-",
				Name:         "h",
				Type:         OptionTypeEarlyArgumentNone,
			},
			{
				DefaultValue: "",
				Prefix:       "--",
				Name:         "help",
				Type:         OptionTypeEarlyArgumentNone,
			},
			{
				DefaultValue: "",
				Prefix:       "+",
				Name:         "short",
				Type:         OptionTypeStandaloneArgumentNone,
			},
		},
	}

	// Define the test cases
	type testcase struct {
		name        string   // name of the test case
		skipCase    bool     // whether to skip the test case
		argv        []string // argument vector -- includes the program name
		expectValue Value    // expected parsed value
	}
	cases := []testcase{
		{
			name:     "successful recognition of --help with an invalid command line",
			skipCase: false,
			argv: []string{
				"program",
				"-x",
				"file1.txt",
				"--verbose",
				"file2.txt",
				"+short",
				"file4.txt",
				"--file",
				"/dev/null",
				"--",
				"-z",
				"file5.txt",
				"--http",
				"--http=2.0",
				"file5.txt",
				"--help",
			},
			expectValue: ValueOption{
				Option: px.Options[1],
				Tok:    scanner.OptionToken{Idx: 15, Prefix: "--", Name: "help"},
				Value:  "",
			},
		},

		{
			name:     "successful recognition of -h with an invalid command line",
			skipCase: false,
			argv: []string{
				"program",
				"-x",
				"file1.txt",
				"--verbose",
				"file2.txt",
				"+short",
				"file4.txt",
				"--file",
				"/dev/null",
				"--",
				"-z",
				"file5.txt",
				"--http",
				"--http=2.0",
				"file5.txt",
				"-h",
			},
			expectValue: ValueOption{
				Option: px.Options[0],
				Tok:    scanner.OptionToken{Idx: 15, Prefix: "-", Name: "h"},
				Value:  "",
			},
		},

		{
			name:     "no early option in the command line",
			skipCase: false,
			argv: []string{
				"program",
				"-x",
				"file1.txt",
				"--verbose",
				"file2.txt",
				"+short",
				"file4.txt",
				"--file",
				"/dev/null",
				"--",
				"-z",
				"file5.txt",
				"--http",
				"--http=2.0",
				"file5.txt",
			},
			expectValue: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Possibly skip the test case if requested
			if tc.skipCase {
				t.Skip("skipping test case:", tc.name)
			}

			// Run the function we are testing
			value, found := searchEarly(px, tc.argv)

			expectFound := tc.expectValue != nil
			if expectFound && !found {
				t.Fatal("expected to find an early option, but did not")
			}

			// Check for success
			if diff := cmp.Diff(tc.expectValue, value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
