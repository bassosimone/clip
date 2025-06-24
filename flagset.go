// flagset.go - Code to parse command line flags.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import "github.com/bassosimone/clip/pkg/flag"

// ErrorHandling is an alias for [flag.ErrorHandling].
type ErrorHandling = flag.ErrorHandling

const (
	// ContinueOnError is an alias for [flag.ContinueOnError].
	ContinueOnError = flag.ContinueOnError

	// ExitOnError is an alias for [flag.ExitOnError].
	ExitOnError = flag.ExitOnError

	// PanicOnError is an alias for [flag.PanicOnError].
	PanicOnError = flag.PanicOnError
)

// NewFlagSet is an alias for [flag.NewFlagSet].
var NewFlagSet = flag.NewFlagSet

// FlagSet is an alias for [flag.FlagSet].
type FlagSet = flag.FlagSet
