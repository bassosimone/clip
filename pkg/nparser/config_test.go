// config_test.go - parser config tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/google/go-cmp/cmp"
)

func TestErrUnknownOption(t *testing.T) {
	err := ErrUnknownOption{
		Name:   "verbose",
		Prefix: "--",
		Token: scanner.OptionToken{
			Idx:    4,
			Prefix: "--",
			Name:   "verbose",
		},
	}

	expect := "unknown option: --verbose"

	if diff := cmp.Diff(expect, err.Error()); diff != "" {
		t.Fatal(diff)
	}
}

func TestErrAmbiguousPrefix(t *testing.T) {
	err := ErrAmbiguousPrefix{
		Prefix: "-",
	}

	expect := `prefix "-" is used for both standalone and groupable options`
	if diff := cmp.Diff(expect, err.Error()); diff != "" {
		t.Fatal(diff)
	}
}

func TestErrMultipleOptionsWithSameName(t *testing.T) {
	expect := `multiple options with "foo" name`

	opt1 := &Option{Name: "foo"}
	opt2 := &Option{Name: "foo"}
	err := ErrMultipleOptionsWithSameName{
		Name:    "foo",
		Options: []*Option{opt1, opt2},
	}

	if got := err.Error(); got != expect {
		t.Fatalf("expected prefix %q, got %q", expect, got)
	}
}

func TestErrTooLongGroupableOptionName(t *testing.T) {
	opt := &Option{Name: "longname"}
	err := ErrTooLongGroupableOptionName{Option: opt}

	expect := "groupable option names should be a single byte, found:"
	if got := err.Error(); len(got) < len(expect) || got[:len(expect)] != expect {
		t.Fatalf("expected prefix %q, got %q", expect, got)
	}
}

func TestErrEmptyOptionName(t *testing.T) {
	opt := &Option{Name: ""}
	err := ErrEmptyOptionName{Option: opt}

	expect := "option name cannot be empty:"
	if got := err.Error(); len(got) < len(expect) || got[:len(expect)] != expect {
		t.Fatalf("expected prefix %q, got %q", expect, got)
	}
}

func TestErrEmptyOptionPrefix(t *testing.T) {
	opt := &Option{Prefix: ""}
	err := ErrEmptyOptionPrefix{Option: opt}

	expect := "option prefix cannot be empty:"
	if got := err.Error(); len(got) < len(expect) || got[:len(expect)] != expect {
		t.Fatalf("expected prefix %q, got %q", expect, got)
	}
}

func Test_config_disablePermute(t *testing.T) {
	cases := []bool{true, false}
	for _, tc := range cases {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			cfg := config{parser: &Parser{DisablePermute: tc}}
			if got := cfg.disablePermute(); got != cfg.parser.DisablePermute {
				t.Fatal("expected", cfg.parser.DisablePermute, "got", got)
			}
		})
	}
}

func Test_config_findOption(t *testing.T) {
	// Create the option we would like to return to the caller
	option := Option{
		DefaultValue: "",
		Prefix:       "--",
		Name:         "verbose",
		Type:         OptionTypeStandaloneArgumentNone,
	}

	// Create a parser with a single option inside
	cfg := config{
		options: map[string]*Option{
			"verbose": &option,
		},
	}

	// Define the test cases
	type testcase struct {
		caseName         string
		tok              scanner.OptionToken
		optName          string
		kind             OptionType
		expectOp         *Option
		expectErrUnknown bool
	}
	cases := []testcase{
		{
			caseName: "successful find",
			tok: scanner.OptionToken{
				Idx:    4,
				Prefix: "--",
				Name:   "verbose",
			},
			optName:          "verbose",
			kind:             optionKindStandalone,
			expectOp:         &option,
			expectErrUnknown: false,
		},

		{
			caseName: "no such option",
			tok: scanner.OptionToken{
				Idx:    4,
				Prefix: "--",
				Name:   "verbose",
			},
			optName:          "file",
			kind:             optionKindStandalone,
			expectOp:         nil,
			expectErrUnknown: true,
		},

		{
			caseName: "the prefix does not match",
			tok: scanner.OptionToken{
				Idx:    4,
				Prefix: "/",
				Name:   "verbose",
			},
			optName:          "verbose",
			kind:             optionKindStandalone,
			expectOp:         nil,
			expectErrUnknown: true,
		},

		{
			caseName: "the option kind does not match",
			tok: scanner.OptionToken{
				Idx:    4,
				Prefix: "--",
				Name:   "verbose",
			},
			optName:          "verbose",
			kind:             optionKindEarly,
			expectOp:         nil,
			expectErrUnknown: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			gotOption, err := cfg.findOption(tc.tok, tc.optName, tc.kind)

			switch {
			case tc.expectErrUnknown:
				var errval ErrUnknownOption
				if !errors.As(err, &errval) {
					t.Fatalf("cannot convert error to ErrUnknownOption: %T", err)
				}
				if errval.Name != tc.optName {
					t.Fatal("expected", tc.optName, "got", errval.Name)
				}
				if errval.Prefix != tc.tok.Prefix {
					t.Fatal("expected", tc.tok.Prefix, "got", errval.Prefix)
				}
				if diff := cmp.Diff(tc.tok, errval.Token); diff != "" {
					t.Fatal(diff)
				}

			case err != nil:
				t.Fatal(err)
			}

			if gotOption != tc.expectOp {
				t.Fatal("expected", tc.expectOp, "got", gotOption)
			}
		})
	}
}

