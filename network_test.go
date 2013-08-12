package main

import (
	"testing"
)

func TestNetworkIDUnique(t *testing.T) {
	const iterations = 1000

	ch := make(chan uint64)

	for i := 0; i < iterations; i++ {
		go func() {
			var n networkID
			ch <- n.NetworkID()
		}()
	}

	seen := make(map[uint64]bool, iterations)
	for i := 0; i < iterations; i++ {
		id := <-ch
		if seen[id] {
			t.Errorf("ID %d duplicated", id)
		}
		seen[id] = true
	}
}

func TestNetworkIDDeterministic(t *testing.T) {
	const iterations = 1000

	ch := make(chan uint64)

	var n networkID
	for i := 0; i < iterations; i++ {
		go func() {
			ch <- n.NetworkID()
		}()
	}

	expected := <-ch
	for i := 1; i < iterations; i++ {
		id := <-ch
		if id != expected {
			t.Errorf("ID %d: %d != %d", i, id, expected)
		}
	}
}
