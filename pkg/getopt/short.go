// short.go - traditional getopt implementation.
// SPDX-License-Identifier: GPL-3.0-or-later

package getopt

import "github.com/bassosimone/clip/pkg/parser"

// Short emulates a subset of the traditional getopt implementation.
//
// This function is implemented in terms of the [Long] function.
func Short(argv []string, optstring string) ([]parser.CommandLineItem, error) {
	return Long(argv, optstring, nil)
}
