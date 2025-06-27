// deque.go - generic deque.
// SPDX-License-Identifier: GPL-3.0-or-later

package nparser

type deque[T any] struct {
	values []T
}

func (d *deque[T]) Empty() bool {
	return len(d.values) <= 0
}

func (d *deque[T]) Front() (value T, ok bool) {
	if !d.Empty() {
		value, ok = d.values[0], true
	}
	return
}

func (d *deque[T]) PopFront() {
	if !d.Empty() {
		d.values = d.values[1:]
	}
}

func (d *deque[T]) PushBack(val T) {
	d.values = append(d.values, val)
}
