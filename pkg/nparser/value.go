// value.go - parsed value.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"fmt"
	"slices"

	"github.com/bassosimone/clip/pkg/scanner"
)

// Value is a value parsed by the [*Parser].
type Value interface {
	// Strings returns the strings to append to a slice
	// to reconstruct the original command line.
	Strings() []string

	// Token returns the scanner token.
	Token() scanner.Token
}

// ValueProgramName is a [Value] containing the parsed program name.
type ValueProgramName struct {
	// Name is the name of the program.
	Name string

	// Tok is the token associated with the value.
	Tok scanner.Token
}

var _ Value = ValueProgramName{}

// Strings implements [Value].
func (val ValueProgramName) Strings() []string {
	return []string{val.Name}
}

// Token implements [Value].
func (val ValueProgramName) Token() scanner.Token {
	return val.Tok
}

// ValueOption is a [Value] containing a parsed [*Option].
type ValueOption struct {
	// Option is the corresponding [*Option].
	Option *Option

	// Tok is the token from which we parsed this [*Option].
	Tok scanner.Token

	// Value is the possibly-empty value. Specifically:
	//
	//	1. For [OptionTypeEarlyArgumentNone] this field is empty.
	//
	//	2. For [OptionTypeStandaloneArgumentNone] this field is empty.
	//
	//	3. For [OptionTypeGroupedArgumentNone] this field is empty.
	//
	//	4. For [OptionTypeStandaloneArgumentRequired] this field
	// 	   contains the value of the parsed argument.
	//
	//	5. For [OptionTypeGroupedArgumentRequired] this field
	// 	   contains the value of the parsed argument.
	//
	//	6. For [OptionTypeStandaloneArgumentOptional] this field
	// 	   contains the value of the parsed argument, if any,
	// 	   or the default value specified in [*Option], otherwise.
	Value string
}

var _ Value = ValueOption{}

// Strings implements [Value].
func (val ValueOption) Strings() []string {
	var output []string
	switch val.Option.Type {
	case OptionTypeEarlyArgumentNone, OptionTypeGroupableArgumentNone, OptionTypeStandaloneArgumentNone:
		output = append(output, val.Option.Prefix+val.Option.Name)

	case OptionTypeStandaloneArgumentOptional:
		output = append(output, val.Option.Prefix+val.Option.Name+"="+val.Value)

	case OptionTypeStandaloneArgumentRequired, OptionTypeGroupableArgumentRequired:
		output = append(output, val.Option.Prefix+val.Option.Name)
		output = append(output, val.Value)

	default:
		panic(fmt.Sprintf("unhandled option type: %d", val.Option.Type))
	}
	return output
}

// Token implements [Value].
func (val ValueOption) Token() scanner.Token {
	return val.Tok
}

// ValuePositionalArgument is a [Value] containing a parsed positional argument.
type ValuePositionalArgument struct {
	// Tok is the token associated with the value.
	Tok scanner.Token

	// Value is the argument value.
	Value string
}

var _ Value = ValuePositionalArgument{}

// Strings implements [Value].
func (a ValuePositionalArgument) Strings() []string {
	return []string{a.Value}
}

// Token implements [Value].
func (val ValuePositionalArgument) Token() scanner.Token {
	return val.Tok
}

// ValueOptionsArgumentsSeparator is a [Value] containing a parsed separator.
type ValueOptionsArgumentsSeparator struct {
	// Separator is the separator value.
	Separator string

	// Tok is the token associated with the item.
	Tok scanner.Token
}

var _ Value = ValueOptionsArgumentsSeparator{}

// Strings implements [Value].
func (s ValueOptionsArgumentsSeparator) Strings() []string {
	return []string{s.Separator}
}

// Token implements [Value].
func (val ValueOptionsArgumentsSeparator) Token() scanner.Token {
	return val.Tok
}

func sortValues(input []Value) {
	slices.SortStableFunc(input, func(a, b Value) int {
		return a.Token().Index() - b.Token().Index()
	})
}
