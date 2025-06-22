// scanner.go - Command line scanner.
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package scanner provides low-level tokenization of command-line arguments.

The [*Scanner.Scan] method breaks command-line arguments into [Token]
based on configurable option prefixes and separators, allowing higher-level parsers
to implement custom parsing logic on top of the tokenized stream.

# Token Types

[*Scanner.Scan] produces these token types:

 1. [ProgramNameToken]: The program name (argv[0])

 2. [OptionToken]: Options started with any configured prefix (e.g., -v, --verbose, +trace)

 3. [SeparatorToken]: Special separators (e.g., -- to stop parsing)

 4. [ArgumentToken]: Everything else (positional arguments)

# Option Prefixes

The [*Scanner] is configured with the option prefixes to use when tokenizing
command-line arguments. Prefixes are sorted by length (longest first) to ensure
correct tokenization when prefixes overlap (e.g., "-" and "--").

This design allows building parsers for different command-line styles:

 1. GNU-style: "-", "--" (e.g., -v, --verbose)

 2. Dig-style: "-", "--", "+" (e.g., -v, --verbose, +trace)

 3. Windows-style: "/" (e.g., /v, /verbose)

 4. Go-style: "-" (e.g., -v, -verbose)

# Separators

The [*Scanner] can be configured to recognize and emit as a token the separator
to stop parsing options and treat all remaining arguments as positional.
However, note that it is not the [*Scanner] job to interpret semantics and
subsequent tokens will still be tokenized as options. This rule should instead
be implemented by higher-level parsers.

# Example

Given the "--" and "-" option prefixes and the "--" separator, the
following command line:

	command --verbose -- othercommand -v --trace file.txt

produces the following tokens:

 1. [ProgramNameToken] command
 2. [OptionToken] verbose
 3. [SeparatorToken] --
 4. [ArgumentToken] othercommand
 5. [OptionToken] v
 7. [OptionToken] trace
 8. [ArgumentToken] file.txt
*/
package scanner

import (
	"errors"
	"sort"
	"strings"
)

// Scanner is a command line scanner.
//
// We check for separators first. Then for option prefixes
// sorted by length (longest first).
type Scanner struct {
	// Prefixes contains the prefixes delimiting options.
	Prefixes []string

	// Separators contains the separators between option arguments.
	Separators []string
}

// Token is a token lexed by [*Scanner.Scan].
type Token interface {
	// String returns the string representation of the token.
	String() string
}

// OptionToken is a [Token] containing an option.
type OptionToken struct {
	// Index is the position in the original command line arguments.
	Index int

	// Prefix is the scanned prefix.
	Prefix string

	// Name is the parsed name.
	Name string
}

var _ Token = OptionToken{}

// String implements [Token].
func (tk OptionToken) String() string {
	return tk.Prefix + tk.Name
}

// ArgumentToken is a [Token] containing a positional argument.
type ArgumentToken struct {
	// Index is the position in the original command line arguments.
	Index int

	// Value is the parsed value.
	Value string
}

var _ Token = ArgumentToken{}

// String implements [Token].
func (tk ArgumentToken) String() string {
	return tk.Value
}

// SeparatorToken is a [Token] containing the separator between options and arguments.
type SeparatorToken struct {
	// Index is the position in the original command line arguments.
	Index int

	// Separator is the parsed separator.
	Separator string
}

var _ Token = SeparatorToken{}

// String implements [Token].
func (tk SeparatorToken) String() string {
	return tk.Separator
}

// ProgramNameToken is the program name [Token].
type ProgramNameToken struct {
	// Index is the position in the original command line arguments.
	Index int

	// Name is the program name.
	Name string
}

var _ Token = ProgramNameToken{}

// String implements [Token].
func (tk ProgramNameToken) String() string {
	return tk.Name
}

// ErrMissingProgramName is returned when the program name is missing. That is when
// [*Scanner.Scan] is passed an empty slice.
var ErrMissingProgramName = errors.New("missing program name")

// Scan scans the command line arguments and returns a list of [Token] or an error.
//
// The argv MUST include the program name as the first argument.
//
// This method does not mutate [*Scanner] and is safe to call concurrently.
//
// The only possible error is [ErrMissingProgramName].
func (sx *Scanner) Scan(argv []string) ([]Token, error) {
	// Create an empty list of tokens
	tokens := make([]Token, 0, len(argv))

	// Ensure there is at least the program name
	if len(argv) <= 0 {
		return nil, ErrMissingProgramName
	}

	// Save the program name
	tokens = append(tokens, ProgramNameToken{Index: 0, Name: argv[0]})
	argv = argv[1:]

	// Create sorted copy of prefixes (longest first)
	prefixes := make([]string, len(sx.Prefixes))
	copy(prefixes, sx.Prefixes)

	// Sort by length descending, then alphabetically for stability
	sort.SliceStable(prefixes, func(i, j int) bool {
		if len(prefixes[i]) == len(prefixes[j]) {
			return prefixes[i] < prefixes[j]
		}
		return len(prefixes[i]) > len(prefixes[j])
	})

	// Cycle through the remaining arguments
Loop:
	for idx, arg := range argv {
		// Calculate the actual index in the original args slice
		actual := idx + 1

		// Check for separators first
		for _, sep := range sx.Separators {
			if arg == sep {
				tokens = append(tokens, SeparatorToken{Index: actual, Separator: arg})
				continue Loop
			}
		}

		// Then, check for (sorted) prefixes
		for _, prefix := range prefixes {
			if strings.HasPrefix(arg, prefix) {
				tokens = append(tokens, OptionToken{Index: actual, Prefix: prefix, Name: arg[len(prefix):]})
				continue Loop
			}
		}

		// Everything else is an argument
		tokens = append(tokens, ArgumentToken{Index: actual, Value: arg})
	}

	return tokens, nil
}
