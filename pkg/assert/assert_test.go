// assert_test.go - Test assertions utilities.
// SPDX-License-Identifier: GPL-3.0-or-later

package assert

import (
	"errors"
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
				switch {
				case tt.wantPanic && r != nil:
					err, ok := r.(error)
					if !ok {
						t.Fatalf("expected panic value to be error, got %T", r)
					}
					if err.Error() != tt.message {
						t.Fatalf("expected panic message %q, got %q", tt.message, err.Error())
					}

				case tt.wantPanic:
					t.Fatalf("expected panic but got none")

				case r != nil:
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			True(tt.condition, tt.message)

			if tt.wantPanic {
				t.Fatalf("expected panic but got none")
			}
		})
	}
}

func TestTrue1(t *testing.T) {
	// Test cases table
	tests := []struct {
		name      string
		condition bool
		wantPanic bool
	}{
		{
			name:      "true condition should not panic",
			condition: true,
			wantPanic: false,
		},
		{
			name:      "false condition should panic",
			condition: false,
			wantPanic: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case tt.wantPanic && r != nil:
					err, ok := r.(error)
					if !ok {
						t.Fatalf("expected panic value to be error, got %T", r)
					}
					expect := "assertion failed"
					if err.Error() != expect {
						t.Fatalf("expected panic message %q, got %q", expect, err.Error())
					}

				case tt.wantPanic:
					t.Fatalf("expected panic but got none")

				case r != nil:
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			got := True1("something", tt.condition)

			if tt.wantPanic {
				t.Fatalf("expected panic but got none")
			}

			if got != "something" {
				t.Fatalf("expected return value %q, got %q", "something", got)
			}
		})
	}
}

func TestNotError(t *testing.T) {
	// Test cases table
	tests := []struct {
		name      string
		err       error
		wantPanic bool
	}{
		{
			name:      "nil error should not panic",
			err:       nil,
			wantPanic: false,
		},
		{
			name:      "non-nil error should panic",
			err:       errors.New("mocked error"),
			wantPanic: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case tt.wantPanic && r != nil:
					err, ok := r.(error)
					if !ok {
						t.Fatalf("expected panic value to be error, got %T", r)
					}
					if !errors.Is(err, tt.err) {
						t.Fatalf("expected panic error %v, got %v", tt.err, err)
					}

				case tt.wantPanic:
					t.Fatalf("expected panic but got none")

				case r != nil:
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			NotError(tt.err)

			if tt.wantPanic {
				t.Fatalf("expected panic but got none")
			}
		})
	}
}

func TestNotError1(t *testing.T) {
	// Test cases table
	tests := []struct {
		name      string
		err       error
		wantPanic bool
	}{
		{
			name:      "nil error should not panic",
			err:       nil,
			wantPanic: false,
		},
		{
			name:      "non-nil error should panic",
			err:       errors.New("mocked error"),
			wantPanic: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				switch {
				case tt.wantPanic && r != nil:
					err, ok := r.(error)
					if !ok {
						t.Fatalf("expected panic value to be error, got %T", r)
					}
					if !errors.Is(err, tt.err) {
						t.Fatalf("expected panic error %v, got %v", tt.err, err)
					}

				case tt.wantPanic:
					t.Fatalf("expected panic but got none")

				case r != nil:
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			got := NotError1("something", tt.err)

			if tt.wantPanic {
				t.Fatalf("expected panic but got none")
			}

			if got != "something" {
				t.Fatalf("expected return value %q, got %q", "something", got)
			}
		})
	}
}
