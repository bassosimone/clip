// bool.go - Boolean flag implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package pflagcompat

import "github.com/bassosimone/clip/pkg/assert"

// Bool adds a long-only boolean flag with a default value and
// returns a pointer to the flag value.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.Bool].
func (fx *FlagSet) Bool(longName string, defvalue bool, usage string) *bool {
	return fx.BoolP(longName, 0, defvalue, usage)
}

// BoolVar adds a long-only boolean flag with a default value and
// arranges for parsing to modify the given pointer.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.BoolVar].
func (fx *FlagSet) BoolVar(valuep *bool, longName string, defvalue bool, usage string) {
	fx.BoolVarP(valuep, longName, 0, defvalue, usage)
}

// BoolP adds a flag with both long and short name and a default value
// and returns to the caller a pointer to the flag value.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.BoolP].
func (fx *FlagSet) BoolP(longName string, shortName byte, defvalue bool, usage string) *bool {
	value := defvalue
	fx.Set.BoolFlagVar(&value, longName, shortName, usage)
	return &value
}

// BoolVarP adds a flag with both long and short name and a default value
// and arranges for parsing to modify the given pointer.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.BoolVarP].
func (fx *FlagSet) BoolVarP(valuep *bool, longName string, shortName byte, defvalue bool, usage string) {
	assert.True(valuep != nil, "valuep cannot be nil")
	*valuep = defvalue
	fx.Set.BoolFlagVar(valuep, longName, shortName, usage)
}
