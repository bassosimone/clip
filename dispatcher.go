// dispatcher.go - subcommand dispatcher.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/textwrap"
)

// DispatcherCommand is a [Command] dispatching to subcommands.
//
// The zero value is not ready to use. Initialize the mandatory fields.
type DispatcherCommand[T ExecEnv] struct {
	// --- mandatory fields ---

	// BriefDescriptionText is the mandatory short description of the dispatcher.
	BriefDescriptionText string

	// --- optional fields ---

	// Commands optionally contains the subcommands.
	Commands map[string]Command[T]

	// Usage is the optional usage string for this dispatcher. If empty, we
	// automatically generate a usage string when needed.
	Usage string

	// ErrorHandling indicates how the dispatcher should handle errors.
	//
	// Added in v0.3.0. When empty the behavior is [ContinueOnError], which
	// is exactly cosistent with the v0.2.0 behavior.
	ErrorHandling ErrorHandling
}

var _ Command[*StdlibExecEnv] = (*DispatcherCommand[*StdlibExecEnv])(nil)

// --- public code ---

// BriefDescription implements [Command].
func (dx *DispatcherCommand[T]) BriefDescription() string {
	return dx.BriefDescriptionText
}

// HelpFlag implements [Command].
func (dx *DispatcherCommand[T]) HelpFlag() string {
	return "--help"
}

// SupportsSubcommands implements [Command].
func (dx *DispatcherCommand[T]) SupportsSubcommands() bool {
	return true
}

// Run implements [Command].
func (dx *DispatcherCommand[T]) Run(ctx context.Context, args *CommandArgs[T]) error {
	return dx.maybeHandleError(args.Env, dx.dispatch(ctx, args))
}

// --- error filtering code ---

func (dx *DispatcherCommand[T]) maybeHandleError(env T, err error) error {
	// Determine what to do based on the policy
	switch dx.ErrorHandling {
	case ContinueOnError:
		return err

	case ExitOnError:
		switch {
		case errors.Is(err, ErrNoSuchCommand):
			env.Exit(2)
		case errors.Is(err, ErrAmbiguousCommandLine):
			env.Exit(2)
		default:
			env.Exit(1)
		}
	}

	// We end up here for [PanicOnError] or whenever env.Exit is so
	// broken that it doesn't actually exit.
	panic(err)
}

// --- dispatching code ---

// ErrAmbiguousCommandLine is returned when the dispatcher encounters an ambiguous command line.
var ErrAmbiguousCommandLine = errors.New("ambiguous command line")

func (dx *DispatcherCommand[T]) dispatch(ctx context.Context, args *CommandArgs[T]) error {
	// Handle the case where there are no arguments
	if len(args.Args) <= 0 {
		return dx.printUsage(args.Env, args.CommandName)
	}

	// Obtain subcommand name and arguments
	subName, subArgs := args.Args[0], args.Args[1:]

	// Attempt an exact match with subName
	if cmd := dx.Commands[subName]; cmd != nil {
		return dx.run(ctx, cmd, args, subName, subArgs)
	}

	// Handle special cases pertaining to obtain --help, -h, and help
	switch subName {
	case "--help", "-h":
		return dx.printUsage(args.Env, args.CommandName)

	case "help":
		return dx.maybeForwardHelp(ctx, args, subArgs)
	}

	// Find all matching subcommands inside the subArgs
	var indexes []int
	for idx, arg := range subArgs {
		if _, found := dx.Commands[arg]; found {
			indexes = append(indexes, idx)
		}
	}

	// If there is a single match, follow the rule of repair and
	// silently reorder the command line to execute it.
	//
	// Here's an example to fully understand what is going on:
	//
	// 	subName = "IN"
	// 	subArgs = ["A", "+short", "dig", "example.com"]
	//
	// In such a case we want to modify the overall state to:
	//
	// 	subName = "dig"
	// 	subArgs = ["IN", "A", "+short", "example.com"]
	if len(indexes) == 1 {
		rewritten := []string{subName}
		subName = subArgs[indexes[0]]
		for idx := 0; idx < len(subArgs); idx++ {
			if idx != indexes[0] {
				rewritten = append(rewritten, subArgs[idx])
			}
		}
		cmd := dx.Commands[subName]
		assert.True(cmd != nil, "expected command to be not nil here")
		subArgs = rewritten
		return dx.run(ctx, cmd, args, subName, subArgs)
	}

	// Let the user know that the command line is ambiguous
	if len(indexes) > 1 {
		fmt.Fprintln(args.Env.Stderr(), dx.formatAmbiguousCommandLine(args.CommandName, subArgs, indexes))
		return ErrAmbiguousCommandLine
	}

	// Otherwise mention that the given command was not found
	return dx.errorNoSuchCommand(args.Env, args.CommandName, subName)
}

