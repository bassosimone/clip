//go:build windows

// signals_windows.go - Windows-specific signals code.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import "syscall"

var interruptSignals = []Signal{syscall.SIGINT}
