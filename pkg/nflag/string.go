// string.go - String flag implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import (
	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nparser"
)

// StringFlag is like [*FlagSet.StringFlagVar] but returns a string variable
// rather than accepting the variable as its first argument.
func (fx *FlagSet) StringFlag(longName string, shortName byte, usage string) *string {
	var value string
	fx.StringFlagVar(&value, longName, shortName, usage)
	return &value
}

// StringFlagVar adds flags for setting the given string variable.
//
// The flag default value is set to the *valuep value.
//
// With longName="header" and shortName='h', the default configuration creates:
//
//  1. a `--header <value>` flag with mandatory argument.
//
//  2. a `-H <value>` flag with mandatory argument.
//
// As a side effect of seeing either flag, the pointee will be set to `<value>`.
//
// If longName and shortName are empty, this method will panic. If just one
// of them is empty, this method skips creating the related flag.
func (fx *FlagSet) StringFlagVar(valuep *string, longName string, shortName byte, usage string) {
	// make sure at least one of the two names is set
	assert.True(longName != "" || shortName != 0, "longName and shortName cannot be both zero values")

	// make sure the pointer is not nil
	assert.True(valuep != nil, "valuep cannot be nil")

	// be prepared for potentially adding two flags
	var long, short *Flag

	// possibly create the long flag value
	if longName != "" {
		long = &Flag{
			Modified: false,
			Option: &nparser.Option{
				Type:   nparser.OptionTypeStandaloneArgumentRequired,
				Prefix: fx.LongFlagPrefix,
				Name:   longName,
			},
			TakesArg: true,
			Value:    stringValue{valuep},
			Usage:    usage,
		}
	}

	// possibly create the short flag value
	if shortName != 0 {
		short = &Flag{
			Modified: false,
			Option: &nparser.Option{
				Type:   nparser.OptionTypeGroupableArgumentRequired,
				Prefix: fx.ShortFlagPrefix,
				Name:   string(shortName),
			},
			TakesArg: true,
			Value:    stringValue{valuep},
			Usage:    usage,
		}
	}

	// add as much as possible
	fx.mustAddLongAndShortFlag(long, short)
}

type stringValue struct {
	valuep *string
}

var _ Value = stringValue{}

func (v stringValue) Set(value string) error {
	*v.valuep = value
	return nil
}

func (v stringValue) String() string {
	return *v.valuep
}
