// flagset.go - FlagSet implementation for command-line flag parsing
// SPDX-License-Identifier: GPL-3.0-or-later

package flagemu

import (
	"errors"
	"fmt"

	"github.com/bassosimone/clip/internal/assert"
	"github.com/bassosimone/clip/pkg/parser"
)

// ErrorHandling defines the behavior of the flag set when an error occurs.
type ErrorHandling int

const (
	// ContinueOnError causes the flag set to continue parsing after an error occurs.
	ContinueOnError = ErrorHandling(iota)
)

// FlagSet represents a set of command-line flags.
//
// The zero value is invalid. Construct using [NewFlagSet].
type FlagSet struct {
	// args contains the positional arguments.
	args []string

	// flags contains flags to modify the parser's behavior.
	flags parser.Flags

	// progname contains the name of the program.
	progname string

	// valuesShort maps short flags (single character) to their values.
	valuesShort map[string]Value

	// valuesLong maps long flags (multi-character) to their values.
	valuesLong map[string]Value
}

// NewFlagSet creates a new flag set with the given program name and error handling.
func NewFlagSet(progname string, handling ErrorHandling) *FlagSet {
	return &FlagSet{
		args:        []string{},
		flags:       0,
		progname:    progname,
		valuesShort: make(map[string]Value),
		valuesLong:  make(map[string]Value),
	}
}

// SetNoPermute sets the flag to disable argument permutation. This flag will only
// apply to parsers created using [*FlagSet.NewParser] and when parsing directly
// using the [*FlagSet.Parse] method.
//
// This method MUST be invoked before parsing the arguments.
func (fs *FlagSet) SetNoPermute() {
	fs.flags |= parser.FlagNoPermute
}

// Bool creates a new boolean flag with the given name, default value, and usage.
//
// This method MUST be invoked before parsing the arguments.
func (fs *FlagSet) Bool(longName string, shortName byte, value bool, usage string) *bool {
	if longName != "" {
		fs.valuesLong[longName] = newBoolValue(&value)
	}
	if shortName != 0 {
		fs.valuesShort[string(shortName)] = newBoolValue(&value)
	}
	return &value
}

// String creates a new string flag with the given name, default value, and usage.
//
// This method MUST be invoked before parsing the arguments.
func (fs *FlagSet) String(longName string, shortName byte, value string, usage string) *string {
	if longName != "" {
		fs.valuesLong[longName] = newStringValue(&value)
	}
	if shortName != 0 {
		fs.valuesShort[string(shortName)] = newStringValue(&value)
	}
	return &value
}

// Args returns the positional arguments.
//
// This method MUST be invoked after parsing the arguments.
func (fs *FlagSet) Args() []string {
	return fs.args
}

// Parse parses the command line arguments.
//
// This method MUST be invoked after setting up the flags.
func (fs *FlagSet) Parse(arguments []string) error {
	return fs.ParseWithParser(fs.NewParser(), arguments)
}

// Parse parses the command line arguments. This method is useful when you need to customize the
// definition of the flags. For example, you may want to allow long flags to start with `+`. In such
// a use case, you construct a new parser using [*FlagSet.NewParser], customize the parser, and
// finally call [*FlagSet.ParseWithParser].
//
// This method MUST be invoked after setting up the flags.
func (fs *FlagSet) ParseWithParser(px *parser.Parser, arguments []string) error {
	// Generate the argv from the program name and the arguments.
	argv := append([]string{fs.progname}, arguments...)

	// Parse the argv
	res, err := px.Parse(argv)
	if err != nil {
		return err
	}

	// Save the results
	return fs.saveResults(res)
}

// NewParser constructs a new parser for the [*FlagSet]. This method is
// useful when you need to customize the definition of the flags. For
// example, you may want to allow long flags to start with `+`. In such
// a use case, you construct a new parser using [*FlagSet.NewParser], customize
// the parser, and finally call [*FlagSet.ParseWithParser].
//
// This method MUST be invoked after setting up the flags.
func (fs *FlagSet) NewParser() *parser.Parser {
	px := &parser.Parser{
		Flags:               fs.flags,
		LongOptions:         map[string]parser.OptionType{},
		LongOptionPrefixes:  []string{"--"},
		Separators:          []string{"--"},
		ShortOptions:        map[string]parser.OptionType{},
		ShortOptionPrefixes: []string{"-"},
	}
	// Populate short options from valuesShort
	for name, value := range fs.valuesShort {
		px.ShortOptions[name] = value.OptionType()
	}
	// Populate long options from valuesLong
	for name, value := range fs.valuesLong {
		px.LongOptions[name] = value.OptionType()
	}
	return px
}

// NewGoStyleParser constructs a new parser configured for Go-style command-line parsing.
// In Go-style parsing, all options use a single dash prefix (e.g., -v, -verbose).
// This is useful for creating Go-compatible command-line interfaces.
//
// Note that this method creates a parser with the FlagNoPermute flag set, which means that
// options are not permuted with non-option arguments, regardless of whether you did or did not
// configure this using [*FlagSet.SetNoPermute].
//
// This method MUST be invoked after setting up the flags.
func (fs *FlagSet) NewGoStyleParser() *parser.Parser {
	px := &parser.Parser{
		Flags:               parser.FlagNoPermute,
		LongOptions:         map[string]parser.OptionType{},
		LongOptionPrefixes:  []string{"-"},
		Separators:          []string{"--"},
		ShortOptions:        map[string]parser.OptionType{},
		ShortOptionPrefixes: []string{},
	}
	// For Go-style parsing, put all options in LongOptions since they all use "-" prefix
	for name, value := range fs.valuesShort {
		px.LongOptions[name] = value.OptionType()
	}
	for name, value := range fs.valuesLong {
		px.LongOptions[name] = value.OptionType()
	}
	return px
}

// ErrUnknownOption is returned when an option is not found in the specification.
var ErrUnknownOption = errors.New("unknown option")

// ErrUnexpectedItemType is returned when an unexpected item type is encountered.
var ErrUnexpectedItemType = errors.New("unexpected item type")

// findOption is a helper method that looks up an option by name in both
// short and long option maps.
func (fs *FlagSet) findOption(name string) (Value, bool) {
	if option, found := fs.valuesShort[name]; found {
		return option, true
	}
	if option, found := fs.valuesLong[name]; found {
		return option, true
	}
	return nil, false
}

// saveResults saves the parsed results back to the [*FlagSet].
func (fs *FlagSet) saveResults(res []parser.CommandLineItem) error {
	// Ensure the results are consistent
	assert.True(len(res) >= 0, "expected at least one result")
	_, ok := res[0].(parser.ProgramNameItem)
	assert.True(ok, "expected program name item")
	res = res[1:]

	// Process each item in the command line
	for _, item := range res {
		switch item := item.(type) {

		// Handle option
		case parser.OptionItem:
			option, found := fs.findOption(item.Name)
			if !found {
				return fmt.Errorf("%w: %s", ErrUnknownOption, item.Name)
			}
			if err := option.Set(item.Value); err != nil {
				return err
			}

		// Handle argument
		case parser.ArgumentItem:
			fs.args = append(fs.args, item.Value)

		// Otherwise it's an error
		default:
			return fmt.Errorf("%w: %T", ErrUnexpectedItemType, item)
		}
	}
	return nil
}
