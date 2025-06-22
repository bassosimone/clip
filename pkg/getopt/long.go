// long.go - getopt_long implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"os"
	"strings"

	"github.com/bassosimone/clip/pkg/parser"
)

// Option describes a long option.
type Option struct {
	// Name is the option name without the `--` prefix.
	Name string

	// HasArg indicates whether the option takes an argument.
	HasArg bool
}

var boolToOptionType = map[bool]parser.OptionType{
	true:  parser.OptionTypeString,
	false: parser.OptionTypeBool,
}

var lookupEnv = os.LookupEnv // for testing

// Long emulates a subset of the getopt_long implementation by GNU.
//
// The function returns the parsed command line arguments.
//
// We honour the POSIXLY_CORRECT environment variable. If the variable is
// set, we disable permutation of the command line arguments.
//
// If the optstring starts with `-`, we also disable permutation.
func Long(argv []string, optstring string, options []Option) ([]parser.CommandLineItem, error) {
	// Honour the POSIXLY_CORRECT environment variable.
	var flags parser.Flags
	if _, found := lookupEnv("POSIXLY_CORRECT"); found {
		flags |= parser.FlagNoPermute
	}

	// Check whether to disable permutation by checking the optstring
	if strings.HasPrefix(optstring, "-") {
		flags |= parser.FlagNoPermute
		optstring = optstring[1:]
	}

	// Instantiate the parser
	px := &parser.Parser{
		Flags:               flags,
		LongOptions:         map[string]parser.OptionType{},
		LongOptionPrefixes:  []string{"--"},
		Separators:          []string{"--"},
		ShortOptions:        map[string]parser.OptionType{},
		ShortOptionPrefixes: []string{"-"},
	}

	// Register the long options
	for _, option := range options {
		px.LongOptions[option.Name] = boolToOptionType[option.HasArg]
	}

	// Register the short options
	for len(optstring) > 0 {
		optname := string(optstring[0])
		optstring = optstring[1:]
		hasArg := false
		if strings.HasPrefix(optstring, ":") {
			optstring = optstring[1:]
			hasArg = true
		}
		px.ShortOptions[optname] = boolToOptionType[hasArg]
	}

	// Parse the arguments vector
	return px.Parse(argv)
}
