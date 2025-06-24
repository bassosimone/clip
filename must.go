// must.go - the must function
// SPDX-License-Identifier: GPL-3.0-or-later

package clip

// Must invokes exit with exit code 1 if err is not nil.
func Must[T ExecEnv](env T, err error) {
	if err != nil {
		env.Exit(1)
	}
}
