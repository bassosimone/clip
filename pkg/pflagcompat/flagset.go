// flagset.go - definition of FlagSet.
// SPDX-License-Identifier: GPL-3.0-or-later

package pflagcompat

import "github.com/bassosimone/clip/pkg/nflag"

// FlagSet is a tiny wrapper around [*nflag.FlagSet].
type FlagSet struct {
	Set *nflag.FlagSet
}

// NewFlagSet constructs a [*FlagSet] following the GNU command
// line parsing conventions. You can still modify the settings by
// editing the underlying [*FlagSet] `Set` field.
func NewFlagSet(progname string, handling nflag.ErrorHandling) FlagSet {
	// Implementation note: the documentation of nflag.NewFlagSet explicitly
	// states that the constructor uses the GNU convention. Changing this would
	// be a huge breaking change. So it feels safe to rely on that specific
	// GNU semantics here without being too paranoid about it.
	return FlagSet{Set: nflag.NewFlagSet(progname, handling)}
}

// Parse invokes the underlying [*nflag.FlagSet.Parse] method.
func (fx *FlagSet) Parse(args []string) error {
	return fx.Set.Parse(args)
}
