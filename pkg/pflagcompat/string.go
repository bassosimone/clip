// string.go - String flag implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package pflagcompat

import "github.com/bassosimone/clip/pkg/assert"

// String adds a long-only string flag with a default value and
// returns a pointer to the flag value.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.String].
func (fx *FlagSet) String(longName string, defvalue string, usage string) *string {
	return fx.StringP(longName, 0, defvalue, usage)
}

// StringVar adds a long-only string flag with a default value and
// arranges for parsing to modify the given pointer.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.StringVar].
func (fx *FlagSet) StringVar(valuep *string, longName string, defvalue string, usage string) {
	fx.StringVarP(valuep, longName, 0, defvalue, usage)
}

// StringP adds a flag with both long and short name and a default value
// and returns to the caller a pointer to the flag value.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.StringP].
func (fx *FlagSet) StringP(longName string, shortName byte, defvalue string, usage string) *string {
	value := defvalue
	fx.Set.StringFlagVar(&value, longName, shortName, usage)
	return &value
}

// StringVarP adds a flag with both long and short name and a default value
// and arranges for parsing to modify the given pointer.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.StringVarP].
func (fx *FlagSet) StringVarP(valuep *string, longName string, shortName byte, defvalue string, usage string) {
	assert.True(valuep != nil, "valuep cannot be nil")
	*valuep = defvalue
	fx.Set.StringFlagVar(valuep, longName, shortName, usage)
}
