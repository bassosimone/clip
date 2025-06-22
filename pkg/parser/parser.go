// parser.go - Command line parser.
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package parser provides high-level command-line argument parsing.

The [*Parser.Parse] method parses command-line arguments into structured
options and arguments implementing [CommandLineItem]. It supports various
command-line styles through configurable prefixes and option types.

# Supported Option Types

The parser supports these option types:

 1. [OptionTypeBool]: Boolean flags that can be true/false.

 2. [OptionTypeString]: Basically any other options.

Users of this package should parse the string values into proper data types
such as integers, floats, or custom types.

# Command Line Items

The parser produces structured [CommandLineItem]:

 1. [ProgramNameItem]: The program name (i.e., argv[0]).

 2. [OptionItem]: Parsed options with their values.

 3. [ArgumentItem]: Positional arguments.

# Parser Behavior

The parser handles various parsing behaviors:

 1. Option bundling: Short options like -abc parsed as -a -b -c

 2. Value assignment: Long options support --option=value syntax

 3. Argument consumption: String options consume the next argument if no value
    provided (e.g., -ffile.text is equivalent to -f file.text)

 4. Separator handling: separators stops option parsing (arguments after
    the separators are treated as positional)

 5. Permutation control: [FlagNoPermute] stops parsing at first non-option
    argument (i.e., command -v -- subcommand -s is parsed such that the
    subcommand and the -s are treated as positional arguments)

# Flexible Configuration

Different command-line styles can be implemented by configuring prefixes:

 1. GNU-style: short options starting with "-" and long options starting with "--"

 2. Dig-style: short options starting with "-" and long options starting with "--" or "+"

 3. Go-style: only long options starting with "-"

 4. Windows-style: only long options starting with "/"

# Example

Assuming [FlagNoPermute] is not set and GNU-style configuration,
the following command line:

	command -v --file=example.txt subcommand -sky

is parsed as:

 1. [ProgramNameItem] command
 2. [OptionItem] -v
 3. [OptionItem] --file=example.txt
 4. [OptionItem] -s
 5. [OptionItem] -k
 6. [OptionItem] -y
 7. [ArgumentItem] subcommand

Instead, assuming [FlagNoPermute] is set, it is parsed as:

 1. [ProgramNameItem] command
 2. [OptionItem] -v
 3. [OptionItem] --file=example.txt
 4. [OptionItem] subcommand
 5. [ArgumentItem] -sky

# Implementation Details

The parser delegates tokenization to the [scanner] package.
*/
package parser

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bassosimone/clip/internal/assert"
	"github.com/bassosimone/clip/pkg/scanner"
)

// OptionType represents the type of an option.
type OptionType int64

// These constants define the available option types.
const (
	// optionTypeNull is a special value indicating that an option does
	// not exist. This value is the zero value for simplicity.
	optionTypeNull = iota

	// OptionTypeString is a string option.
	OptionTypeString

	// OptionTypeBool is a boolean option.
	OptionTypeBool
)

// Flags contains flags modifying the parser behavior.
type Flags int64

// These constants define the available flags.
const (
	// FlagNoPermute disables permutation of options and arguments.
	FlagNoPermute = Flags(1 << iota)
)

// Parser is a command line parser.
type Parser struct {
	// Flags contains flags modifying the parser behavior.
	//
	// By default, the parser will permute options and arguments. Use
	// the [FlagNoPermute] flag to disable this.
	Flags Flags

	// LongOptions is a map of long options to their types.
	//
	// Long options are prefixed with one of the prefixes in [LongOptionPrefixes]
	// and are never grouped. Therefore, to implement Go style command-line
	// parsing, you need to use the "-" prefix and configure all options, including
	// those consisting of a single character (e.g., "-v") as long options.
	//
	// There are two ways to specify a long option:
	//
	//   - In the same token: "--option=value"
	//
	//   - In subsequent tokens: "--option" "value"
	//
	// In both cases, we parse this as:
	//
	//   - [OptionItem] name="option" value="value"
	//
	// Boolean long options are presence-based, meaning their presence indicates
	// true and their absence indicates false. For example, "--verbose" sets
	// the option to true.
	LongOptions map[string]OptionType

	// LongOptionPrefixes contains the prefixes for long options.
	LongOptionPrefixes []string

	// Separators contains the separators for options and arguments.
	//
	// Typically, this field is empty or contains "--".
	Separators []string

	// ShortOptions is a map of short options to their types.
	//
	// Short options are prefixed with one of the prefixes in [ShortOptionPrefixes].
	//
	// Boolean short options are presence-based, meaning their presence indicates
	// true and their absence indicates false. For example, "-v" sets the option
	// to true.
	//
	// Short options can be grouped (e.g., "-abc" is equivalent to "-a -b -c"). When
	// an option takes an argument, there are two distinct cases:
	//
	//   - Argument on a subsequent token: "-vf file.txt"
	//
	//   - Argument on the same token: "-vffile.txt"
	//
	// In both cases, we will parse this as follows:
	//
	//   - [OptionItem] name="v" value="true"
	//
	//   - [OptionItem] name="f" value="file.txt"
	ShortOptions map[string]OptionType

	// ShortOptionPrefixes contains the prefixes for short options.
	ShortOptionPrefixes []string
}

