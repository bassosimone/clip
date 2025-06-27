// config.go - parser config.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"fmt"

	"github.com/bassosimone/clip/pkg/scanner"
)

// ErrAmbiguousPrefix indicates that the options contain ambiguous prefixes.
type ErrAmbiguousPrefix struct {
	// Prefix is the prefix that is used for both standalone and groupable options.
	Prefix string
}

var _ error = ErrAmbiguousPrefix{}

// Error returns a string representation of this error.
func (err ErrAmbiguousPrefix) Error() string {
	return fmt.Sprintf("prefix %q is used for both standalone and groupable options", err.Prefix)
}

// ErrMultipleOptionWithSameName indicates that there are multiple options with the same name.
type ErrMultipleOptionsWithSameName struct {
	// Name is the name of the option that appears multiple times.
	Name string

	// Options is a slice of options with the same name.
	Options []*Option
}

var _ error = ErrMultipleOptionsWithSameName{}

// Error returns a string representation of this error.
func (err ErrMultipleOptionsWithSameName) Error() string {
	return fmt.Sprintf("multiple options with %q name", err.Name)
}

// ErrTooLongGroupableOptionName indicates that a groupable option name is longer than one byte.
type ErrTooLongGroupableOptionName struct {
	Option *Option
}

var _ error = ErrTooLongGroupableOptionName{}

// Error returns a string representation of this error.
func (err ErrTooLongGroupableOptionName) Error() string {
	return fmt.Sprintf("groupable option names should be a single byte, found: %+v", err.Option)
}

// ErrEmptyOptionName indicates that an option name is empty.
type ErrEmptyOptionName struct {
	// Option is the option with the empty name.
	Option *Option
}

var _ error = ErrEmptyOptionName{}

// Error returns a string representation of this error.
func (err ErrEmptyOptionName) Error() string {
	return fmt.Sprintf("option name cannot be empty: %+v", err.Option)
}

// ErrEmptyOptionPrefix indicates that an option prefix is empty.
type ErrEmptyOptionPrefix struct {
	// Option is the option with the empty prefix.
	Option *Option
}

var _ error = ErrEmptyOptionPrefix{}

// Error returns a string representation of this error.
func (err ErrEmptyOptionPrefix) Error() string {
	return fmt.Sprintf("option prefix cannot be empty: %+v", err.Option)
}

// ErrUnknownOption indicates that an option is unknown.
type ErrUnknownOption struct {
	// Name is the name of the unknown option.
	Name string

	// Prefix is the prefix of the unknown option.
	Prefix string

	// Token is the token of the unknown option.
	Token scanner.Token
}

var _ error = ErrUnknownOption{}

// Error returns a string representation of this error.
func (err ErrUnknownOption) Error() string {
	return fmt.Sprintf("unknown option: %s%s", err.Prefix, err.Name)
}

type config struct {
	// options maps option names to options
	options map[string]*Option

	// parser is the parent parser
	parser *Parser

	// prefixes maps a prefix to its option type
	prefixes map[string]OptionType
}

// newConfig creates a new [*config] instance.
func newConfig(px *Parser) (*config, error) {
	// Make sure that groupable options have a single-byte name.
	for _, opt := range px.Options {
		if len(opt.Name) > 1 && opt.Type.isGroupable() {
			return nil, ErrTooLongGroupableOptionName{opt}
		}
	}

	// Make sure each option name appears exactly once to avoid ambiguity.
	names := make(map[string][]*Option)
	for _, opt := range px.Options {
		switch {
		case len(opt.Name) <= 0:
			return nil, ErrEmptyOptionName{opt}
		case len(opt.Prefix) <= 0:
			return nil, ErrEmptyOptionPrefix{opt}
		default:
			names[opt.Name] = append(names[opt.Name], opt)
		}
	}
	for name, options := range names {
		if len(options) != 1 {
			return nil, ErrMultipleOptionsWithSameName{Name: name, Options: options}
		}
	}

	// Collect unique prefixes, ensure they are used consistently across
	// standalone and groupable options, and configure the scanner for
	// scanning them. Note that we treat the early options as a special case
	// since they are checked ahead of proper parsing.
	prefixes := make(map[string]OptionType)
	for _, opt := range px.Options {
		switch {
		case opt.Type.isGroupable():
			prefixes[opt.Prefix] |= optionKindGroupable
		case opt.Type.isStandalone():
			prefixes[opt.Prefix] |= optionKindStandalone
		}
	}
	offending := optionKindGroupable | optionKindStandalone
	for prefix, flags := range prefixes {
		if (flags & offending) == offending {
			return nil, ErrAmbiguousPrefix{prefix}
		}
	}

	// If no prefixes have been defined, we assume that the user wants
	// a GNU-style parser, so we add the GNU-style prefixes.
	//
	// This is an edge case. Usually you want to use a parser to
	// parse *some* options but, anyway, it can happen.
	if len(prefixes) <= 0 {
		prefixes["-"] = optionKindGroupable
		prefixes["--"] = optionKindStandalone
	}

	// Create a map between option names and their spec.
	options := make(map[string]*Option)
	for _, opt := range px.Options {
		options[opt.Name] = opt
	}

	// Build the worker instance.
	cfg := &config{
		parser:   px,
		prefixes: prefixes,
		options:  options,
	}

	// Return the worker instance.
	return cfg, nil
}

func (cfg config) disablePermute() bool {
	return cfg.parser.DisablePermute
}

func (cfg *config) findOption(tok scanner.OptionToken, optname string, kind OptionType) (*Option, error) {
	option := cfg.options[optname]
	if option == nil || option.Prefix != tok.Prefix || (option.Type&kind) == 0 {
		err := ErrUnknownOption{Name: optname, Prefix: tok.Prefix, Token: tok}
		return nil, err
	}
	return option, nil
}
