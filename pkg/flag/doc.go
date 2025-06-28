// doc.go - package documentation.
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package flag provides facilities for command-line flag parsing with support for both
short and long options with an API similar to the standard library's flag package.

Deprecated: Use [github.com/bassosimone/clip/pkg/nflag] instead.

The [FlagSet] type represents a set of command-line flags and provides methods
to define boolean and string flags with various combinations of short (-f) and
long (--flag) option names. The package supports GNU-style flag parsing with
customizable option prefixes and separators through the underlying [parser.Parser].

The [NewFlagSet] function creates a new flag set with configurable error
handling behavior. Use [FlagSet.Bool], [FlagSet.String], and their variants
to define flags, then call [FlagSet.Parse] to parse command-line arguments.

The package provides comprehensive flag definition methods including [FlagSet.BoolLong],
[FlagSet.BoolShort], [FlagSet.StringLong], [FlagSet.StringShort], and their
corresponding Var variants that accept existing variable pointers.
*/
package flag
