// leaf_test.go - leaf command tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"errors"
	"testing"
)

func TestLeafCommand(t *testing.T) {
	t.Run("BriefDescription", func(t *testing.T) {
		expect := "foo"
		lc := &LeafCommand[*StdlibExecEnv]{BriefDescriptionText: expect}
		got := lc.BriefDescription()
		if got != expect {
			t.Errorf("BriefDescription() = %q, want %q", got, expect)
		}
	})

	t.Run("HelpFlag when explicitly set", func(t *testing.T) {
		expect := "/help"
		lc := &LeafCommand[*StdlibExecEnv]{HelpFlagValue: expect}
		got := lc.HelpFlag()
		if got != expect {
			t.Errorf("HelpFlag() = %q, want %q", got, expect)
		}
	})

	t.Run("HelpFlag default value", func(t *testing.T) {
		expect := "--help"
		lc := &LeafCommand[*StdlibExecEnv]{}
		got := lc.HelpFlag()
		if got != expect {
			t.Errorf("HelpFlag() = %q, want %q", got, expect)
		}
	})

	t.Run("LongDescription when present", func(t *testing.T) {
		expect := "foo"
		lc := &LeafCommand[*StdlibExecEnv]{LongDescriptionText: expect}
		got := lc.LongDescription()
		if got != expect {
			t.Errorf("LongDescription() = %q, want %q", got, expect)
		}
	})

	t.Run("LongDescription when not present", func(t *testing.T) {
		expect := "foo"
		lc := &LeafCommand[*StdlibExecEnv]{BriefDescriptionText: expect}
		got := lc.LongDescription()
		if got != expect {
			t.Errorf("LongDescription() = %q, want %q", got, expect)
		}
	})

	t.Run("SupportsSubcommands", func(t *testing.T) {
		expect := false
		lc := &LeafCommand[*StdlibExecEnv]{}
		got := lc.SupportsSubcommands()
		if got != expect {
			t.Errorf("SupportsSudo() = %v, want %v", got, expect)
		}
	})

	t.Run("Run", func(t *testing.T) {
		expect := errors.New("mock error")
		lc := &LeafCommand[*StdlibExecEnv]{
			RunFunc: func(ctx context.Context, args *CommandArgs[*StdlibExecEnv]) error {
				return expect
			},
		}
		err := lc.Run(context.Background(), &CommandArgs[*StdlibExecEnv]{})
		if !errors.Is(err, expect) {
			t.Errorf("Run() = %v, want %v", err, expect)
		}
	})
}
