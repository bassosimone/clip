// serialize.go - CommandLineItem serialization.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import (
	"github.com/bassosimone/clip/pkg/parser"
)

// Serialize takes in input the [parser.CommandLineItem] parsed by [Long]
// or [Short] and serializes it to an argv such that:
//
//  1. the first item is the program name
//
//  2. options follow the program name
//
//  3. a separator is provided if it was originally provided
//
//  4. positional arguments follow the separator
func Serialize(items []parser.CommandLineItem) []string {
	// Create an empty slice for the results
	output := []string{}

	// Assume that the parser was correct and just append to
	// the output slice trusting parser serialization
	for _, item := range items {
		output = append(output, item.Strings()...)
	}

	// Return rewritten argv
	return output
}
