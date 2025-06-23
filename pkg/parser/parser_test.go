// parser_test.go - Tests for command line parser.
// SPDX-License-Identifier: GPL-3.0-or-later

package parser

import (
	"errors"
	"testing"

	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/google/go-cmp/cmp"
)

func TestProgramNameItemStrings(t *testing.T) {
	testCases := []struct {
		name     string
		item     ProgramNameItem
		expected []string
	}{
		{
			name: "simple program name",
			item: ProgramNameItem{
				Name:  "myprogram",
				Token: scanner.ProgramNameToken{Name: "myprogram"},
			},
			expected: []string{"myprogram"},
		},
		{
			name: "program name with path",
			item: ProgramNameItem{
				Name:  "./bin/myprogram",
				Token: scanner.ProgramNameToken{Name: "./bin/myprogram"},
			},
			expected: []string{"./bin/myprogram"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.item.Strings()
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("ProgramNameItem.Strings() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOptionItemStrings(t *testing.T) {
	testCases := []struct {
		name     string
		item     OptionItem
		expected []string
	}{
		{
			name: "short boolean option",
			item: OptionItem{
				Name:    "v",
				Value:   "true",
				Token:   scanner.OptionToken{Name: "v", Prefix: "-"},
				IsShort: true,
				Type:    OptionTypeBool,
				Prefix:  "-",
			},
			expected: []string{"-v"},
		},
		{
			name: "long boolean option",
			item: OptionItem{
				Name:    "verbose",
				Value:   "true",
				Token:   scanner.OptionToken{Name: "verbose", Prefix: "--"},
				IsShort: false,
				Type:    OptionTypeBool,
				Prefix:  "--",
			},
			expected: []string{"--verbose"},
		},
		{
			name: "short string option",
			item: OptionItem{
				Name:    "f",
				Value:   "example.txt",
				Token:   scanner.OptionToken{Name: "f", Prefix: "-"},
				IsShort: true,
				Type:    OptionTypeString,
				Prefix:  "-",
			},
			expected: []string{"-f", "example.txt"},
		},
		{
			name: "long string option",
			item: OptionItem{
				Name:    "file",
				Value:   "example.txt",
				Token:   scanner.OptionToken{Name: "file", Prefix: "--"},
				IsShort: false,
				Type:    OptionTypeString,
				Prefix:  "--",
			},
			expected: []string{"--file", "example.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.item.Strings()
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("OptionItem.Strings() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestArgumentItemStrings(t *testing.T) {
	testCases := []struct {
		name     string
		item     ArgumentItem
		expected []string
	}{
		{
			name: "simple argument",
			item: ArgumentItem{
				Value: "file.txt",
				Token: scanner.ArgumentToken{Value: "file.txt"},
			},
			expected: []string{"file.txt"},
		},
		{
			name: "argument with spaces",
			item: ArgumentItem{
				Value: "file with spaces.txt",
				Token: scanner.ArgumentToken{Value: "file with spaces.txt"},
			},
			expected: []string{"file with spaces.txt"},
		},
		{
			name: "empty argument",
			item: ArgumentItem{
				Value: "",
				Token: scanner.ArgumentToken{Value: ""},
			},
			expected: []string{""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.item.Strings()
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("ArgumentItem.Strings() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSeparatorItemStrings(t *testing.T) {
	testCases := []struct {
		name     string
		item     SeparatorItem
		expected []string
	}{
		{
			name: "double dash separator",
			item: SeparatorItem{
				Separator: "--",
				Token:     scanner.SeparatorToken{Separator: "--"},
			},
			expected: []string{"--"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.item.Strings()
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("SeparatorItem.Strings() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParserParse(t *testing.T) {
	testCases := []struct {
		name    string
		parser  Parser
		argv    []string
		want    []CommandLineItem
		wantErr error
	}{
		{
			name:    "missing program name",
			parser:  Parser{},
			argv:    []string{},
			wantErr: scanner.ErrMissingProgramName,
		},

		{
			name:   "program name only",
			parser: Parser{},
			argv:   []string{"myprogram"},
			want: []CommandLineItem{
				ProgramNameItem{Name: "myprogram", Token: scanner.ProgramNameToken{Name: "myprogram"}},
			},
		},

		{
			name: "GNU-style with short options",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"v": OptionTypeBool,
					"f": OptionTypeString,
				},
			},
			argv: []string{"myprogram", "-v", "-f", "file.txt"},
			want: []CommandLineItem{
				ProgramNameItem{Name: "myprogram", Token: scanner.ProgramNameToken{Name: "myprogram"}},
				OptionItem{
					Name:    "v",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 1, Name: "v", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeBool,
					Prefix:  "-",
				},
				OptionItem{
					Name:    "f",
					Value:   "file.txt",
					Token:   scanner.OptionToken{Index: 2, Name: "f", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeString,
					Prefix:  "-",
				},
			},
		},

		{
			name: "GNU-style with long options",
			parser: Parser{
				LongOptionPrefixes: []string{"--"},
				LongOptions: map[string]OptionType{
					"verbose": OptionTypeBool,
					"file":    OptionTypeString,
				},
			},
			argv: []string{"myprogram", "target.txt", "--verbose", "--file=test.txt"},
			want: []CommandLineItem{
				ProgramNameItem{Name: "myprogram", Token: scanner.ProgramNameToken{Name: "myprogram"}},
				OptionItem{
					Name:    "verbose",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 2, Name: "verbose", Prefix: "--"},
					IsShort: false,
					Type:    OptionTypeBool,
					Prefix:  "--",
				},
				OptionItem{
					Name:    "file",
					Value:   "test.txt",
					Token:   scanner.OptionToken{Index: 3, Name: "file=test.txt", Prefix: "--"},
					IsShort: false,
					Type:    OptionTypeString,
					Prefix:  "--",
				},
				ArgumentItem{
					Value: "target.txt",
					Token: scanner.ArgumentToken{Index: 1, Value: "target.txt"},
				},
			},
		},

		{
			name: "options with separator",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"v": OptionTypeBool,
				},
				Separators: []string{"--"},
			},
			argv: []string{"myprogram", "-v", "--", "-v", "file.txt"},
			want: []CommandLineItem{
				ProgramNameItem{Name: "myprogram", Token: scanner.ProgramNameToken{Name: "myprogram"}},
				OptionItem{
					Name:    "v",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 1, Name: "v", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeBool,
					Prefix:  "-",
				},
				SeparatorItem{
					Separator: "--",
					Token:     scanner.SeparatorToken{Index: 2, Separator: "--"},
				},
				ArgumentItem{
					Value: "-v",
					Token: scanner.OptionToken{Index: 3, Name: "v", Prefix: "-"},
				},
				ArgumentItem{
					Value: "file.txt",
					Token: scanner.ArgumentToken{Index: 4, Value: "file.txt"},
				},
			},
		},

		{
			name: "no permute flag",
			parser: Parser{
				Flags:               FlagNoPermute,
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"v": OptionTypeBool,
				},
			},
			argv: []string{"myprogram", "-v", "file.txt", "-v"},
			want: []CommandLineItem{
				ProgramNameItem{Name: "myprogram", Token: scanner.ProgramNameToken{Name: "myprogram"}},
				OptionItem{
					Name:    "v",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 1, Name: "v", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeBool,
					Prefix:  "-",
				},
				ArgumentItem{
					Value: "file.txt",
					Token: scanner.ArgumentToken{Index: 2, Value: "file.txt"},
				},
				ArgumentItem{
					Value: "-v",
					Token: scanner.OptionToken{Index: 3, Name: "v", Prefix: "-"},
				},
			},
		},

		{
			name: "unknown short option",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"v": OptionTypeBool,
				},
			},
			argv:    []string{"myprogram", "-x"},
			wantErr: ErrUnknownOption,
		},

		{
			name: "unknown long option",
			parser: Parser{
				LongOptionPrefixes: []string{"--"},
				LongOptions: map[string]OptionType{
					"verbose": OptionTypeBool,
				},
			},
			argv:    []string{"myprogram", "--unknown"},
			wantErr: ErrUnknownOption,
		},

		{
			name: "missing value for string option",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"f": OptionTypeString,
				},
			},
			argv:    []string{"myprogram", "-f"},
			wantErr: ErrOptionRequiresValue,
		},

		{
			name: "no value for short string option",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"f": OptionTypeString,
				},
			},
			argv:    []string{"myprogram", "-f"},
			wantErr: ErrOptionRequiresValue,
		},

		{
			name: "no value for long string option",
			parser: Parser{
				LongOptionPrefixes: []string{"-"},
				LongOptions: map[string]OptionType{
					"file": OptionTypeString,
				},
			},
			argv:    []string{"myprogram", "-file"},
			wantErr: ErrOptionRequiresValue,
		},

		{
			name: "separator as value for short string option",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"f": OptionTypeString,
				},
				Separators: []string{"--"},
			},
			argv:    []string{"myprogram", "-f", "--"},
			wantErr: ErrInvalidOptionValue,
		},

		{
			name: "separator as value for long string option",
			parser: Parser{
				LongOptionPrefixes: []string{"-"},
				LongOptions: map[string]OptionType{
					"file": OptionTypeString,
				},
				Separators: []string{"--"},
			},
			argv:    []string{"myprogram", "-file", "--"},
			wantErr: ErrInvalidOptionValue,
		},

		{
			name: "attempt to set value for boolean option",
			parser: Parser{
				LongOptionPrefixes: []string{"--"},
				LongOptions: map[string]OptionType{
					"verbose": OptionTypeBool,
				},
			},
			argv:    []string{"myprogram", "--verbose=true"},
			wantErr: ErrInvalidOptionValue,
		},

		{
			name: "option bundling",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"a": OptionTypeBool,
					"b": OptionTypeBool,
					"c": OptionTypeBool,
				},
			},
			argv: []string{"myprogram", "-abc"},
			want: []CommandLineItem{
				ProgramNameItem{Name: "myprogram", Token: scanner.ProgramNameToken{Name: "myprogram"}},
				OptionItem{
					Name:    "a",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 1, Name: "abc", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeBool,
					Prefix:  "-",
				},
				OptionItem{
					Name:    "b",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 1, Name: "abc", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeBool,
					Prefix:  "-",
				},
				OptionItem{
					Name:    "c",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 1, Name: "abc", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeBool,
					Prefix:  "-",
				},
			},
		},

		{
			name: "option bundling with string option",
			parser: Parser{
				ShortOptionPrefixes: []string{"-"},
				ShortOptions: map[string]OptionType{
					"v": OptionTypeBool,
					"f": OptionTypeString,
				},
			},
			argv: []string{"myprogram", "-vffile.txt"},
			want: []CommandLineItem{
				ProgramNameItem{Name: "myprogram", Token: scanner.ProgramNameToken{Name: "myprogram"}},
				OptionItem{
					Name:    "v",
					Value:   "true",
					Token:   scanner.OptionToken{Index: 1, Name: "vffile.txt", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeBool,
					Prefix:  "-",
				},
				OptionItem{
					Name:    "f",
					Value:   "file.txt",
					Token:   scanner.OptionToken{Index: 1, Name: "vffile.txt", Prefix: "-"},
					IsShort: true,
					Type:    OptionTypeString,
					Prefix:  "-",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.parser.Parse(tc.argv)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestErrUnknownOptionContext(t *testing.T) {
	// create an error that is actually an [ErrUnknownOptionContext] instance
	var err error = ErrUnknownOptionContext{
		OptionName: "help",
		IsShort:    true,
		Token:      scanner.OptionToken{Index: 1, Name: "help", Prefix: "-"},
	}

	// Ensure that [errors.Is] continues to see this as [ErrUnknownOption]
	if !errors.Is(err, ErrUnknownOption) {
		t.Fatalf("expected error %v, got %v", ErrUnknownOption, err)
	}

	// Ensure that we can unwrap with [errors.As]
	var unknownOptionContext ErrUnknownOptionContext
	if !errors.As(err, &unknownOptionContext) {
		t.Fatalf("expected error %T, got %v", unknownOptionContext, err)
	}
	if diff := cmp.Diff(err, unknownOptionContext); diff != "" {
		t.Fatalf("ErrUnknownOptionContext mismatch (-want +got):\n%s", diff)
	}

	// Ensure that the error message is like before
	got := err.Error()
	expect := "unknown option: help"
	if diff := cmp.Diff(got, expect); diff != "" {
		t.Fatalf("ErrUnknownOptionContext mismatch (-want +got):\n%s", diff)
	}
}
