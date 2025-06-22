// value.go - Value interface and implementation for command-line flag values
// SPDX-License-Identifier: GPL-3.0-or-later

package flagemu

import "github.com/bassosimone/clip/pkg/parser"

// Value represents a flag value.
type Value interface {
	// OptionType returns the type of the option.
	OptionType() parser.OptionType

	// String returns the string representation of the value.
	String() string

	// Set sets the value of the flag.
	//
	// This method MAY be called multiple times.
	Set(value string) error
}
