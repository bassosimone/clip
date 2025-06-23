// assert.go - Utilities to write runtime assertions.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package assert provides utilities to write runtime assertions.
package assert

import "errors"

// True panics with the given message if the condition is false.
func True(condition bool, message string) {
	if !condition {
		panic(errors.New(message))
	}
}