func (dx *DispatcherCommand[T]) printUsage(env T, commandName string) error {
	_, err := fmt.Fprintln(env.Stdout(), dx.formatUsage(commandName))
	return err
}

func (dx *DispatcherCommand[T]) maybeForwardHelp(
	ctx context.Context, args *CommandArgs[T], subArgs []string) error {
	// We enter into this function with the following state:
	//
	//	subName = "help"
	//	subArgs = ???
	//
	// So the first action is to check whether there's anything
	// after `help` otherwise it's equivalent to `--help`
	if len(subArgs) <= 0 {
		return dx.printUsage(args.Env, args.CommandName)
	}
	subName, subArgs := subArgs[0], subArgs[1:]

	// Attempt to locate the command
	cmd := dx.Commands[subName]

	switch {
	// We don't have a subcommand with the provided name
	case cmd == nil:
		return dx.errorNoSuchCommand(args.Env, args.CommandName, subName)

	// The subcommand supports subcommands so we can ask for it to provide help
	case cmd.SupportsSubcommands():
		subArgs = append([]string{"help"}, subArgs...)
		return dx.run(ctx, cmd, args, subName, subArgs)

	// Otherwise just introduce an `--help` flag equivalent
	default:
		subArgs = append(subArgs, cmd.HelpFlag())
		return dx.run(ctx, cmd, args, subName, subArgs)
	}
}

func (dx *DispatcherCommand[T]) run(
	ctx context.Context, cmd Command[T], args *CommandArgs[T], subName string, subArgs []string) error {
	nargs := &CommandArgs[T]{
		Args:        subArgs,
		Command:     cmd,
		CommandName: args.CommandName + " " + subName,
		Env:         args.Env,
		Parent:      dx,
	}
	return cmd.Run(ctx, nargs)
}

// ErrNoSuchCommand is returned when a command is not found.
var ErrNoSuchCommand = errors.New("no such command")

func (dx *DispatcherCommand[T]) errorNoSuchCommand(env T, commandName, subcommandName string) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s: no such command: %s\n", commandName, subcommandName)
	fmt.Fprintf(&sb, "Try '%s --help' for more information.\n", commandName)
	fmt.Fprintln(env.Stderr(), sb.String())
	return ErrNoSuchCommand
}

// --- formatting code ---

func (dx *DispatcherCommand[T]) formatAmbiguousCommandLine(commandName string, subArgs []string, indexes []int) string {
	// List the ambiguous subcommands
	var sb strings.Builder
	fmt.Fprintf(
		&sb,
		"%s: fatal: ambiguous command line: found multiple subcommands (%s)",
		commandName,
		strings.Join(func() []string {
			out := make([]string, len(indexes))
			for i, idx := range indexes {
				out[i] = subArgs[idx]
			}
			return out
		}(), ", "))

	// Mention how to obtain further help
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "Try '%s help' for more information on '%s'.\n", commandName, commandName)

	// Return a trimmed string to avoid messing up with newlines
	return strings.TrimSpace(sb.String())
}

func (dx *DispatcherCommand[T]) formatUsage(commandName string) string {
	// If the user configured the usage string, use it
	if dx.Usage != "" {
		return dx.Usage
	}

	// Otherwise, create a simple usage string message
	var sb strings.Builder
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "Usage: %s [command] [args]\n", commandName)
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "%s\n", textwrap.Do(dx.BriefDescriptionText, 72, ""))
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "Commands:\n")
	for _, name := range dx.sortedSubcommandNames() {
		fmt.Fprintf(&sb, "  %s\n", name)
		cmd := dx.Commands[name]
		fmt.Fprintf(&sb, "%s\n\n", textwrap.Do(cmd.BriefDescription(), 72, "    "))
	}
	fmt.Fprintf(&sb, "Try '%s help COMMAND' for more information on COMMAND.\n", commandName)
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "Use '%s help' to show this help screen.\n", commandName)
	return strings.TrimSpace(sb.String())
}

func (dx *DispatcherCommand[T]) sortedSubcommandNames() []string {
	names := make([]string, 0, len(dx.Commands))
	for name := range dx.Commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
