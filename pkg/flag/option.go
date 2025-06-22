// option.go - Option interface and implementation for command-line flag values
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

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
