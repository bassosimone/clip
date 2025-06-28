// main.go - getopt(1) implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"errors"
	"math"
	"strings"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nparser"
)

// ErrExpectedSeparator is returned when the separator is not found.
var ErrExpectedSeparator = errors.New("expected separator")

// Main implements the getopt(1) tool.
//
// The usage is as follows:
//
//	getopt [-o|--options optstring] [-l|--longoptions longopts] [--] [params]
//
// The optstring is like in getopt(3). The longopts is a list of comma separated
// names, followed by `:` if the option takes an argument. Multiple --longoptions
// flags can be specified and add to the already specified options.
//
// The provided params must not contain the program name, like in getopt(1).
func Main(argv []string) ([]string, error) {
	// Instantiate the parser
	px := &nparser.Parser{
		DisablePermute:            true,
		MaxPositionalArguments:    math.MaxInt,
		MinPositionalArguments:    0,
		OptionsArgumentsSeparator: "--",
		Options: []*nparser.Option{
			{
				Name:   "longoptions",
				Prefix: "--",
				Type:   nparser.OptionTypeStandaloneArgumentRequired,
			},
			{
				Name:   "options",
				Prefix: "--",
				Type:   nparser.OptionTypeStandaloneArgumentRequired,
			},
			{
				Name:   "l",
				Prefix: "-",
				Type:   nparser.OptionTypeGroupableArgumentRequired,
			},
			{
				Name:   "o",
				Prefix: "-",
				Type:   nparser.OptionTypeGroupableArgumentRequired,
			},
		},
	}

	// Parse the command line arguments
	values, err := px.Parse(argv)
	if err != nil {
		return nil, err
	}

	// Initialize the options, longoptions, and parameters
	var (
		optstring  string
		options    []Option
		parameters []string
	)

	// Skip over the program name.
	//
	// Note: using assert here because the parser guarantees that
	// the first value is always the program name.
	assert.True(len(values) >= 1, "program name not found")
	_, ok := values[0].(nparser.ValueProgramName)
	assert.True(ok, "program name not found")
	values = values[1:]

	// Process the option values
	for len(values) > 0 {
		opt, ok := values[0].(nparser.ValueOption)
		if !ok {
			break
		}
		values = values[1:]
		switch opt.Option.Name {
		case "o", "options":
			optstring = opt.Value
		case "l", "longoptions":
			options = append(options, parseLongOptions(opt.Value)...)
		}
	}

	// Handle the case where we are out of values
	if len(values) <= 0 {
		return []string{}, nil
	}

	// Expect to see the separator
	if _, ok := values[0].(nparser.ValueOptionsArgumentsSeparator); !ok {
		return nil, ErrExpectedSeparator
	}
	values = values[1:]

	// Process the command line values
	for len(values) > 0 {
		value, ok := values[0].(nparser.ValuePositionalArgument)
		assert.True(ok, "expected ValuePositionalArgument") // parser should guarantee this
		values = values[1:]
		parameters = append(parameters, value.Value)
	}

	// Reorder the arguments
	return reorder(parameters, optstring, options)
}

func parseLongOptions(optValue string) []Option {
	options := []Option{}
	// strings.Split and strings.SplitSet return a single
	// value when there are no commas
	values := strings.SplitSeq(optValue, ",")
	for value := range values {
		hasArg := false
		if strings.HasSuffix(value, ":") {
			value = strings.TrimSuffix(value, ":")
			hasArg = true
		}
		options = append(options, Option{
			Name:   value,
			HasArg: hasArg,
		})
	}
	return options
}

func reorder(parameters []string, optstring string, options []Option) ([]string, error) {
	// append a dummy program name to the command line
	argv := append([]string{"dummy"}, parameters...)

	// parse using the [Long] function
	values, err := Long(argv, optstring, options)

	// handle any parsing errors
	if err != nil {
		return nil, err
	}

	// remove the dummy program name
	//
	// Note: using assert here because the parser guarantees that
	// the first value is always the program name.
	//
	// Also the parser should return at least one value
	assert.True(len(values) >= 1, "no values found after parsing")
	_, ok := values[0].(nparser.ValueProgramName)
	assert.True(ok, "program name not found after parsing")
	values = values[1:]

	// serialize the command line values w/o the dummy program name
	return Serialize(values), nil
}
