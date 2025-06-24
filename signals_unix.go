//go:build unix

// signals_unix.go - Unix-specific signals code.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import "syscall"

var interruptSignals = []Signal{syscall.SIGINT, syscall.SIGTERM}
