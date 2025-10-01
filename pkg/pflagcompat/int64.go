// int64.go - int64 flag implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package pflagcompat

import "github.com/bassosimone/clip/pkg/assert"

// Int64 adds a long-only int64 flag with a default value and
// returns a pointer to the flag value.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.Int64].
func (fx *FlagSet) Int64(longName string, defvalue int64, usage string) *int64 {
	return fx.Int64P(longName, 0, defvalue, usage)
}

// Int64Var adds a long-only int64 flag with a default value and
// arranges for parsing to modify the given pointer.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.Int64Var].
func (fx *FlagSet) Int64Var(valuep *int64, longName string, defvalue int64, usage string) {
	fx.Int64VarP(valuep, longName, 0, defvalue, usage)
}

// Int64P adds a flag with both long and short name and a default value
// and returns to the caller a pointer to the flag value.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.Int64P].
func (fx *FlagSet) Int64P(longName string, shortName byte, defvalue int64, usage string) *int64 {
	value := defvalue
	fx.Set.Int64FlagVar(&value, longName, shortName, usage)
	return &value
}

// Int64VarP adds a flag with both long and short name and a default value
// and arranges for parsing to modify the given pointer.
//
// This method is equivalent to [*github.com/spf13/pflag.FlagSet.Int64VarP].
func (fx *FlagSet) Int64VarP(valuep *int64, longName string, shortName byte, defvalue int64, usage string) {
	assert.True(valuep != nil, "valuep cannot be nil")
	*valuep = defvalue
	fx.Set.Int64FlagVar(valuep, longName, shortName, usage)
}
