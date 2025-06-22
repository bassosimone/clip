// bool_test.go - Tests for boolean flag value implementation
// SPDX-License-Identifier: GPL-3.0-or-later

package flag

import (
	"testing"
)

func TestBoolValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected string
	}{
		{
			name:     "true value",
			value:    true,
			expected: "true",
		},
		{
			name:     "false value",
			value:    false,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := boolValue{
				value: &tt.value,
			}

			result := v.String()
			if result != tt.expected {
				t.Errorf("boolValue.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestBoolValue_String_PointerModification(t *testing.T) {
	// Test that String() reflects changes to the underlying value
	value := false
	v := boolValue{
		value: &value,
	}

	// Check initial value
	if got := v.String(); got != "false" {
		t.Errorf("boolValue.String() = %q, want %q", got, "false")
	}

	// Modify the underlying value
	value = true

	// Check that String() reflects the change
	if got := v.String(); got != "true" {
		t.Errorf("boolValue.String() after modification = %q, want %q", got, "true")
	}

	// Modify back to false
	value = false

	// Check that String() reflects the change again
	if got := v.String(); got != "false" {
		t.Errorf("boolValue.String() after second modification = %q, want %q", got, "false")
	}
}

func TestBoolValue_String_Consistency(t *testing.T) {
	// Test that String() is consistent with strconv.FormatBool
	testValues := []bool{true, false}

	for _, val := range testValues {
		v := boolValue{
			value: &val,
		}

		result := v.String()

		// The String() method should produce the same result as strconv.FormatBool
		if val {
			if result != "true" {
				t.Errorf("boolValue.String() for true = %q, want %q", result, "true")
			}
		} else {
			if result != "false" {
				t.Errorf("boolValue.String() for false = %q, want %q", result, "false")
			}
		}
	}
}
