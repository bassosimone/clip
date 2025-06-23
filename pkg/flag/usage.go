// usage.go - Code to print the usage information for a flag set.
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bassosimone/clip/pkg/parser"
	"github.com/bassosimone/clip/pkg/textwrap"
)

// Options returns all the configured [*Option] in the [*FlagSet].
//
// You typically want to call this method to print the usage.
func (fx *FlagSet) Options() []*Option {
	// Collect options and deduplicate them
	uniq := make(map[*Option]struct{})
	for _, option := range fx.optionsLong {
		uniq[option] = struct{}{}
	}
	for _, option := range fx.optionsShort {
		uniq[option] = struct{}{}
	}

	// Transform the map to a slice
	options := make([]*Option, 0, len(uniq))
	for option := range uniq {
		options = append(options, option)
	}

	// Sort options by their shortest name (short or long)
	sort.SliceStable(options, func(i, j int) bool {
		// Get shortest name for sorting
		iName := options[i].LongName
		if options[i].ShortName != 0 {
			iName = string(options[i].ShortName)
		}
		jName := options[j].LongName
		if options[j].ShortName != 0 {
			jName = string(options[j].ShortName)
		}
		// Compare shortest names (ignoring prefix like - or --)
		return iName < jName
	})

	return options
}

// Usage returns a string containing the [*FlagSet] usage information.
func (fx *FlagSet) Usage() string {
	var sb strings.Builder

	// Print the usage string
	fmt.Fprintf(&sb, "%s", fx.UsageSynopsis())

	// Print the options
	fmt.Fprintf(&sb, "options:\n")
	fmt.Fprintf(&sb, "%s\n", fx.UsageOptions())

	return sb.String()
}

// UsageSynopsis returns a string containing the [*FlagSet] usage synopsis.
func (fx *FlagSet) UsageSynopsis() string {
	var sb strings.Builder

	// Gather the separator to use (pick the first one for simplicity)
	var sep string
	if len(fx.parser.Separators) > 0 {
		sep = " [" + fx.parser.Separators[0] + "] "
	}

	// Print the synopsis string
	fmt.Fprintf(&sb, "\nusage: %s [options]%s[arguments]\n\n", fx.progname, sep)
	return sb.String()
}

// UsageOptions formats the usage information for the options in the [*FlagSet].
func (fx *FlagSet) UsageOptions() string {
	var sb strings.Builder

	// Gather the short option prefix (pick the first one for simplicity)
	var spref string
	if len(fx.parser.ShortOptionPrefixes) > 0 {
		spref = fx.parser.ShortOptionPrefixes[0]
	}

	// Gather the long option prefix (pick the first one for simplicity)
	var lpref string
	if len(fx.parser.LongOptionPrefixes) > 0 {
		lpref = fx.parser.LongOptionPrefixes[0]
	}

	// Print the options
	for _, opt := range fx.Options() {
		// Determine whether the option has a value
		var value string
		if opt.Value.OptionType() != parser.OptionTypeBool {
			value = " VALUE"
		}

		// Customize formatting depending on how the option is defined
		switch {
		case opt.ShortName != 0 && opt.LongName != "":
			fmt.Fprintf(&sb, "  %s%s, %s%s%s\n", spref, string(opt.ShortName), lpref, opt.LongName, value)
			fmt.Fprintf(&sb, "%s\n\n", textwrap.Do(opt.Usage, 72, "    "))

		case opt.ShortName != 0:
			fmt.Fprintf(&sb, "  %s%s%s\n", spref, string(opt.ShortName), value)
			fmt.Fprintf(&sb, "%s\n\n", textwrap.Do(opt.Usage, 72, "    "))

		case opt.LongName != "":
			fmt.Fprintf(&sb, "  %s%s%s\n", lpref, opt.LongName, value)
			fmt.Fprintf(&sb, "%s\n\n", textwrap.Do(opt.Usage, 72, "    "))
		}
	}
	return sb.String()
}
