// parse.go - parsing algorithm tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"errors"
	"testing"

	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/google/go-cmp/cmp"
)

func TestErrOptionRequiresNoArgument(t *testing.T) {
	err := ErrOptionRequiresNoArgument{
		Option: &Option{
			DefaultValue: "",
			Prefix:       "--",
			Name:         "verbose",
			Type:         OptionTypeStandaloneArgumentNone,
		},
		Token: scanner.OptionToken{
			Idx:    4,
			Prefix: "--",
			Name:   "verbose",
		},
	}

	expect := "option requires no argument: --verbose"

	if diff := cmp.Diff(expect, err.Error()); diff != "" {
		t.Fatal(diff)
	}
}

func TestErrOptionRequiresArgument(t *testing.T) {
	err := ErrOptionRequiresArgument{
		Option: &Option{
			DefaultValue: "",
			Prefix:       "--",
			Name:         "file",
			Type:         OptionTypeStandaloneArgumentRequired,
		},
		Token: scanner.OptionToken{
			Idx:    4,
			Prefix: "--",
			Name:   "file",
		},
	}

	expect := "option requires an argument: --file"

	if diff := cmp.Diff(expect, err.Error()); diff != "" {
		t.Fatal(diff)
	}
}

func Test_parse(t *testing.T) {
	// Define the configuration to use
	cfg := &config{
		parser: &Parser{
			DisablePermute:            false, // we may mutate this inside the test cases
			OptionsArgumentsSeparator: "--",
		},
		prefixes: map[string]OptionType{
			"--": optionKindStandalone,
			"-":  optionKindGroupable,
			"/":  optionKindEarly, // this causes a panic but newConfig should actually prevent this
		},
		options: map[string]*Option{
			"__panic": {
				DefaultValue: "",
				Prefix:       "--",
				Name:         "__panic",
				Type:         optionKindStandalone, // this option type is invalid and triggers a panic
			},
			"file": {
				DefaultValue: "",
				Prefix:       "--",
				Name:         "file",
				Type:         OptionTypeStandaloneArgumentRequired,
			},
			"http": {
				DefaultValue: "1.1",
				Prefix:       "--",
				Name:         "http",
				Type:         OptionTypeStandaloneArgumentOptional,
			},
			"verbose": {
				DefaultValue: "",
				Prefix:       "--",
				Name:         "verbose",
				Type:         OptionTypeStandaloneArgumentNone,
			},
			"_": {
				DefaultValue: "",
				Prefix:       "-",
				Name:         "_",
				Type:         optionKindGroupable, // this option type is invalid and triggers a panic
			},
			"k": {
				DefaultValue: "",
				Prefix:       "-",
				Name:         "k",
				Type:         OptionTypeGroupableArgumentRequired,
			},
			"z": {
				DefaultValue: "",
				Prefix:       "-",
				Name:         "z",
				Type:         OptionTypeGroupableArgumentNone,
			},
			"x": {
				DefaultValue: "",
				Prefix:       "-",
				Name:         "x",
				Type:         OptionTypeGroupableArgumentRequired,
			},
		},
	}

	// Define the test cases
	type testcase struct {
		name              string                // name of the test case
		skipCase          bool                  // whether to skip the test case
		disablePermute    bool                  // whether to disable permutation
		cfg               *config               // configuration to use
		input             *deque[scanner.Token] // input tokens to parse
		expectOptions     *deque[Value]         // expected parsed options
		expectPositionals *deque[Value]         // expected parsed positionals
		expectErrUnknown  bool                  // whether to expect an unknown option error
		expectErrArg      bool                  // whether to expect an error for an option that requires an argument
		expectErrNoArg    bool                  // whether to expect an error for an option that requires no argument
		expectPanic       bool                  // whether to expect a panic
	}
	cases := []testcase{
		{
			name:           "successful parsing with permutation without separator",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
				scanner.PositionalArgumentToken{Idx: 2, Value: "file1.txt"},
				scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
				scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
				scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
				scanner.PositionalArgumentToken{Idx: 6, Value: "file4.txt"},
				scanner.OptionToken{Idx: 7, Prefix: "--", Name: "file"},
				scanner.PositionalArgumentToken{Idx: 8, Value: "/dev/null"},
				scanner.OptionToken{Idx: 9, Prefix: "-", Name: "z"},
				scanner.PositionalArgumentToken{Idx: 10, Value: "file5.txt"},
				scanner.OptionToken{Idx: 11, Prefix: "--", Name: "http"},
				scanner.OptionToken{Idx: 12, Prefix: "--", Name: "http=2.0"},
				scanner.PositionalArgumentToken{Idx: 13, Value: "file5.txt"},
			}},
			expectOptions: &deque[Value]{values: []Value{
				ValueOption{
					Option: cfg.options["z"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["x"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "file1.txt",
				},
				ValueOption{
					Option: cfg.options["verbose"],
					Tok:    scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["k"],
					Tok:    scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
					Value:  "file4.txt",
				},
				ValueOption{
					Option: cfg.options["file"],
					Tok:    scanner.OptionToken{Idx: 7, Prefix: "--", Name: "file"},
					Value:  "/dev/null",
				},
				ValueOption{
					Option: cfg.options["z"],
					Tok:    scanner.OptionToken{Idx: 9, Prefix: "-", Name: "z"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["http"],
					Tok:    scanner.OptionToken{Idx: 11, Prefix: "--", Name: "http"},
					Value:  "1.1",
				},
				ValueOption{
					Option: cfg.options["http"],
					Tok:    scanner.OptionToken{Idx: 12, Prefix: "--", Name: "http=2.0"},
					Value:  "2.0",
				},
			}},
			expectPositionals: &deque[Value]{values: []Value{
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
					Value: "file2.txt",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 10, Value: "file5.txt"},
					Value: "file5.txt",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 13, Value: "file5.txt"},
					Value: "file5.txt",
				},
			}},
			expectErrUnknown: false,
			expectErrArg:     false,
			expectErrNoArg:   false,
			expectPanic:      false,
		},

		{
			name:           "successful parsing without permutation",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: true,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
				scanner.PositionalArgumentToken{Idx: 2, Value: "file1.txt"},
				scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
				scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
				scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
				scanner.PositionalArgumentToken{Idx: 6, Value: "file4.txt"},
				scanner.OptionToken{Idx: 7, Prefix: "--", Name: "file"},
				scanner.PositionalArgumentToken{Idx: 8, Value: "/dev/null"},
				scanner.OptionToken{Idx: 9, Prefix: "-", Name: "z"},
				scanner.PositionalArgumentToken{Idx: 10, Value: "file5.txt"},
				scanner.OptionToken{Idx: 11, Prefix: "--", Name: "http"},
				scanner.OptionToken{Idx: 12, Prefix: "--", Name: "http=2.0"},
				scanner.PositionalArgumentToken{Idx: 13, Value: "file5.txt"},
			}},
			expectOptions: &deque[Value]{values: []Value{
				ValueOption{
					Option: cfg.options["z"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["x"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "file1.txt",
				},
				ValueOption{
					Option: cfg.options["verbose"],
					Tok:    scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
					Value:  "",
				},
			}},
			expectPositionals: &deque[Value]{values: []Value{
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
					Value: "file2.txt",
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
					Value: "-k",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 6, Value: "file4.txt"},
					Value: "file4.txt",
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 7, Prefix: "--", Name: "file"},
					Value: "--file",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 8, Value: "/dev/null"},
					Value: "/dev/null",
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 9, Prefix: "-", Name: "z"},
					Value: "-z",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 10, Value: "file5.txt"},
					Value: "file5.txt",
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 11, Prefix: "--", Name: "http"},
					Value: "--http",
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 12, Prefix: "--", Name: "http=2.0"},
					Value: "--http=2.0",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 13, Value: "file5.txt"},
					Value: "file5.txt",
				},
			}},
			expectErrUnknown: false,
			expectErrArg:     false,
			expectErrNoArg:   false,
		},

		{
			name:           "successful parsing with permutation with separator",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
				scanner.PositionalArgumentToken{Idx: 2, Value: "file1.txt"},
				scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
				scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
				scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
				scanner.PositionalArgumentToken{Idx: 6, Value: "file4.txt"},
				scanner.OptionToken{Idx: 7, Prefix: "--", Name: "file"},
				scanner.PositionalArgumentToken{Idx: 8, Value: "/dev/null"},
				scanner.OptionsArgumentsSeparatorToken{Idx: 9, Separator: "--"},
				scanner.OptionToken{Idx: 10, Prefix: "-", Name: "z"},
				scanner.PositionalArgumentToken{Idx: 11, Value: "file5.txt"},
				scanner.OptionToken{Idx: 12, Prefix: "--", Name: "http"},
				scanner.OptionToken{Idx: 13, Prefix: "--", Name: "http=2.0"},
				scanner.PositionalArgumentToken{Idx: 14, Value: "file5.txt"},
			}},
			expectOptions: &deque[Value]{values: []Value{
				ValueOption{
					Option: cfg.options["z"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["x"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "file1.txt",
				},
				ValueOption{
					Option: cfg.options["verbose"],
					Tok:    scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["k"],
					Tok:    scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
					Value:  "file4.txt",
				},
				ValueOption{
					Option: cfg.options["file"],
					Tok:    scanner.OptionToken{Idx: 7, Prefix: "--", Name: "file"},
					Value:  "/dev/null",
				},
			}},
			expectPositionals: &deque[Value]{values: []Value{
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
					Value: "file2.txt",
				},
				ValueOptionsArgumentsSeparator{
					Separator: "--",
					Tok:       scanner.OptionsArgumentsSeparatorToken{Idx: 9, Separator: "--"},
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 10, Prefix: "-", Name: "z"},
					Value: "-z",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 11, Value: "file5.txt"},
					Value: "file5.txt",
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 12, Prefix: "--", Name: "http"},
					Value: "--http",
				},
				ValuePositionalArgument{
					Tok:   scanner.OptionToken{Idx: 13, Prefix: "--", Name: "http=2.0"},
					Value: "--http=2.0",
				},
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 14, Value: "file5.txt"},
					Value: "file5.txt",
				},
			}},
			expectErrUnknown: false,
			expectErrArg:     false,
			expectErrNoArg:   false,
		},

		{
			name:           "nonexitent standalone option",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
				scanner.PositionalArgumentToken{Idx: 2, Value: "file1.txt"},
				scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
				scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
				scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
				scanner.PositionalArgumentToken{Idx: 6, Value: "file4.txt"},
				scanner.OptionToken{Idx: 7, Prefix: "--", Name: "output"}, // non-existent option
				scanner.PositionalArgumentToken{Idx: 8, Value: "/dev/null"},
				scanner.OptionsArgumentsSeparatorToken{Idx: 9, Separator: "--"},
				scanner.OptionToken{Idx: 10, Prefix: "-", Name: "z"},
				scanner.PositionalArgumentToken{Idx: 11, Value: "file5.txt"},
				scanner.OptionToken{Idx: 12, Prefix: "--", Name: "http"},
				scanner.OptionToken{Idx: 13, Prefix: "--", Name: "http=2.0"},
				scanner.PositionalArgumentToken{Idx: 14, Value: "file5.txt"},
			}},
			expectOptions: &deque[Value]{values: []Value{
				ValueOption{
					Option: cfg.options["z"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["x"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "file1.txt",
				},
				ValueOption{
					Option: cfg.options["verbose"],
					Tok:    scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["k"],
					Tok:    scanner.OptionToken{Idx: 5, Prefix: "-", Name: "k"},
					Value:  "file4.txt",
				},
			}},
			expectPositionals: &deque[Value]{values: []Value{
				ValuePositionalArgument{
					Tok:   scanner.PositionalArgumentToken{Idx: 4, Value: "file2.txt"},
					Value: "file2.txt",
				},
			}},
			expectErrUnknown: true,
			expectErrArg:     false,
			expectErrNoArg:   false,
		},

		{
			name:           "argument to standalone option that requires no argument",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
				scanner.PositionalArgumentToken{Idx: 2, Value: "file1.txt"},
				scanner.OptionToken{Idx: 3, Prefix: "--", Name: "verbose=true"}, // this option requires no argument
			}},
			expectOptions: &deque[Value]{values: []Value{
				ValueOption{
					Option: cfg.options["z"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["x"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "file1.txt",
				},
			}},
			expectPositionals: &deque[Value]{values: nil},
			expectErrUnknown:  false,
			expectErrArg:      false,
			expectErrNoArg:    true,
		},

		{
			name:           "no argument to standalone option that requires argument",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
				scanner.PositionalArgumentToken{Idx: 2, Value: "file1.txt"},
				scanner.OptionToken{Idx: 3, Prefix: "--", Name: "file"}, // this option requires argument
			}},
			expectOptions: &deque[Value]{values: []Value{
				ValueOption{
					Option: cfg.options["z"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "",
				},
				ValueOption{
					Option: cfg.options["x"],
					Tok:    scanner.OptionToken{Idx: 1, Prefix: "-", Name: "zx"},
					Value:  "file1.txt",
				},
			}},
			expectPositionals: &deque[Value]{values: nil},
			expectErrUnknown:  false,
			expectErrArg:      true,
			expectErrNoArg:    false,
		},

		{
			name:           "panic because we see the program name",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.ProgramNameToken{Idx: 1, Name: "curl"},
			}},
			expectOptions:     &deque[Value]{values: nil},
			expectPositionals: &deque[Value]{values: nil},
			expectErrUnknown:  false,
			expectErrArg:      false,
			expectErrNoArg:    false,
			expectPanic:       true, // we expect a panic because we see the program name token
		},

		{
			name:           "panic because of unhandled standalone option type",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "--", Name: "__panic"},
			}},
			expectOptions:     &deque[Value]{values: nil},
			expectPositionals: &deque[Value]{values: nil},
			expectErrUnknown:  false,
			expectErrArg:      false,
			expectErrNoArg:    false,
			expectPanic:       true, // we expect a panic because of an invalid option type
		},

		{
			name:           "panic because of unhandled groupable option type",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "-", Name: "_"},
			}},
			expectOptions:     &deque[Value]{values: nil},
			expectPositionals: &deque[Value]{values: nil},
			expectErrUnknown:  false,
			expectErrArg:      false,
			expectErrNoArg:    false,
			expectPanic:       true, // we expect a panic because of an invalid option type
		},

		{
			name:           "panic because of prefix not bound to any option",
			skipCase:       false,
			cfg:            cfg,
			disablePermute: false,
			input: &deque[scanner.Token]{values: []scanner.Token{
				scanner.OptionToken{Idx: 1, Prefix: "/", Name: "help"},
			}},
			expectOptions:     &deque[Value]{values: nil},
			expectPositionals: &deque[Value]{values: nil},
			expectErrUnknown:  false,
			expectErrArg:      false,
			expectErrNoArg:    false,
			expectPanic:       true, // we expect a panic because of an unbound prefix
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Possibly skip the test case if requested
			if tc.skipCase {
				t.Skip("skipping test case:", tc.name)
			}

			// Be prepared to handle a panic
			defer func() {
				r := recover()
				if got := (r != nil); got != tc.expectPanic {
					t.Fatal("expected", tc.expectPanic, "got", got, "panic is", r)
				}
			}()

			// Honor request to disable permutation unconditionally
			// so that we always set it to a known-good value
			cfg.parser.DisablePermute = tc.disablePermute

			// Run the function we are testing
			var (
				options     deque[Value]
				positionals deque[Value]
			)
			err := parse(tc.cfg, tc.input, &options, &positionals)

			// Check for errors
			switch {
			case tc.expectErrArg:
				var errvalue ErrOptionRequiresArgument
				if !errors.As(err, &errvalue) {
					t.Fatalf("expected ErrOptionRequiresArgument, got %T", err)
				}

			case tc.expectErrNoArg:
				var errvalue ErrOptionRequiresNoArgument
				if !errors.As(err, &errvalue) {
					t.Fatalf("expected ErrOptionRequiresNoArgument, got %T", err)
				}

			case tc.expectErrUnknown:
				var errval ErrUnknownOption
				if !errors.As(err, &errval) {
					t.Fatalf("cannot convert error to ErrUnknownOption: %T", err)
				}

			case err != nil:
				t.Fatal(err)
			}

			// Check for success
			if diff := cmp.Diff(tc.expectOptions.values, options.values); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(tc.expectPositionals.values, positionals.values); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
