package main

import "testing"

func TestPrintCommandRegistered(t *testing.T) {
	t.Parallel()

	if _, ok := Methods["print"]; !ok {
		t.Fatal("expected print command to be registered")
	}
}
