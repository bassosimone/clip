// positional.go - Positional arguments checks
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"errors"
	"fmt"
)

// NArg returns the number of positional arguments.
func (fx *FlagSet) NArg() int {
	return len(fx.args)
}

// ErrTooManyPositionalArgs is returned when the number of positional
// arguments is greater than the maximum allowed.
var ErrTooManyPositionalArgs = errors.New("too many positional arguments")

// ErrTooFewPositionalArgs is returned when the number of positional
// arguments is less than the minimum allowed.
var ErrTooFewPositionalArgs = errors.New("too few positional arguments")

// ErrUnexpectedNumberOfPositionalArgs is returned when the number of
// positional arguments is not the expected value.
var ErrUnexpectedNumberOfPositionalArgs = errors.New("unexpected number of positional arguments")

// PositionalArgsRangeCheck checks whether the number of positional
// arguments is within the given closed interval.
//
// Note: this function honors the [ErrorHandling] settings:
//
//  1. Returns the error with [ContinueOnError].
//
//  2. Invokes exit with [ExitOnError].
//
//  3. And calls panic with [PanicOnError].
func (fx *FlagSet) PositionalArgsRangeCheck(minargs, maxargs int) error {
	return fx.maybeHandleError(fx.positionalArgsRangeCheck(minargs, maxargs))
}

func (fx *FlagSet) positionalArgsRangeCheck(minargs, maxargs int) error {
	if minargs == maxargs && fx.NArg() != minargs {
		return fmt.Errorf("%w: expected %d, got %d", ErrUnexpectedNumberOfPositionalArgs, minargs, fx.NArg())
	}
	if fx.NArg() < minargs {
		return fmt.Errorf("%w: expected at least %d, got %d", ErrTooFewPositionalArgs, minargs, fx.NArg())
	}
	if fx.NArg() > maxargs {
		return fmt.Errorf("%w: expected at most %d, got %d", ErrTooManyPositionalArgs, maxargs, fx.NArg())
	}
	return nil
}

// PositionalArgsEqualCheck checks whether the number of positional
// arguments is equal to the given value.
func (fx *FlagSet) PositionalArgsEqualCheck(nargs int) error {
	return fx.PositionalArgsRangeCheck(nargs, nargs)
}
