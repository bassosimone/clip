// value.go - Value interface and implementation for command-line flag values
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

// Value represents a generic flag value.
type Value interface {
	// Set sets the value of the flag.
	//
	// This method MAY be called multiple times if
	// the command-line flag is repeated.
	Set(value string) error

	// String returns the string representation of the value.
	String() string
}
