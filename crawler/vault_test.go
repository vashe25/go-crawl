package crawler

import (
	"testing"
)

func TestVault(t *testing.T) {
	v := newVault()

	value := "value"
	v.addVisited(value)

	if !v.isVisited(value) {
		t.Error("Expected value to be visited")
	}
}

func TestVaultCollected(t *testing.T) {
	v := newVault()

	value := "value"
	v.collect(value)

	l := v.collected()
	if l.Length() != 1 {
		t.Error("Expected length to be 1")
	}
}
