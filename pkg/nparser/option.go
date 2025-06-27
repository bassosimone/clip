// option.go - option type.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

// Option specifies the kind of option to parse.
type Option struct {
	// DefaultValue is the default value assigned to the option
	// [Value] when the option argument is optional.
	DefaultValue string

	// Prefix is the prefix to use for parsing this option (e.g., `-`)
	Prefix string

	// Name is the option name without the prefix (e.g., `f`).
	Name string

	// Type is the option type.
	Type OptionType
}
