// io.go - Code to manage I/O streams
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import "io"

// SetStderr sets the error stream for the flagset.
func (fx *FlagSet) SetStderr(w io.Writer) {
	fx.stderr = w
}

// SetStdout sets the output stream for the flagset.
func (fx *FlagSet) SetStdout(w io.Writer) {
	fx.stdout = w
}
