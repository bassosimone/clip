// bool.go - Boolean flag value implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flagemu

import (
	"strconv"

	"github.com/bassosimone/clip/pkg/parser"
)

// boolValue implements [Value] for bool.
type boolValue struct {
	value *bool
}

var _ Value = boolValue{}

// newBoolValue creates a new BoolValue with the given default value.
func newBoolValue(value *bool) boolValue {
	return boolValue{value: value}
}

// OptionType implements [Value].
func (v boolValue) OptionType() parser.OptionType {
	return parser.OptionTypeBool
}

// String implements [Value].
func (v boolValue) String() string {
	return strconv.FormatBool(*v.value)
}

// Set implements [Value].
func (v boolValue) Set(value string) error {
	*v.value = value == "true"
	return nil
}
