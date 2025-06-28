// flagset.go - Code to parse command line flags.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import "github.com/bassosimone/clip/pkg/flag"

// ErrorHandling is an alias for [flag.ErrorHandling].
//
// Deprecated: use pkg/nflag or pkg/flag directly instead.
type ErrorHandling = flag.ErrorHandling

const (
	// ContinueOnError is an alias for [flag.ContinueOnError].
	//
	// Deprecated: use pkg/nflag or pkg/flag directly instead.
	ContinueOnError = flag.ContinueOnError

	// ExitOnError is an alias for [flag.ExitOnError].
	//
	// Deprecated: use pkg/nflag or pkg/flag directly instead.
	ExitOnError = flag.ExitOnError

	// PanicOnError is an alias for [flag.PanicOnError].
	//
	// Deprecated: use pkg/nflag or pkg/flag directly instead.
	PanicOnError = flag.PanicOnError
)

// NewFlagSet is an alias for [flag.NewFlagSet].
//
// Deprecated: use pkg/nflag or pkg/flag directly instead.
var NewFlagSet = flag.NewFlagSet

// FlagSet is an alias for [flag.FlagSet].
//
// Deprecated: use pkg/nflag or pkg/flag directly instead.
type FlagSet = flag.FlagSet