// CommandLineItem is an item present in the command line.
type CommandLineItem interface {
	// Strings returns the strings to append to a slice
	// to reconstruct the original command line.
	Strings() []string
}

// ProgramNameItem is the name of the program
type ProgramNameItem struct {
	// Name is the name of the program.
	Name string

	// Token is the token associated with the item.
	Token scanner.Token
}

var _ CommandLineItem = ProgramNameItem{}

// Strings implements [CommandLineItem].
func (p ProgramNameItem) Strings() []string {
	return []string{p.Name}
}

// OptionItem is an option present on the command line.
type OptionItem struct {
	// Name is the option name without any prefix (e.g., "verbose" or "v").
	Name string

	// Token is the token associated with the item.
	Token scanner.Token

	// Value is the option value. For boolean options this is always "true"
	// when the option is present.
	Value string

	// IsShort indicates whether this is a short option (e.g., -v vs --verbose).
	IsShort bool

	// Type indicates the type of the option (bool or string).
	Type OptionType

	// Prefix is the prefix used for the option name (e.g., "-" or "--").
	Prefix string
}

var _ CommandLineItem = OptionItem{}

// Strings implements [CommandLineItem].
func (o OptionItem) Strings() []string {
	// For boolean options, just return the option name
	if o.Type == OptionTypeBool {
		return []string{o.Prefix + o.Name}
	}

	// Otherwise, return the name and the value in two distinct tokens
	return []string{o.Prefix + o.Name, o.Value}
}

// ArgumentItem is an argument present on the command line.
type ArgumentItem struct {
	// Token is the token associated with the item.
	Token scanner.Token

	// Value is the argument value.
	Value string
}

var _ CommandLineItem = ArgumentItem{}

// Strings implements [CommandLineItem].
func (a ArgumentItem) Strings() []string {
	return []string{a.Value}
}

// Parse parses the command line arguments.
//
// This method does not mutate [*Parser] and is safe to call concurrently.
//
// The argv MUST include the program name as the first argument.
//
// The possible errors are:
//
//  1. [scanner.ErrMissingProgramName]
//  2. [ErrInvalidOptionValue]
//  3. [ErrOptionRequiresValue]
//  4. [ErrUnknownOption]
func (px *Parser) Parse(argv []string) ([]CommandLineItem, error) {
	// Create the initial empty list of items
	rv := []CommandLineItem{}

	// Create the argv scanner and configure the prefixes
	sx := &scanner.Scanner{
		Prefixes:   []string{},
		Separators: px.Separators,
	}
	sx.Prefixes = append(sx.Prefixes, px.LongOptionPrefixes...)
	sx.Prefixes = append(sx.Prefixes, px.ShortOptionPrefixes...)

	// Scan the argv
	tokens, err := sx.Scan(argv)
	if err != nil {
		return nil, err
	}

	// Remember the program name and advance
	//
	// Here we assert because the scanner guarantees that the
	// first token is a program name token
	assert.True(len(tokens) >= 1, "expected at least one token")
	pname, ok := tokens[0].(scanner.ProgramNameToken)
	assert.True(ok, "expected program name token")
	rv = append(rv, ProgramNameItem{Name: pname.Name, Token: tokens[0]})
	tokens = tokens[1:]

	// Process the options
	rv, err = px.parse(tokens, rv)
	if err != nil {
		return nil, err
	}
	return rv, nil
}

func (px *Parser) parse(tokens []scanner.Token, rv []CommandLineItem) ([]CommandLineItem, error) {
	// We start parsing and stop when we see a separator
	parse := true

	for len(tokens) > 0 {
		// Get the current token and advance argv
		cur := tokens[0]
		tokens = tokens[1:]

		// Stop parsing in any case when we encounter a separator
		if _, ok := cur.(scanner.SeparatorToken); ok {
			parse = false
			continue
		}

		// If we're not parsing, just collect the argument
		if !parse {
			rv = append(rv, ArgumentItem{Token: cur, Value: cur.String()})
			continue
		}

		// If the argument is not an option we collect it.
		//
		// Note: stop parsing if we're not permuting.
		curopt, ok := cur.(scanner.OptionToken)
		if !ok {
			rv = append(rv, ArgumentItem{Token: cur, Value: cur.String()})
			if (px.Flags & FlagNoPermute) != 0 {
				parse = false
			}
			continue
		}

		// Select the proper parser to use
		var pf optionParser
		if slices.Contains(px.LongOptionPrefixes, curopt.Prefix) {
			pf = px.parseLong
		} else {
			pf = px.parseShort
		}

		// Parse either a long or short option
		var err error
		tokens, rv, err = pf(tokens, rv, curopt)

		// Handle parse errors
		if err != nil {
			return nil, err
		}
	}

	// Return success when out of tokens
	return rv, nil
}

