// root_test.go - root command tests
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

import (
	"context"
	"sync"
	"syscall"
	"testing"
)

func TestRoot_Main(t *testing.T) {
	t.Run("interrupted by signal", func(t *testing.T) {
		// Create a mockable environment
		env := NewStdlibExecEnv()
		var (
			sigchan chan<- Signal
			sigmu   = &sync.Mutex{}
		)
		env.SignalNotifyFunc = func(c chan<- Signal, sig ...Signal) {
			sigmu.Lock()
			sigchan = c
			sigmu.Unlock()
		}

		// Add channel to unblock the background goroutines
		readych := make(chan struct{})

		// Create a Root instance with the mockExecEnv
		root := &RootCommand[*StdlibExecEnv]{
			Command: &LeafCommand[*StdlibExecEnv]{
				RunFunc: func(ctx context.Context, args *CommandArgs[*StdlibExecEnv]) error {
					close(readych)
					<-ctx.Done()
					return nil
				},
			},
			AutoCancel: true,
		}

		// Add background goroutine that simulates sending a signal
		go func() {
			<-readych
			sigmu.Lock()
			ch := sigchan
			sigmu.Unlock()
			ch <- syscall.SIGINT
		}()

		// Call the Main method with the mockExecEnv
		root.Main(env)
	})
}
