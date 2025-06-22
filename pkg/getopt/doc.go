// doc.go - package documentation.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package getopt provides facilities to implement the getopt(1) tool.
//
// The [Short] and [Long] functions emulate a subset of the traditional getopt
// implementation and of getopt_long by GNU. They check options for correctness
// and return a parsed list of [parser.CommandLineItem].
//
// The [Main] function implements the getopt(1) tool.
package getopt
