// dispatcher.go - subcommand dispatcher.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/bassosimone/clip/pkg/assert"
	"github.com/bassosimone/clip/pkg/nflag"
	"github.com/bassosimone/clip/pkg/scanner"
	"github.com/bassosimone/clip/pkg/textwrap"
	"github.com/kballard/go-shellquote"
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
	ErrorHandling nflag.ErrorHandling

	// Version is the program version string. If this field is set, the
	// dispatcher will implement the following algorithm:
	//
	// If `--version` is the first argument, it will behave like the
	// `version` command had been specified instead.
	//
	// If the `version` command exists, the dispatcher will invoke it.
	//
	// Otherwise, the dispatcher will create a [*VersionCommand] on
	// the fly, configured with the version, and invoke it.
	//
	// Added in v0.4.0. When empty, we don't handle `--version` or `version.
	Version string

	// OptionPrefixes contains the option prefixes used when trying
	// to find a matching subcommand when the command line is not
	// correctly ordered with the subcommand being the first entry.
	//
	// If empty, we do not attempt to reorder the command line.
	//
	// New in v0.6.0. Before, we always tried to reorder, but that
	// was flawed since we did not know the prefixes.
	OptionPrefixes []string

	// OptionsArgumentsSeparator optionally specifies the separator
	// used to separate options from positional arguments.
	//
	// When empty, there is no separator.
	//
	// New in v0.6.0 and tied to OptionPrefixes.
	OptionsArgumentsSeparator string
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
	switch {
	case err == nil:
		return nil

	case dx.ErrorHandling == nflag.ContinueOnError:
		return err

	case dx.ErrorHandling == nflag.ExitOnError:
		switch {
		case errors.Is(err, ErrNoSuchCommand):
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

func (dx *DispatcherCommand[T]) dispatch(ctx context.Context, args *CommandArgs[T]) error {
	// Handle the case where there are no arguments
	if len(args.Args) <= 0 {
		return dx.printUsage(args.Env, args.CommandName)
	}

	// Unconditionally scan the command line. Note that, with empty separators
	// and prefixes, the scanner will return all positional arguments.
	//
	// Also, the scanner only fails if the program name is missing, which is
	// something we already checked for above, hence the assert.
	//
	// Also, the scanner returns the program name as the first token, we
	// know that, so we can assert again for the type and then skip it.
	sx := &scanner.Scanner{Prefixes: dx.OptionPrefixes, Separators: []string{}}
	if dx.OptionsArgumentsSeparator != "" {
		sx.Separators = append(sx.Separators, dx.OptionsArgumentsSeparator)
	}
	argv := append([]string{args.CommandName}, args.Args...)
	tokens := assert.NotError1(sx.Scan(argv))
	_, ok := tokens[0].(scanner.ProgramNameToken)
	assert.True(ok, "the first token must be a ProgramNameToken")
	tokens = tokens[1:]

	// Now, scan the tokens to find the subcommand name.
	commandIdx := -1
scannerLoop:
	for idx := 0; idx < len(tokens); idx++ {
		switch tokens[idx].(type) {
		case scanner.PositionalArgumentToken:
			commandIdx = idx // first positional argument stops the search
			break scannerLoop
		case scanner.OptionsArgumentsSeparatorToken:
			break scannerLoop
		}
	}

	// With a subcommand name, we're in business.
	if commandIdx >= 0 {
		// Reorder the command line arguments to move the subcommand at the beginning
		var subArgs []string
		subName := tokens[commandIdx].String()
		for idx := 0; idx < len(tokens); idx++ {
			if idx != commandIdx {
				subArgs = append(subArgs, tokens[idx].String())
			}
		}

		// Attempt an exact match with subName
		if cmd := dx.Commands[subName]; cmd != nil {
			return dx.run(ctx, cmd, args, subName, subArgs)
		}

		// Special case: synthesize `help` and `version` commands
		// when they have not been explicitly defined
		switch {
		case subName == "help":
			return dx.maybeForwardHelp(ctx, args, subArgs)

		case dx.Version != "" && subName == "version":
			return dx.handleVersionCommand(ctx, args, subArgs)
		}

		// Otherwise mention that the given command was not found
		return dx.errorNoSuchCommand(args.Env, args.CommandName, subName)
	}

	// If we have tokens, attempt to find a `--help` or `--version` or `-h`
	// flag (or equivalents considering the configured prefixes).
	for idx := 0; idx < len(tokens); idx++ {
		switch tok := tokens[idx].(type) {
		case scanner.OptionToken:
			switch {
			case tok.Name == "help" || tok.Name == "h":
				return dx.printUsage(args.Env, args.CommandName)
			case dx.Version != "" && tok.Name == "version":
				return dx.handleVersionFlag(args.Env)
			}
		}
	}

	return dx.errorInvalidFlags(args.Env, args.CommandName, args.Args)
}

// ErrInvalidFlags is returned when the command line contains invalid
// flags and no subcommand is specified.
var ErrInvalidFlags = errors.New("invalid flags")

func (dx *DispatcherCommand[T]) errorInvalidFlags(env T, commandName string, args []string) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s: invalid flags: %s\n", commandName, shellquote.Join(args...))
	fmt.Fprintf(&sb, "Try '%s --help' for more information.\n", commandName)
	fmt.Fprintln(env.Stderr(), strings.TrimSpace(sb.String()))
	return ErrInvalidFlags
}

func (dx *DispatcherCommand[T]) printUsage(env T, commandName string) error {
	_, err := fmt.Fprintln(env.Stdout(), dx.formatUsage(commandName))
	return err
}

func (dx *DispatcherCommand[T]) handleVersionFlag(env T) error {
	command := &VersionCommand[T]{Version: dx.Version}
	return command.PrintVersion(env)
}

func (dx *DispatcherCommand[T]) handleVersionCommand(
	ctx context.Context, args *CommandArgs[T], subArgs []string) error {
	command := &VersionCommand[T]{ErrorHandling: dx.ErrorHandling, Version: dx.Version}
	return dx.run(ctx, command, args, "version", subArgs)
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
	// Handle the case of autogenerated version subcommand
	case cmd == nil && dx.Version != "" && subName == "version":
		cmd = &VersionCommand[T]{ErrorHandling: dx.ErrorHandling, Version: dx.Version}
		subArgs = append(subArgs, cmd.HelpFlag())
		return dx.run(ctx, cmd, args, subName, subArgs)

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
	fmt.Fprintln(env.Stderr(), strings.TrimSpace(sb.String()))
	return ErrNoSuchCommand
}

// --- formatting code ---

func (dx *DispatcherCommand[T]) formatUsage(commandName string) string {
	// If the user configured the usage string, use it
	if dx.Usage != "" {
		return dx.Usage
	}

	// Otherwise, create a simple usage string message
	var sb strings.Builder

	// Synopsis
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "Usage: %s [command] [args]\n", commandName)

	// Description
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "%s\n", textwrap.Do(dx.BriefDescriptionText, 72, ""))

	// Commands
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "Commands:\n")
	commands := dx.cloneSubcommandsForUsage()
	for _, name := range sortedSubcommandNames(commands) {
		fmt.Fprintf(&sb, "  %s\n", name)
		cmd := commands[name]
		fmt.Fprintf(&sb, "%s\n\n", textwrap.Do(cmd.BriefDescription(), 72, "    "))
	}

	// Conclusion
	fmt.Fprintf(&sb, "Try '%s help COMMAND' for more information on COMMAND.\n", commandName)
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "Use '%s help' to show this help screen.\n", commandName)
	if dx.Version != "" {
		fmt.Fprintf(&sb, "\n")
		fmt.Fprintf(&sb, "Use '%s --version` to show the command version.\n", commandName)
	}

	return strings.TrimSpace(sb.String())
}

func (dx *DispatcherCommand[T]) cloneSubcommandsForUsage() map[string]Command[T] {
	output := maps.Clone(dx.Commands)
	if dx.Version != "" && output["version"] == nil {
		output["version"] = &VersionCommand[T]{Version: dx.Version}
	}
	return output
}

func sortedSubcommandNames[T ExecEnv](commands map[string]Command[T]) []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