// ErrUnknownOption is returned when an option is not found in the [*Parser] config.
var ErrUnknownOption = errors.New("unknown option")

// ErrOptionRequiresValue is returned when an option requires a value but none is provided.
var ErrOptionRequiresValue = errors.New("option requires a value")

// ErrInvalidOptionValue is returned when an option value is invalid.
var ErrInvalidOptionValue = errors.New("invalid option value")

type optionParser func(tokens []scanner.Token, rv []CommandLineItem,
	cur scanner.OptionToken) ([]scanner.Token, []CommandLineItem, error)

func (px *Parser) parseLong(tokens []scanner.Token, rv []CommandLineItem,
	cur scanner.OptionToken) ([]scanner.Token, []CommandLineItem, error) {
	// The option may contain a value, account for this
	var optname, optvalue string
	index := strings.Index(cur.Name, "=")
	if index > 0 {
		optname = cur.Name[:index]
		optvalue = cur.Name[index+1:]
	} else {
		optname = cur.Name
	}

	// Determine what to do based on the option kind
	optkind := px.LongOptions[optname]
	switch optkind {

	// Handle the case of boolean option
	case OptionTypeBool:
		if optvalue != "" {
			return nil, nil, fmt.Errorf("%w for option %s: %s", ErrInvalidOptionValue, optname, optvalue)
		}
		rv = append(rv, OptionItem{
			Name:    optname,
			Token:   cur,
			Value:   "true",
			IsShort: false,
			Type:    OptionTypeBool,
			Prefix:  cur.Prefix,
		})
		return tokens, rv, nil

	// Handle the case of a string option
	case OptionTypeString:
		// The value has been provided in the same token
		if optvalue != "" {
			rv = append(rv, OptionItem{
				Name:    optname,
				Token:   cur,
				Value:   optvalue,
				IsShort: false,
				Type:    OptionTypeString,
				Prefix:  cur.Prefix,
			})
			return tokens, rv, nil
		}

		// Otherwise try to use the next entry in the argv
		return px.getOptionValueFromNextToken(tokens, rv, cur, optname, false, OptionTypeString)

	// Otherwise, the option does not exist
	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrUnknownOption, optname)
	}
}

func (px *Parser) parseShort(tokens []scanner.Token, rv []CommandLineItem,
	cur scanner.OptionToken) ([]scanner.Token, []CommandLineItem, error) {
	// Process each character in the option string
	optstr := cur.Name
	for len(optstr) > 0 {
		// Get the character and advance
		optname := string(optstr[0])
		optstr = optstr[1:]

		// Determine what to do based on the option kind
		optkind := px.ShortOptions[optname]
		switch optkind {

		// If the option does not need an argument, advance
		case OptionTypeBool:
			rv = append(rv, OptionItem{
				Name:    optname,
				Token:   cur,
				Value:   "true",
				IsShort: true,
				Type:    OptionTypeBool,
				Prefix:  cur.Prefix,
			})
			continue

		// If the option needs an argument, fetch it
		case OptionTypeString:
			// GNU getopt compatible short options processing: just consume the remainder
			if len(optstr) > 0 {
				rv = append(rv, OptionItem{
					Name:    optname,
					Token:   cur,
					Value:   optstr,
					IsShort: true,
					Type:    OptionTypeString,
					Prefix:  cur.Prefix,
				})
				return tokens, rv, nil
			}

			// Otherwise try to use the next entry in the argv
			return px.getOptionValueFromNextToken(tokens, rv, cur, optname, true, OptionTypeString)

		// Otherwise, it does not exist
		default:
			return nil, nil, fmt.Errorf("%w: %s", ErrUnknownOption, optname)
		}
	}

	// Return the updated arguments vector
	return tokens, rv, nil
}

func (px *Parser) getOptionValueFromNextToken(tokens []scanner.Token, rv []CommandLineItem,
	cur scanner.OptionToken, optname string, isShort bool, optType OptionType) ([]scanner.Token, []CommandLineItem, error) {
	// Make sure there is at least one token left
	if len(tokens) < 1 {
		return nil, nil, fmt.Errorf("%w: %s", ErrOptionRequiresValue, optname)
	}

	// Make sure the token is either an argument or an option
	//
	// Specifically, we don't want the program name or the separator
	// or anything else to be a valid option value
	switch tokens[0].(type) {
	case scanner.ArgumentToken:
	case scanner.OptionToken:
	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrInvalidOptionValue, optname)
	}

	// Add the option value to the list of parsed items
	rv = append(rv, OptionItem{
		Name:    optname,
		Token:   cur,
		Value:   tokens[0].String(),
		IsShort: isShort,
		Type:    optType,
		Prefix:  cur.Prefix,
	})

	// Advance the token pointer
	tokens = tokens[1:]

	// Return the updated tokens and results
	return tokens, rv, nil
}
