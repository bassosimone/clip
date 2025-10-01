// help.go - Helpean flag implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import (
	"strconv"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nparser"
)

// AutoHelp adds a flag indicating that the user has requested for help.
//
// Assuming longName="help" and shortName='h', the default configuration creates:
//
//  1. a `--help` long boolean flag.
//
//  2. a `-h` short boolean flag.
//
// As a side effect of seeing either flag, the [*FlagSet] prints the usage message
// on the configured stdout and then implements this algorithm:
//
//  1. with [ContinueOnError], [ErrHelp] is returned.
//
//  2. with [ExitOnError], we call the configured exit function with status 2.
//
//  3. with [PanicOnError], we invoke panic.
//
// The help flag will be recognized and handled even when the command line is wrong
// and would not otherwise parse, this providing a nice UX.
//
// If longName and shortName are empty, this method will panic. If just one
// of them is empty, this method skips creating the related flag.
func (fx *FlagSet) AutoHelp(longName string, shortName byte, usage string) {
	// make sure at least one of the two names is set
	assert.True(longName != "" || shortName != 0, "longName and shortName cannot be both zero values")

	// be prepared for potentially adding two flags
	var (
		long  *Flag
		short *Flag
		value bool
	)

	// create a single underlying value for both flags
	mvalue := &helpValue{false, &value}

	// possibly create the long flag value
	if longName != "" {
		long = &Flag{
			Option: &nparser.Option{
				Type:   nparser.OptionTypeEarlyArgumentNone,
				Prefix: fx.LongFlagPrefix,
				Name:   longName,
			},
			TakesArg: false,
			Usage:    usage,
			Value:    mvalue,
		}
	}

	// possibly create the short flag value
	if shortName != 0 {
		short = &Flag{
			Option: &nparser.Option{
				Type:   nparser.OptionTypeEarlyArgumentNone,
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

type helpValue struct {
	modified bool
	valuep   *bool
}

var _ Value = &helpValue{}

func (v *helpValue) Modified() bool {
	return v.modified
}

func (v *helpValue) Set(value string) error {
	*v.valuep = true
	v.modified = true
	return nil
}

func (v *helpValue) String() string {
	return strconv.FormatBool(*v.valuep)
}
