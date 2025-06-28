// serialize.go - CommandLineItem serialization.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import "github.com/bassosimone/clip/pkg/nparser"

// Serialize takes in input the [nparser.Value] parsed by [Long]
// or [Short] and serializes it to an argv such that:
//
//  1. the first value is the program name
//
//  2. options follow the program name
//
//  3. a separator is provided if it was originally provided
//
//  4. positional arguments follow the separator
func Serialize(input []nparser.Value) []string {
	// Create an empty slice for the results
	output := []string{}

	// Assume that the parser was correct and just append to
	// the output slice trusting parser serialization
	for _, value := range input {
		output = append(output, value.Strings()...)
	}

	// Return rewritten argv
	return output
}
