// string.go - String flag value implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flagemu

import "github.com/bassosimone/clip/pkg/parser"

// stringValue implements [Value] for string.
type stringValue struct {
	value *string
}

var _ Value = stringValue{}

// newStringValue creates a new StringValue with the given default value.
func newStringValue(value *string) stringValue {
	return stringValue{value: value}
}

// OptionType implements [Value].
func (v stringValue) OptionType() parser.OptionType {
	return parser.OptionTypeString
}

// String implements [Value].
func (v stringValue) String() string {
	return *v.value
}

// Set implements [Value].
func (v stringValue) Set(value string) error {
	*v.value = value
	return nil
}
