// version_test.go - tests for version handling
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"errors"
	"testing"

	"github.com/bassosimone/clip/pkg/flag"
	"github.com/bassosimone/clip/pkg/parser"
)

func TestVersionCommand(t *testing.T) {
	t.Run("BriefDescription", func(t *testing.T) {
		expect := "foo"
		vc := &VersionCommand[*StdlibExecEnv]{BriefDescriptionText: expect}
		got := vc.BriefDescription()
		if got != expect {
			t.Fatalf("BriefDescription() = %q, want %q", got, expect)
		}
	})

	t.Run("HelpFlag when explicitly set", func(t *testing.T) {
		expect := "/help"
		vc := &VersionCommand[*StdlibExecEnv]{HelpFlagValue: expect}
		got := vc.HelpFlag()
		if got != expect {
			t.Fatalf("HelpFlag() = %q, want %q", got, expect)
		}
	})

	t.Run("HelpFlag default value", func(t *testing.T) {
		expect := "--help"
		vc := &VersionCommand[*StdlibExecEnv]{}
		got := vc.HelpFlag()
		if got != expect {
			t.Fatalf("HelpFlag() = %q, want %q", got, expect)
		}
	})

	t.Run("SupportsSubcommands", func(t *testing.T) {
		expect := false
		vc := &VersionCommand[*StdlibExecEnv]{}
		got := vc.SupportsSubcommands()
		if got != expect {
			t.Fatalf("SupportsSubcommands() = %v, want %v", got, expect)
		}
	})

	t.Run("Run with unexpected flags", func(t *testing.T) {
		vc := &VersionCommand[*StdlibExecEnv]{}
		args := &CommandArgs[*StdlibExecEnv]{
			Args:    []string{"--unexpected"},
			Command: vc,
			Env:     NewStdlibExecEnv(),
		}
		err := vc.Run(context.Background(), args)
		if !errors.Is(err, parser.ErrUnknownOption) {
			t.Fatalf("Run() = %v, want %v", err, parser.ErrUnknownOption)
		}
	})

	t.Run("Run with unexpected positional arguments", func(t *testing.T) {
		vc := &VersionCommand[*StdlibExecEnv]{}
		args := &CommandArgs[*StdlibExecEnv]{
			Args:    []string{"unexpected"},
			Command: vc,
			Env:     NewStdlibExecEnv(),
		}
		err := vc.Run(context.Background(), args)
		if !errors.Is(err, flag.ErrUnexpectedNumberOfPositionalArgs) {
			t.Fatalf("Run() = %v, want %v", err, parser.ErrUnknownOption)
		}
	})

	t.Run("Run with no arguments and no flags", func(t *testing.T) {
		vc := &VersionCommand[*StdlibExecEnv]{}
		args := &CommandArgs[*StdlibExecEnv]{
			Args:    []string{},
			Command: vc,
			Env:     NewStdlibExecEnv(),
		}
		err := vc.Run(context.Background(), args)
		if err != nil {
			t.Fatal(err)
		}
	})
}
