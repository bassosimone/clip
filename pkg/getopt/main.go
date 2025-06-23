// main.go - getopt(1) implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"errors"
	"strings"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/parser"
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
	px := &parser.Parser{
		Flags: parser.FlagNoPermute,
		LongOptions: map[string]parser.OptionType{
			"longoptions": parser.OptionTypeString,
			"options":     parser.OptionTypeString,
		},
		LongOptionPrefixes: []string{"--"},
		Separators:         []string{"--"},
		ShortOptions: map[string]parser.OptionType{
			"l": parser.OptionTypeString,
			"o": parser.OptionTypeString,
		},
		ShortOptionPrefixes: []string{"-"},
	}

	// Parse the command line arguments
	items, err := px.Parse(argv)
	if err != nil {
		return nil, err
	}

	// Initalize the options, longoptions, and parameters
	var (
		optstring  string
		options    []Option
		parameters []string
	)

	// Skip over the program name
	assert.True(len(items) >= 1, "program name not found")
	_, ok := items[0].(parser.ProgramNameItem)
	assert.True(ok, "program name not found")
	items = items[1:]

	// Process the option items
	for len(items) > 0 {
		opt, ok := items[0].(parser.OptionItem)
		if !ok {
			break
		}
		items = items[1:]
		switch opt.Name {
		case "o", "options":
			optstring = opt.Value
		case "l", "longoptions":
			options = append(options, parseLongOptions(opt)...)
		}
	}

	// Handle the case where we are out of items
	if len(items) <= 0 {
		return []string{}, nil
	}

	// Expect to see the separator
	if _, ok := items[0].(parser.SeparatorItem); !ok {
		return nil, ErrExpectedSeparator
	}
	items = items[1:]

	// Process the command line items
	for len(items) > 0 {
		item, ok := items[0].(parser.ArgumentItem)
		assert.True(ok, "expected argument item") // parser should guarantee this
		items = items[1:]
		parameters = append(parameters, item.Value)
	}

	// Reorder the arguments
	return reorder(parameters, optstring, options)
}

func parseLongOptions(opt parser.OptionItem) []Option {
	options := []Option{}
	values := strings.SplitSeq(opt.Value, ",")
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
	items, err := Long(argv, optstring, options)

	// handle any parsing errors
	if err != nil {
		return nil, err
	}

	// remove the dummy program name
	assert.True(len(items) >= 1, "no items found after parsing")
	_, ok := items[0].(parser.ProgramNameItem)
	assert.True(ok, "program name not found after parsing")
	items = items[1:]

	// serialize the command line items w/o the dummy program name
	return Serialize(items), nil
}
