// usage.go - Code to print the usage information for a flag set.
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"fmt"
	"sort"
	"strings"

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

// SetDescription sets the description of the [*FlagSet].
func (fx *FlagSet) SetDescription(description string) {
	fx.description = description
}

// Description returns the description of the [*FlagSet].
func (fx *FlagSet) Description() string {
	return fx.description
}

// ArgsDocs returns the documentation for the arguments.
func (fx *FlagSet) ArgsDocs() string {
	return fx.argsdocs
}

// SetArgsDocs sets the documentation for the arguments.
func (fx *FlagSet) SetArgsDocs(docs string) {
	fx.argsdocs = docs
}

// ProgramName returns the program name configured in the [*FlagSet].
func (fx *FlagSet) ProgramName() string {
	return fx.progname
}

// SetCustomUsage sets the custom usage function.
func (fx *FlagSet) SetCustomUsage(customusage func(fx *FlagSet) string) {
	fx.customusage = customusage
}

// Usage returns a string containing the [*FlagSet] usage information.
func (fx *FlagSet) Usage() string {
	// Honor the custom usage function if it's set
	if fx.customusage != nil {
		return fx.customusage(fx)
	}

	// Print the usage string
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s", fx.UsageSynopsis())

	// If there's a description, print it
	if desc := fx.Description(); desc != "" {
		fmt.Fprintf(&sb, "%s\n\n", desc)
	}

	// Print the options
	if optsusage := fx.UsageOptions(); optsusage != "" {
		fmt.Fprintf(&sb, "Options:\n")
		fmt.Fprintf(&sb, "%s", optsusage)
	}

	// Remind the user how to get help
	if lpref := fx.firstLongOptionsPrefix(); lpref != "" {
		fmt.Fprintf(&sb, "Use '%s %shelp' to show this help screen.", fx.ProgramName(), lpref)
	}
	return strings.TrimSpace(sb.String())
}

func (fx *FlagSet) firstShortOptionsPrefix() string {
	if len(fx.parser.ShortOptionPrefixes) > 0 {
		return fx.parser.ShortOptionPrefixes[0]
	}
	return ""
}

func (fx *FlagSet) firstLongOptionsPrefix() string {
	if len(fx.parser.LongOptionPrefixes) > 0 {
		return fx.parser.LongOptionPrefixes[0]
	}
	return ""
}

func (fx *FlagSet) formatFirstSeparator() string {
	if len(fx.parser.Separators) > 0 && fx.argsdocs != "" {
		return " [" + fx.parser.Separators[0] + "] "
	}
	return ""
}

// UsageSynopsis returns a string containing the [*FlagSet] usage synopsis.
func (fx *FlagSet) UsageSynopsis() string {
	var sb strings.Builder

	// Gather the separator to use (pick the first one for simplicity)
	sep := fx.formatFirstSeparator()

	// Gather the arguments documentation
	argsdoc := fx.ArgsDocs()

	// Print the synopsis string
	fmt.Fprintf(&sb, "Usage: %s [options]%s%s\n\n", fx.ProgramName(), sep, argsdoc)
	return sb.String()
}

// UsageOptions formats the usage information for the options in the [*FlagSet].
func (fx *FlagSet) UsageOptions() string {
	var sb strings.Builder

	// Gather the short option prefix (pick the first one for simplicity)
	spref := fx.firstShortOptionsPrefix()

	// Gather the long option prefix (pick the first one for simplicity)
	lpref := fx.firstLongOptionsPrefix()

	// Print the options
	for _, opt := range fx.Options() {
		// Customize formatting depending on how the option is defined
		value := opt.FormatParamName()
		if value != "" {
			value = " " + value
		}
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
