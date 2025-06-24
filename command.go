// command.go - generic command interface
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import "context"

// CommandArgs contains the arguments passed to a [Command].
type CommandArgs[T ExecEnv] struct {
	// Args contains the arguments. This slice does not
	// include the command name.
	Args []string

	// Command is the commant itself. This field is useful when
	// the command main function is a standalone func.
	Command Command[T]

	// CommandName is the name of the command.
	CommandName string

	// Env is the execution environment.
	Env T

	// Parent is the possibly `nil` parent command.
	Parent Command[T]
}

// Command is the generic command interface.
type Command[T ExecEnv] interface {
	// BriefDescription returns a brief description of the command.
	BriefDescription() string

	// HelpFlag returns the help used by the the command.
	HelpFlag() string

	// Run is the command main entry point. The args contains the arguments
	// passed to the subcommand and does not include the command name.
	Run(ctx context.Context, args *CommandArgs[T]) error

	// SupportsSubcommands returns true if the command supports subcommands.
	SupportsSubcommands() bool
}
