// doc.go - package documentation.
// SPDX-License-Identifier: GPL-3.0-or-later

/*
Package clip provides a framework for building command-line
applications with subcommands and flags parsing.

To use this package proceed as follows:

 1. Create a [RootCommand] containing either a [DispatcherCommand] or a
    [LeafCommand], depending on whether you need subcommands.

 2. If using a [DispatcherCommand], initialize its Commands map with
    either [LeafCommand] or [DispatcherCommand] instances.

 3. Inside a [LeafCommand] Run function, implement the flags parsing
    logic using a [*FlagSet] created with [clip.NewFlagSet].

See the package examples for more information.

# RootCommand

The [*RootCommand] optionally allows to react to signals and otherwise
dispatches the work to the [*DispatcherCommand] or [*LeafCommand]
contained within.

# DispatcherCommand

A [*DispatcherCommand] dispatches the execution to subcommands by
inspecting the command-line arguments and finding a match to a known
subcommand. By default, [*DispatcherCommand] assumes that the
command to execute is the first arguments slice entry.

However, the [*DispatcherCommand] can also handle cases in which
options and subcommands are interleaved, if you configure the
prefixes used to recognize options. For example, if following the
GNU conventions, those prefixes are `-` and `--`.

Specifically, if we know the prefixes, we scan for the first
non-option token and use it as the command to execute. If the
there is no such command, but the token is "help" we automatically
print the help message for the [*DispatcherCommand]. If the token
is "version and you configured the version string, we print the
version string. Otherwise, we return [ErrNoSuchCommand].

When no command is specified and you configured the options
prefixes, we also scan the command-line for a flag named `help`
or `h` and print help in such a case. If you configured the
version string, and we see a flag named `version`, we print it.

Otherwise, we return [ErrInvalidFlags]. Apart from these
special cases, there is no flag associated with a
[*DispatcherCommand]. This is by design: we keep the [*DispatcherCommand]
very cleanly separated from [*LeafCommand]. You will need to
manually enforce a reasonably consistent convention for
the command-line interface of your subcommands. But, on the
flip side, the resulting code is straightforward and readable.

# LeafCommand and FlagSet

To parse command line flags in a [LeafCommand], you should
use the [pkg/nflag] package, which provides a functionality
similar to the standard library `flag` package, but with
the possibility of customizing the options prefixes.

# Testability

All top-level types depend on an abstract T type, bounded by the
[ExecEnv] interface. The default implementation, using the standard
library, but highly customizable is [*StdlibExecEnv]. By using
such an interface, it is possible to write highly testable code
where most of the environment dependencies can be mocked.
*/
package clip
