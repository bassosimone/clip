// int64.go - int64 flag implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import (
	"strconv"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nparser"
)

// Int64Flag is like [*FlagSet.Int64FlagVar] but returns an int64 variable
// rather than accepting the variable as its first argument.
func (fx *FlagSet) Int64Flag(longName string, shortName byte, usage string) *int64 {
	var value int64
	fx.Int64FlagVar(&value, longName, shortName, usage)
	return &value
}

// Int64FlagVar adds flags for setting the given int64 variable.
//
// The flag default value is set to the *valuep value.
//
// Assuming longName="count" and shortName='c', the default configuration creates:
//
//  1. a `--count` long int64 flag.
//
//  2. a `-c` short int64 flag.
//
// As a side effect of seeing either flags, the pointee will be set to to the
// given value, if possible, otherwise an error is returned.
//
// If longName and shortName are empty, this method will panic. If just one
// of them is empty, this method skips creating the related flag.
func (fx *FlagSet) Int64FlagVar(valuep *int64, longName string, shortName byte, usage string) {
	// make sure at least one of the two names is set
	assert.True(longName != "" || shortName != 0, "longName and shortName cannot be both zero values")

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
			Value:    int64Value{valuep},
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
			Value:    int64Value{valuep},
			Usage:    usage,
		}
	}

	// add as much as possible
	fx.mustAddLongAndShortFlag(long, short)
}

type int64Value struct {
	valuep *int64
}

var _ Value = int64Value{}

func (v int64Value) Set(value string) error {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*v.valuep = parsed
	return nil
}

func (v int64Value) String() string {
	return strconv.FormatInt(*v.valuep, 10)
}
