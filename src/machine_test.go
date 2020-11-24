// vim:set noexpandtab :
package main

import (
	"testing"
)

func TestMachinePlatform(t *testing.T) {
	got := machinePlatform()
	if got != "linux#alpine_linux" {
		t.Errorf("machinePlatform() = %s, want 'linux#alpine_linux'", got)
	}
}

func TestMachineEnvId(t *testing.T) {
	got := machineEnvId()
	if got != "amd64@linux#alpine_linux" {
		t.Errorf("machineEnvId() = %s, want 'amd64@linux#alpine_linux'", got)
	}
}

func TestMachinePrivilege(t *testing.T) {
	got := machinePrivilege()
	if got {
		t.Errorf("machinePrivilege() = true, want false")
	}
}
