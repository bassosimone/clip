// bool.go - Boolean flag implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import (
	"strconv"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nparser"
)

// BoolFlag is like [*FlagSet.BoolFlagVar] but returns a bool variable
// rather than accepting the variable as its first argument.
func (fx *FlagSet) BoolFlag(longName string, shortName byte, usage string) *bool {
	var value bool
	fx.BoolFlagVar(&value, longName, shortName, usage)
	return &value
}

// BoolFlagVar adds flags for setting the given boolean variable.
//
// The flag default value is set to the *valuep value.
//
// Assuming longName="verbose" and shortName='v', the default configuration creates:
//
//  1. a `--verbose` long boolean flag.
//
//  2. a `-v` short boolean flag.
//
// As a side effect of seeing either flags, the pointee will be set to `true`.
//
// If longName and shortName are empty, this method will panic. If just one
// of them is empty, this method skips creating the related flag.
func (fx *FlagSet) BoolFlagVar(valuep *bool, longName string, shortName byte, usage string) {
	// make sure at least one of the two names is set
	assert.True(longName != "" || shortName != 0, "longName and shortName cannot be both zero values")

	// be prepared for potentially adding two flags
	var long, short *Flag

	// create a single underlying value for both flags
	mvalue := &boolValue{false, valuep}

	// possibly create the long flag value
	if longName != "" {
		long = &Flag{
			Option: &nparser.Option{
				Type:   nparser.OptionTypeStandaloneArgumentNone,
				Prefix: fx.LongFlagPrefix,
				Name:   longName,
			},
			TakesArg: false,
			Value:    mvalue,
			Usage:    usage,
		}
	}

	// possibly create the short flag value
	if shortName != 0 {
		short = &Flag{
			Option: &nparser.Option{
				Type:   nparser.OptionTypeGroupableArgumentNone,
				Prefix: fx.ShortFlagPrefix,
				Name:   string(shortName),
			},
			TakesArg: false,
			Value:    mvalue,
			Usage:    usage,
		}
	}

	// add as much as possible
	fx.mustAddLongAndShortFlag(long, short)
}

type boolValue struct {
	modified bool
	valuep   *bool
}

var _ Value = &boolValue{}

func (v *boolValue) Modified() bool {
	return v.modified
}

func (v *boolValue) Set(value string) error {
	*v.valuep = true
	v.modified = true
	return nil
}

func (v *boolValue) String() string {
	return strconv.FormatBool(*v.valuep)
}
