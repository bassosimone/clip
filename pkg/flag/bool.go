// bool.go - Boolean flag value implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"strconv"

	"github.com/bassosimone/clip/pkg/parser"
)

// BoolLong creates a new boolean flag available as --name with the
// given usage string and a false default value. The actual flag prefixes
// may vary if you modify the parser by using [*FlagSet.Parser].
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) BoolLong(name string, usage string) *bool {
	return fx.Bool(name, 0, usage)
}

// BoolLongVar is like [*FlagSet.BoolLong] but accepts a pointer to a bool
// variable instead of returning a pointer to a bool.
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) BoolLongVar(p *bool, name string, usage string) {
	fx.BoolVar(p, name, 0, usage)
}

// BoolShort creates a new boolean flag available as -name with the
// given usage string and a false default value. The actual flag prefixes
// may vary if you modify the parser by using [*FlagSet.Parser].
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) BoolShort(name byte, usage string) *bool {
	return fx.Bool("", name, usage)
}

// BoolShortVar is like [*FlagSet.BoolShort] but accepts a pointer to a bool
// variable instead of returning a pointer to a bool.
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) BoolShortVar(p *bool, name byte, usage string) {
	fx.BoolVar(p, "", name, usage)
}

// Bool creates a new boolean flag available as --longName and -shortName
// with the given usage string and false default value. The actual flag prefixes
// may vary if you modify the parser by using [*FlagSet.Parser].
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) Bool(longName string, shortName byte, usage string) *bool {
	var value bool
	fx.BoolVar(&value, longName, shortName, usage)
	return &value
}

// BoolVar is like [*FlagSet.Bool] but accepts a pointer to a bool
// variable instead of returning a pointer to a bool.
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) BoolVar(p *bool, longName string, shortName byte, usage string) {
	option := &Option{
		LongName:  longName,
		ShortName: shortName,
		Usage:     usage,
		Value: &boolValue{
			value: p,
		},
	}
	fx.AddOption(option)
}

// boolValue implements [Value] for bool.
type boolValue struct {
	value *bool
}

var _ Value = boolValue{}

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
