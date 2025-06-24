// root.go - root command.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"

	"github.com/bassosimone/clip/pkg/assert"
)

// RootCommand is the root [Command] of the application.
//
// The zero value is not ready to use. Initialize the mandatory fields.
type RootCommand[T ExecEnv] struct {
	// --- mandatory fields ---

	// Command is the mandatory command to execute. The code panics if
	// the Command to execute is a nil pointer.
	Command Command[T]

	// --- optional fields ---

	// AutoCancel optionally cancels the command if the user interrupts
	// its execution using signals (e.g., SIGINT, SIGTERM).
	AutoCancel bool
}

// Main is the root command entry point.
func (rx *RootCommand[T]) Main(env T) {
	// start with creating a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// possibly allow interrupting the command
	if rx.AutoCancel {
		sch := make(chan Signal, 1)
		env.SignalNotify(sch, interruptSignals...)
		go func() {
			defer cancel()
			<-sch
		}()
	}

	// create the command arguments
	argv := env.Args()
	assert.True(len(argv) >= 1, "the program name is required")
	args := &CommandArgs[T]{
		Env:         env,
		Args:        argv[1:],
		Command:     rx.Command,
		CommandName: argv[0],
		Parent:      nil,
	}

	// run the command
	assert.True(rx.Command != nil, "the command to execute is required")
	Must(env, rx.Command.Run(ctx, args))
}
