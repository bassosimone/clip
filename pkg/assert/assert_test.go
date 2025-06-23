// assert_test.go - Test assertions utilities.
// SPDX-License-Identifier: GPL-3.0-or-later

package assert

import (
	"testing"
)

func TestTrue(t *testing.T) {
	// Test cases table
	tests := []struct {
		name      string
		condition bool
		message   string
		wantPanic bool
	}{
		{
			name:      "true condition should not panic",
			condition: true,
			message:   "this should not panic",
			wantPanic: false,
		},
		{
			name:      "false condition should panic",
			condition: false,
			message:   "expected panic message",
			wantPanic: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.wantPanic {
					if r == nil {
						t.Error("expected panic but got none")
						return
					}
					err, ok := r.(error)
					if !ok {
						t.Errorf("expected panic value to be error, got %T", r)
						return
					}
					if err.Error() != tt.message {
						t.Errorf("expected panic message %q, got %q", tt.message, err.Error())
					}
				} else if r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()

			True(tt.condition, tt.message)
		})
	}
}
