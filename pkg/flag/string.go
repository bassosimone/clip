// string.go - String flag value implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import "github.com/bassosimone/clip/pkg/parser"

// StringLong creates a new string flag available as --name with the
// given usage string and default value. The actual flag prefixes may
// vary if you modify the parser by using [*FlagSet.Parser].
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) StringLong(name string, value string, usage string) *string {
	return fx.String(name, 0, value, usage)
}

// StringLongVar is like [*FlagSet.StringLong] but accepts a pointer to a string
// variable instead of returning a pointer to a string.
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) StringLongVar(p *string, name string, usage string) {
	fx.StringVar(p, name, 0, usage)
}

// StringShort creates a new string flag available as -name with the
// given usage string and an empty default value. The actual flag prefixes
// may vary if you modify the parser by using [*FlagSet.Parser].
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) StringShort(name byte, usage string) *string {
	return fx.String("", name, "", usage)
}

// StringShortVar is like [*FlagSet.StringShort] but accepts a pointer to a string
// variable instead of returning a pointer to a string.
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) StringShortVar(p *string, name byte, usage string) {
	fx.StringVar(p, "", name, usage)
}

// String creates a new string flag available as --longName and -shortName
// with the given usage string and default value. The actual flag prefixes
// may vary if you modify the parser by using [*FlagSet.Parser].
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) String(longName string, shortName byte, value string, usage string) *string {
	fx.StringVar(&value, longName, shortName, usage)
	return &value
}

// StringVar is like [*FlagSet.String] but accepts a pointer to a string
// variable instead of returning a pointer to a string.
//
// This method MUST be invoked before calling [FlagSet.Parse].
func (fx *FlagSet) StringVar(p *string, longName string, shortName byte, usage string) {
	option := &Option{
		LongName:  longName,
		ShortName: shortName,
		Usage:     usage,
		Value: &stringValue{
			value: p,
		},
	}
	fx.AddOption(option)
}

// stringValue implements [Value] for string.
type stringValue struct {
	value *string
}

var _ Value = stringValue{}

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
