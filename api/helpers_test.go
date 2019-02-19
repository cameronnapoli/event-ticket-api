package main

import (
	"testing"
)

func TestConcatStrings(t *testing.T) {
	concat := concatStrings("a", "b")
	if concat != "ab" {
		t.Errorf("%s concat %s should have been %s.", "a", "b", "ab")
	}
}
