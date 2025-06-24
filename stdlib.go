// stdlib.go - standard library execution environment.
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"io"
	"os"
	"os/signal"
)

// Signal is an alias for [os.Signal].
type Signal = os.Signal

// StdlibExecEnv is a highly-customizable [ExecEnv] whose default
// implementation uses the standard library.
//
// The zero value is not ready to use. Use [NewStdlibExecEnv]
// to create a new instance. Customize fields as needed.
type StdlibExecEnv struct {
	// OSArgs is initialized with [os.Args].
	OSArgs []string

	// OSExit is initialized with [os.Exit].
	OSExit func(exitcode int)

	// OSLookupEnv is initialized with [os.LookupEnv].
	OSLookupEnv func(key string) (string, bool)

	// SignalNotifyFunc is initialized with [signal.Notify].
	SignalNotifyFunc func(c chan<- Signal, sig ...Signal)

	// OSStderr is initialized with [os.Stderr].
	OSStderr io.Writer

	// OSStdout is initialized with [os.Stdout].
	OSStdout io.Writer

	// OSStdin is initialized with [os.Stdin].
	OSStdin io.Reader
}

var _ ExecEnv = &StdlibExecEnv{}

// NewStdlibExecEnv creates a new [StdlibExecEnv] instance.
func NewStdlibExecEnv() *StdlibExecEnv {
	return &StdlibExecEnv{
		OSArgs:           os.Args,
		OSExit:           os.Exit,
		OSLookupEnv:      os.LookupEnv,
		SignalNotifyFunc: signal.Notify,
		OSStderr:         os.Stderr,
		OSStdout:         os.Stdout,
		OSStdin:          os.Stdin,
	}
}

// Args implements [ExecEnv].
func (ee *StdlibExecEnv) Args() []string {
	return ee.OSArgs
}

// Exit implements [ExecEnv].
func (ee *StdlibExecEnv) Exit(exitcode int) {
	ee.OSExit(exitcode)
}

// LookupEnv implements [ExecEnv].
func (ee *StdlibExecEnv) LookupEnv(key string) (string, bool) {
	return ee.OSLookupEnv(key)
}

// SignalNotify implements [ExecEnv].
func (ee *StdlibExecEnv) SignalNotify(c chan<- Signal, sig ...Signal) {
	ee.SignalNotifyFunc(c, sig...)
}

// Stderr implements [ExecEnv].
func (ee *StdlibExecEnv) Stderr() io.Writer {
	return ee.OSStderr
}

// Stdin implements [ExecEnv].
func (ee *StdlibExecEnv) Stdin() io.Reader {
	return ee.OSStdin
}

// Stdout implements [ExecEnv].
func (ee *StdlibExecEnv) Stdout() io.Writer {
	return ee.OSStdout
}
