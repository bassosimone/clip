// parser.go - parser definition.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"fmt"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/scanner"
)

// Parser is a command line parser.
type Parser struct {
	// DisablePermute optionally disables permuting options and arguments.
	//
	// Consider the following command line:
	//
	// 	curl https://www.google.com/ -H 'Host: google.com'
	//
	// The default behavior is to permute this to:
	//
	// 	curl -H 'Host: google.com' https://www.google.com/
	//
	// However, when DisablePermute is true, we keep the command
	// line unmodified. While permuting is a nice-to-have property
	// in general, consider instead the following case:
	//
	// 	multirepo foreach -kx git status -v
	//
	// With permutation, this command line would become:
	//
	// 	multirepo foreach git status -kv -v
	//
	// This is not the desired behavior if the foreach command
	// takes another command and its options as arguments.
	//
	// To make the above command line work with permutation, a
	// user would instead need to write this:
	//
	// 	multirepo foreach -kv -- git status -v
	//
	// By setting DisablePermute to true, the `--` separator
	// becomes unnecessary and the UX is improved.
	DisablePermute bool

	// MaxPositionalArguments is the maximum number of positional
	// arguments allowed by the parser. The default is zero, meaning
	// that the parser won't accept more than zero positionals.
	MaxPositionalArguments int

	// MinPositionalArguments is the minimum number of positional
	// arguments allowed by the parser. The default is zero, meaning
	// that the parser won't accept less than zero positionals.
	MinPositionalArguments int

	// OptionsArgumentsSeparator is the optional separator that terminates
	// the parsing of options, treating all remaining tokens in the command
	// line as positional arguments. The default is empty, meaning that
	// the parser will always parse all the available options.
	OptionsArgumentsSeparator string

	// Options contains the optional options configured for this parser.
	//
	// When parsing, we will ensure there are no duplicate option names or
	// ambiguous separators across all options.
	//
	// If you don't set this field, the parser will automatically
	// configure itself to parse GNU-style options, meaning that it
	// will use `-` as the prefix for short options and `--` as
	// the prefix for long options. No options will be defined so
	// any option will be considered unknown.
	Options []*Option
}

// ErrTooFewPositionalArguments is returned when the number of positional
// arguments is less than the configured minimum.
type ErrTooFewPositionalArguments struct {
	// Min is the minimum number of positional arguments required.
	Min int

	// Have is the number of positional arguments provided.
	Have int
}

var _ error = ErrTooFewPositionalArguments{}

// Error returns a string representation of this error.
func (err ErrTooFewPositionalArguments) Error() string {
	return fmt.Sprintf("too few positional arguments: expected at least %d, got %d", err.Min, err.Have)
}

// ErrTooManyPositionalArguments is returned when the number of positional
// arguments is greater than the configured maximum.
type ErrTooManyPositionalArguments struct {
	// Max is the maximum number of positional arguments allowed.
	Max int

	// Have is the number of positional arguments provided.
	Have int
}

var _ error = ErrTooManyPositionalArguments{}

// Error returns a string representation of this error.
func (err ErrTooManyPositionalArguments) Error() string {
	return fmt.Sprintf("too many positional arguments: expected at most %d, got %d", err.Max, err.Have)
}

// Parse parses the command line arguments.
//
// This method does not mutate [*Parser] and is safe to call concurrently.
//
// The argv MUST include the program name as the first argument.
func (px *Parser) Parse(argv []string) ([]Value, error) {
	// Create a new configuration for the parser.
	cfg, err := newConfig(px)
	if err != nil {
		return nil, err
	}

	// Preflight the command line arguments searching for early options.
	if value, found := searchEarly(px, argv); found {
		// If we found an early option, return it immediately.
		result := make([]Value, 0, 2)
		result = append(result, ValueProgramName{Name: argv[0]})
		result = append(result, value)
		return result, nil
	}

	// Create scanner for the parser.
	sx := &scanner.Scanner{
		Separators: []string{},
		Prefixes:   []string{},
	}
	if len(px.OptionsArgumentsSeparator) > 0 {
		sx.Separators = append(sx.Separators, px.OptionsArgumentsSeparator)
	}
	for prefix := range cfg.prefixes {
		sx.Prefixes = append(sx.Prefixes, prefix)
	}

	// Tokenize the command line arguments.
	tokens, err := sx.Scan(argv)
	if err != nil {
		return nil, err
	}

	// Remember the program name and advance
	//
	// Here we assert because the scanner guarantees that the
	// first token is a program name token
	assert.True(len(tokens) >= 1, "expected at least one token")
	pnameToken, ok := tokens[0].(scanner.ProgramNameToken)
	assert.True(ok, "expected program name token")
	programName := ValueProgramName{Name: pnameToken.Name, Tok: pnameToken}
	tokens = tokens[1:]

	// Create a deque containing the values to parse.
	input := &deque[scanner.Token]{values: tokens}

	// Parse the command line.
	var (
		positionals = &deque[Value]{}
		options     = &deque[Value]{}
	)
	if err := parse(cfg, input, options, positionals); err != nil {
		return nil, err
	}

	// Ensure this stage has emptied the input deque.
	assert.True(input.Empty(), "expected no unparsed tokens left after parsing standalone options")

	// Ensure the number of positional arguments is within the limits.
	if len(positionals.values) < px.MinPositionalArguments {
		return nil, ErrTooFewPositionalArguments{
			Min:  px.MinPositionalArguments,
			Have: len(positionals.values),
		}
	}
	if len(positionals.values) > px.MaxPositionalArguments {
		return nil, ErrTooManyPositionalArguments{
			Max:  px.MaxPositionalArguments,
			Have: len(positionals.values),
		}
	}

	// Create the result slice.
	result := permute(cfg, programName, options.values, positionals.values)
	return result, nil
}
