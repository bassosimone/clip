// parser_test.go - parser tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestErrTooFewPositionalArguments(t *testing.T) {
	err := ErrTooFewPositionalArguments{Min: 3, Have: 1}
	want := "too few positional arguments: expected at least 3, got 1"
	if got := err.Error(); got != want {
		t.Errorf("unexpected error message: got %q, want %q", got, want)
	}
}

func TestErrTooManyPositionalArguments(t *testing.T) {
	err := ErrTooManyPositionalArguments{Max: 3, Have: 5}
	want := "too many positional arguments: expected at most 3, got 5"
	if got := err.Error(); got != want {
		t.Errorf("unexpected error message: got %q, want %q", got, want)
	}
}

func TestParser_Parse(t *testing.T) {
	// Define the test case structure
	type testcase struct {
		argv        []string // argument vector to parse (including program name)
		skipCase    bool     // whether to skip the test case
		px          *Parser  // parser to use
		expectValue []string // expected parsed and reserialized values
		expectErr   error    // expected error, if any
	}

	cases := []testcase{

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSLOfile.txt"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "O",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{"curl", "-f", "-s", "-S", "-L", "-O", "file.txt", "https://example.com/file.txt"},
			expectErr:   nil,
		},

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSLO", "file.txt"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "O",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{"curl", "-f", "-s", "-S", "-L", "-O", "file.txt", "https://example.com/file.txt"},
			expectErr:   nil,
		},

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSL", "--output=file.txt"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "output",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentRequired,
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{"curl", "-f", "-s", "-S", "-L", "--output", "file.txt", "https://example.com/file.txt"},
			expectErr:   nil,
		},

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSL", "--output", "file.txt"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "output",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentRequired,
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{"curl", "-f", "-s", "-S", "-L", "--output", "file.txt", "https://example.com/file.txt"},
			expectErr:   nil,
		},

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSL", "--output"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:         "output",
						Prefix:       "--",
						Type:         OptionTypeStandaloneArgumentOptional,
						DefaultValue: "FILE",
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{"curl", "-f", "-s", "-S", "-L", "--output=FILE", "https://example.com/file.txt"},
			expectErr:   nil,
		},

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSL", "--output"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "output",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentRequired,
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{},
			expectErr:   errors.New("option requires an argument: --output"),
		},

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSL", "--output=FOO"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "output",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentNone,
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{},
			expectErr:   errors.New("option requires no argument: --output"),
		},

		{
			argv:     []string{"curl", "https://example.com/file.txt", "-fsSL", "--output2"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "f",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "L",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "output",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentNone,
					},
					{
						Name:   "S",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "s",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{},
			expectErr:   errors.New("unknown option: --output2"),
		},
		{
			argv:     []string{"multirepo", "foreach", "-kx", "--", "git", "pull", "--rebase"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "k",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "x",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{"multirepo", "-k", "-x", "foreach", "--", "git", "pull", "--rebase"},
			expectErr:   nil,
		},

		{
			argv:     []string{"multirepo", "foreach", "-kx", "--", "git", "pull", "--rebase", "--help"},
			skipCase: false,
			px: &Parser{
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "help",
						Prefix: "--",
						Type:   OptionTypeEarlyArgumentNone,
					},
				},
			},
			expectValue: []string{"multirepo", "--help"},
			expectErr:   nil,
		},

		{
			argv:     []string{"multirepo-foreach", "-kx", "git", "pull", "--rebase"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            true,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    math.MaxInt,
				Options: []*Option{
					{
						Name:   "k",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
					{
						Name:   "x",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentNone,
					},
				},
			},
			expectValue: []string{"multirepo-foreach", "-k", "-x", "git", "pull", "--rebase"},
			expectErr:   nil,
		},

		{
			argv:     []string{"dig", "@8.8.8.8", "-p53", "IN", "+short", "A", "example.com"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    4,
				Options: []*Option{
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
					{
						Name:   "short",
						Prefix: "+",
						Type:   OptionTypeStandaloneArgumentNone,
					},
				},
			},
			expectValue: []string{"dig", "-p", "53", "+short", "@8.8.8.8", "IN", "A", "example.com"},
			expectErr:   nil,
		},

		{
			argv:     []string{"dig", "@8.8.8.8", "-p53", "IN", "+short", "A", "--", "-h"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    4,
				Options: []*Option{
					{
						Name:   "h",
						Prefix: "-",
						Type:   OptionTypeEarlyArgumentNone,
					},
				},
			},
			expectValue: []string{"dig", "-h"},
			expectErr:   nil,
		},

		{
			argv:     []string{"dig", "@8.8.8.8", "-P53", "IN", "+short", "A", "--"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    4,
				Options: []*Option{
					{
						Name:   "h",
						Prefix: "-",
						Type:   OptionTypeEarlyArgumentNone,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("unknown option: -P"),
		},

		{
			argv:     []string{"dig", "@8.8.8.8", "-p53", "IN", "+short", "A", "-p"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    4,
				Options: []*Option{
					{
						Name:   "h",
						Prefix: "-",
						Type:   OptionTypeEarlyArgumentNone,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("option requires an argument: -p"),
		},

		{
			argv:     []string{},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    4,
				Options: []*Option{
					{
						Name:   "h",
						Prefix: "-",
						Type:   OptionTypeEarlyArgumentNone,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("missing program name"),
		},

		{
			argv:     []string{"dig"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    4,
				Options: []*Option{
					{
						Name:   "h",
						Prefix: "-",
						Type:   OptionTypeEarlyArgumentNone,
					},
					{
						Name:   "port",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("groupable option names should be a single byte, found: &{DefaultValue: Prefix:- Name:port Type:66}"),
		},

		{
			argv:     []string{"dig"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    4,
				Options: []*Option{
					{
						Name:   "h",
						Prefix: "-",
						Type:   OptionTypeEarlyArgumentNone,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("too few positional arguments: expected at least 1, got 0"),
		},

		{
			argv:     []string{"dig", "IN", "A", "example.com"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    1,
				Options: []*Option{
					{
						Name:   "h",
						Prefix: "-",
						Type:   OptionTypeEarlyArgumentNone,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("too many positional arguments: expected at most 1, got 3"),
		},

		{
			argv:     []string{"dig", "IN", "A", "example.com"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    1,
				Options: []*Option{
					{
						Name:   "p",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentRequired,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("multiple options with \"p\" name"),
		},

		{
			argv:     []string{"dig", "IN", "A", "example.com"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    1,
				Options: []*Option{
					{
						Name:   "",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentRequired,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("option name cannot be empty: &{DefaultValue: Prefix:-- Name: Type:34}"),
		},

		{
			argv:     []string{"dig", "IN", "A", "example.com"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    1,
				Options: []*Option{
					{
						Name:   "short",
						Prefix: "",
						Type:   OptionTypeStandaloneArgumentRequired,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("option prefix cannot be empty: &{DefaultValue: Prefix: Name:short Type:34}"),
		},

		{
			argv:     []string{"dig", "IN", "A", "example.com"},
			skipCase: false,
			px: &Parser{
				DisablePermute:            false,
				OptionsArgumentsSeparator: "--",
				MinPositionalArguments:    1,
				MaxPositionalArguments:    1,
				Options: []*Option{
					{
						Name:   "short",
						Prefix: "-",
						Type:   OptionTypeStandaloneArgumentRequired,
					},
					{
						Name:   "p",
						Prefix: "-",
						Type:   OptionTypeGroupableArgumentRequired,
					},
				},
			},
			expectValue: nil,
			expectErr:   errors.New("prefix \"-\" is used for both standalone and groupable options"),
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%+v", tc.argv), func(t *testing.T) {
			// Possibly skip the test case if asked to do so
			if tc.skipCase {
				t.Skip("skipping test case")
			}

			// Parse the arguments using the parser
			values, err := tc.px.Parse(tc.argv)

			// Check for expected error
			switch {
			case tc.expectErr == nil && err != nil:
				t.Fatalf("unexpected error: %v", err)

			case tc.expectErr != nil && err == nil:
				t.Fatalf("expected error %v, got nil", tc.expectErr)

			case err != nil && err.Error() != tc.expectErr.Error():
				t.Fatalf("expected error %+v, got %+v", tc.expectErr, err)

			case err != nil:
				return
			}

			// Reserialize the values to another argument vector
			got := []string{}
			for _, entry := range values {
				got = append(got, entry.Strings()...)
			}

			// Compare to the expectations
			if diff := cmp.Diff(tc.expectValue, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestParserEmptyDefaultsToGNUStyleOptions(t *testing.T) {
	// Create a new empty parser with no options
	px := &Parser{}

	// Parse an empty argument vector
	values, err := px.Parse([]string{"program", "--option", "value"})

	// Check for errors
	var unknownOption ErrUnknownOption
	if !errors.As(err, &unknownOption) {
		t.Fatalf("expected ErrUnknownOption, got %v", err)
	}

	if values != nil {
		t.Fatalf("expected no values, got %v", values)
	}
}