func Test_newConfig(t *testing.T) {
	// Define the structure of the test cases
	type testcase struct {
		caseName       string                // Name of the test case
		options        []*Option             // Options to be used in the parser
		expectErr      error                 // Expected error, if any
		expectPrefixes map[string]OptionType // Expected prefixes and their types
		expectOptions  map[string]*Option    // Expected options by name
	}

	// Define the test cases
	cases := []testcase{
		{
			caseName: "groupable option with multi-byte name",
			options: []*Option{
				{
					Name:   "longname",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
			expectErr: ErrTooLongGroupableOptionName{
				Option: &Option{
					Name:   "longname",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "empty option name",
			options: []*Option{
				{
					Name:   "",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectErr: ErrEmptyOptionName{
				Option: &Option{
					Name:   "",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "empty option prefix",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectErr: ErrEmptyOptionPrefix{
				Option: &Option{
					Name:   "verbose",
					Prefix: "",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "multiple options with same name",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "--",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				{
					Name:   "verbose",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
			},
			expectErr: ErrMultipleOptionsWithSameName{
				Name: "verbose",
				Options: []*Option{
					{
						Name:   "verbose",
						Prefix: "--",
						Type:   OptionTypeStandaloneArgumentNone,
					},
					{
						Name:   "verbose",
						Prefix: "-",
						Type:   OptionTypeStandaloneArgumentNone,
					},
				},
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "ambiguous parsing prefixes",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "-",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				{
					Name:   "v",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
			expectErr: ErrAmbiguousPrefix{
				Prefix: "-",
			},
			expectPrefixes: map[string]OptionType{},
			expectOptions:  map[string]*Option{},
		},

		{
			caseName: "valid configuration",
			options: []*Option{
				{
					Name:   "verbose",
					Prefix: "--",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				{
					Name:   "v",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
				{
					Name:   "help",
					Prefix: "--",
					Type:   OptionTypeEarlyArgumentNone,
				},
				{
					Name:   "h",
					Prefix: "-",
					Type:   OptionTypeEarlyArgumentNone,
				},
			},
			expectErr: nil,
			expectPrefixes: map[string]OptionType{
				"--": optionKindStandalone,
				"-":  optionKindGroupable,
			},
			expectOptions: map[string]*Option{
				"h": {
					Name:   "h",
					Prefix: "-",
					Type:   OptionTypeEarlyArgumentNone,
				},
				"help": {
					Name:   "help",
					Prefix: "--",
					Type:   OptionTypeEarlyArgumentNone,
				},
				"verbose": {
					Name:   "verbose",
					Prefix: "--",
					Type:   OptionTypeStandaloneArgumentNone,
				},
				"v": {
					Name:   "v",
					Prefix: "-",
					Type:   OptionTypeGroupableArgumentNone,
				},
			},
		},
	}

	// Run through each test case
	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			// Create a parser with the provided options
			parser := &Parser{Options: tc.options}

			// Attempt to create a new config
			cfg, err := newConfig(parser)

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

			// Check the prefixes and options in the config
			if diff := cmp.Diff(tc.expectPrefixes, cfg.prefixes); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.expectOptions, cfg.options); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
