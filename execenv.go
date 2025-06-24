// execenv.go - execution environment.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"io"
	"os"
)

// ExecEnv is the execution environment used by [Command].
type ExecEnv interface {
	// Args returns the system arguments.
	Args() []string

	// Exit terminates the program.
	Exit(exitcode int)

	// LookupEnv returns the value of the environment variable named by the key.
	LookupEnv(key string) (string, bool)

	// SignalNotify registers the specified signals to the channel.
	SignalNotify(c chan<- os.Signal, sig ...os.Signal)

	// Stdin is the standard input of the command.
	Stdin() io.Reader

	// Stdout is the standard output of the command.
	Stdout() io.Writer

	// Stderr is the standard error of the command.
	Stderr() io.Writer
}
