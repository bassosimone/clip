// leaf.go - leaf command.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import "context"

// LeafCommand is a leaf [Command]. Leaf commands take flags and
// arguments and should parse the command line.
//
// The zero value is not usable. Initialize all the mandatory fields.
type LeafCommand[T ExecEnv] struct {
	// --- mandatory fields ---

	// BriefDescriptionText is the mandatory short description.
	BriefDescriptionText string

	// RunFunc is the mandatory main function of this command.
	RunFunc func(ctx context.Context, args *CommandArgs[T]) error

	// --- optional fields ---

	// HelpFlagValue is the optional command-line flag used to request help.
	//
	// If empty, we use the default: "--help".
	HelpFlagValue string

	// LongDescriptionText is the optional long description.
	LongDescriptionText string
}

var _ Command[*StdlibExecEnv] = &LeafCommand[*StdlibExecEnv]{}

// BriefDescription implements [Command].
func (c *LeafCommand[T]) BriefDescription() string {
	return c.BriefDescriptionText
}

// HelpFlag implements [Command].
func (c *LeafCommand[T]) HelpFlag() string {
	value := c.HelpFlagValue
	if value == "" {
		value = "--help"
	}
	return value
}

// LongDescription returns the long description, if available, and
// otherwise defaults to use the brief description.
func (c *LeafCommand[T]) LongDescription() string {
	text := c.LongDescriptionText
	if text == "" {
		text = c.BriefDescriptionText
	}
	return text
}

// SupportsSubcommands implements [Command].
func (c *LeafCommand[T]) SupportsSubcommands() bool {
	return false
}

// Run implements [Command].
func (c *LeafCommand[T]) Run(ctx context.Context, args *CommandArgs[T]) error {
	return c.RunFunc(ctx, args)
}
