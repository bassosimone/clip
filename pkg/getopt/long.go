// long.go - getopt_long implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"math"
	"os"
	"strings"

	"github.com/bassosimone/clip/pkg/nparser"
)

// Option describes a long option.
type Option struct {
	// Name is the option name without the `--` prefix.
	Name string

	// HasArg indicates whether the option takes an argument.
	HasArg bool

	// IsArgOptional flags optional arguments.
	//
	// Added in v0.6.0 - before that, option arguments
	// were always required.
	IsArgOptional bool

	// DefaultValue is the default value used
	// when the argument is optional.
	//
	// Added in v0.6.0 - before that, option arguments
	// were always required.
	DefaultValue string
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
func Long(argv []string, optstring string, options []Option) ([]nparser.Value, error) {
	// Honour the POSIXLY_CORRECT environment variable.
	var disablePermute bool
	if _, found := lookupEnv("POSIXLY_CORRECT"); found {
		disablePermute = true
	}

	// Check whether to disable permutation by checking the optstring
	if strings.HasPrefix(optstring, "-") {
		disablePermute = true
		optstring = optstring[1:]
	}

	// Instantiate the parser
	px := &nparser.Parser{
		DisablePermute:            disablePermute,
		MaxPositionalArguments:    math.MaxInt,
		MinPositionalArguments:    0,
		OptionsArgumentsSeparator: "--",
		Options:                   []*nparser.Option{},
	}

	// Register the long options
	for _, option := range options {
		px.Options = append(px.Options, &nparser.Option{
			DefaultValue: option.DefaultValue,
			Name:         option.Name,
			Prefix:       "--",
			Type: (func(o Option) nparser.OptionType {
				switch {
				case o.HasArg && o.IsArgOptional:
					return nparser.OptionTypeStandaloneArgumentOptional
				case o.HasArg:
					return nparser.OptionTypeStandaloneArgumentRequired
				default:
					return nparser.OptionTypeStandaloneArgumentNone
				}
			}(option)),
		})
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
		px.Options = append(px.Options, &nparser.Option{
			Name:   optname,
			Prefix: "-",
			Type: (func(hasArg bool) nparser.OptionType {
				switch hasArg {
				case true:
					return nparser.OptionTypeGroupableArgumentRequired
				default:
					return nparser.OptionTypeGroupableArgumentNone
				}
			}(hasArg)),
		})
	}

	// Parse the arguments vector
	return px.Parse(argv)
}
