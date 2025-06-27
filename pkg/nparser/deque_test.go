// deque_test.go - tokens deque tests.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_deque(t *testing.T) {
	// Start with a deque containing three elements
	original := []Value{
		ValueProgramName{Name: "curl"},
		ValueOption{Option: &Option{Prefix: "-", Name: "o"}, Value: "FILE"},
		ValuePositionalArgument{Value: "http://www.google.com/"},
	}
	input := deque[Value]{values: original}

	// Extract from the deque like we're going to do when parsing
	var output deque[Value]
	for !input.Empty() {
		value, good := input.Front()
		if !good {
			t.Fatal("expected to be able to extract an element")
		}
		input.PopFront()
		output.PushBack(value)
	}

	// Compare the results
	if diff := cmp.Diff(original, output.values); diff != "" {
		t.Fatal(diff)
	}
}
