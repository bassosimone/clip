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

// True1 is like [True] but returns the given [T] on success.
func True1[T any](value T, condition bool) T {
	True(condition, "assertion failed")
	return value
}

// NotError panics if the given error is not nil.
func NotError(err error) {
	if err != nil {
		panic(err)
	}
}

// NotError1 is like [NotError] but returns the given [T] on success.
func NotError1[T any](value T, err error) T {
	NotError(err)
	return value
}
