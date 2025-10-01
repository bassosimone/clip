// usage.go - Code to print the usage information.
// SPDX-License-Identifier: GPL-3.0-or-later

package nflag

import (
	"fmt"
	"io"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/textwrap"
)

// PrintUsage prints the usage string to the given [io.Writer].
//
// The template we use is roughly this:
//
//	usage: <programName> [options] [<separator>] <arguments>
//
//	<description>
//
//	Options:
//	  <shortPrefix><shortName>, <longPrefix><longName> <argument>
//	    <description>
//
//	<examples>
//
// We adapt it depending on the [*FlagSet] configuration. For example,
// we don't print the separator if none is defined.
//
// This method panics in case of I/O error.
func (fx *FlagSet) PrintUsage(w io.Writer) {

	// construct the synopsis line
	assert.NotError1(fmt.Fprintf(w, "Usage: %s", fx.ProgramName))
	if len(fx.usageView) > 0 {
		assert.NotError1(fmt.Fprint(w, " [options]"))
	}
	if minimum := fx.MinPositionalArgs; minimum >= 0 {
		if maximum := fx.MaxPositionalArgs; maximum >= minimum {
			usage := fx.PositionalArgumentsUsage
			assert.NotError1(fmt.Fprintf(w, " %s", usage))
		}
	}
	assert.NotError1(fmt.Fprint(w, "\n\n"))

	// optionally print the description
	if descr := fx.Description; descr != "" {
		descr = textwrap.Do(descr, 72, "")
		assert.NotError1(fmt.Fprintf(w, "%s\n\n", descr))
	}

	// optionally print the options
	if len(fx.usageView) > 0 {
		assert.NotError1(fmt.Fprintf(w, "Options:\n"))
		for _, pair := range fx.usageView {
			long, short := pair.LongFlag, pair.ShortFlag
			assert.NotError1(fmt.Fprint(w, "  "))
			if short != nil {
				assert.NotError1(fmt.Fprintf(w, "%s%s", short.Option.Prefix, short.Option.Name))
			}
			if short != nil && long != nil {
				assert.NotError1(fmt.Fprint(w, ", "))
			}
			if long != nil {
				assert.NotError1(fmt.Fprintf(w, "%s%s", long.Option.Prefix, long.Option.Name))
			}
			if pair.TakesArg {
				space := map[bool]string{true: "=", false: " "}
				assert.NotError1(fmt.Fprintf(w, "%sVALUE", space[long != nil]))
			}
			assert.NotError1(fmt.Fprintf(w, "\n"))
			usage := textwrap.Do(pair.Usage, 72, "    ")
			assert.NotError1(fmt.Fprint(w, usage))
			assert.NotError1(fmt.Fprintf(w, "\n\n"))
		}
	}

	// optionally print the examples
	if examples := fx.Examples; examples != "" {
		assert.NotError1(fmt.Fprint(w, examples))
	}
}

// PrintHelpHint prints the help hint to the given [io.Writer].
//
// The template is roughly:
//
//	Try '<programName> <flagPrefix><flagName>' for more help.\n
//
// This method panics in case of I/O error.
func (fx *FlagSet) PrintHelpHint(w io.Writer) {
	ishelp := func(fp *Flag) bool {
		_, ok := fp.Value.(*helpValue)
		return ok
	}

	for _, pair := range fx.usageView {
		long, short := pair.LongFlag, pair.ShortFlag
		assert.True(long != nil || short != nil, "PrintHelpHint found an invalid LongShortFlag")
		var x *Flag
		switch {
		case long != nil:
			x = long
		case short != nil:
			x = short
		}
		if ishelp(x) {
			assert.NotError1(fmt.Fprintf(w, "Try '%s %s%s' for more help.\n",
				fx.ProgramName, x.Option.Prefix, x.Option.Name))
			return
		}
	}
}
