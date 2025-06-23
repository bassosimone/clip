// option.go - Option interface and implementation for command-line flag values
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import "github.com/bassosimone/clip/pkg/parser"

// Option represents a command-line flag option.
type Option struct {
	// LongName is the long option name. Long options are parsed
	// independently and are never grouped. A long option can have
	// a single-character name. You can emulate the Go flag
	// parser by configuring the [*parser.Parser] to use `-` and
	// `--` as the long option prefix.
	//
	// Leave this field empty if the option does not have a long name.
	LongName string

	// Modified indicates that the option has been modified
	// through values provided using the command-line.
	Modified bool

	// ParamName is the optional parameter name associated with
	// the option. If not set, we use 'VALUE'.
	ParamName string

	// ShortName is the short option name. Short options are parsed
	// together and may be grouped like in GNU getopt.
	//
	// Leave this field empty if the option does not have a short name.
	ShortName byte

	// Usage is the usage string of the option.
	Usage string

	// Value is the value of the option.
	Value Value
}

// FormatParamName returns the string associated with the [Option] parameter
// name. The return value is an empty string for boolean options, which are
// not associated with a parameter name. Instead, for string options, the return
// value is the ParamName, if set, and otherwise 'VALUE'.
func (opt *Option) FormatParamName() string {
	switch {
	case opt.Value.OptionType() == parser.OptionTypeBool:
		return ""
	case opt.ParamName != "":
		return opt.ParamName
	default:
		return "VALUE"
	}
}
