// vim:set noexpandtab :
package main

import (
	"testing"
)

func TestHomeDir(t *testing.T) {
	got := homeDir()
	if got != "/home/mopmuser" {
		t.Errorf("homeDir() = %s, want '/home/mopmuser'", got)
	}
}

func TestMopmDir(t *testing.T) {
	got := mopmDir()
	if got != "/home/mopmuser/.mopm" {
		t.Errorf("homeDir() = %s, want '/home/mopmuser/.mopm'", got)
	}
}
