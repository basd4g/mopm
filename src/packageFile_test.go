// vim:set noexpandtab :
package main

import (
	"os"
	"testing"
)

func TestReadPackageFile(t *testing.T) {
	dir, _ := os.Getwd()
	got, err := readPackageFile(dir + "/../definitions/sample.yaml")
	if err != nil {
		t.Errorf("readPackageFile() return error: %s, want nil", err)
	}
	if got.Package.Name != "sample" {
		t.Errorf("readPackageFile() return got.Package.Name = %s, want 'sample'", got.Package.Name)
	}
	if got.Package.Url != "https://github.com/basd4g/mopm" {
		t.Errorf("readPackageFile() return got.Package.Url = %s, want 'https://github.com/basd4g/mopm'", got.Package.Url)
	}
	if got.Package.Description != "This is sample package definition script. It cannot be installed." {
		t.Errorf("readPackageFile() return got.Package.Description = %s, want 'This is sample package definition script. It cannot be installed.'", got.Package.Description)
	}
	if got.Package.Environments[0].Architecture != "amd64" {
		t.Errorf("readPackageFile() return got.Package.Environments[0].Architecture = %s, wants 'amd64'", got.Package.Environments[0].Architecture)
	}
	if got.Package.Environments[0].Platform != "darwin" {
		t.Errorf("readPackageFile() return got.Package.Environments[0].Platform = %s, wants 'darwin'", got.Package.Environments[0].Platform)
	}
	if got.Package.Environments[0].Dependencies != nil {
		t.Errorf("readPackageFile() return got.Package.Environments[0].Dependencies = %s, nil", got.Package.Environments[0].Dependencies)
	}
	if got.Package.Environments[0].Verification != "false && false" {
		t.Errorf("readPackageFile() return got.Package.Environments[0].Verification = %s, wants 'false && false'", got.Package.Environments[0].Verification)
	}
	if got.Package.Environments[0].Privilege != false {
		t.Errorf("readPackageFile() return got.Package.Environments[0].Privilege = true, wants false")
	}
	if got.Package.Environments[0].Script != "echo \"This is sample install script. It is no excution anyware.\"\n" {
		t.Errorf("readPackageFile() return got.Package.Environments[0].Script = %s, wants echo \"This is sample install script. It is no excution anyware.\"\n", got.Package.Environments[0].Script)
	}
	if got.Package.Environments[1].Architecture != "amd64" {
		t.Errorf("readPackageFile() return got.Package.Environments[1].Architecture = %s, wants 'amd64'", got.Package.Environments[1].Architecture)
	}
	if got.Package.Environments[1].Platform != "linux/ubuntu" {
		t.Errorf("readPackageFile() return got.Package.Environments[1].Platform = %s, wants 'linux/ubuntu'", got.Package.Environments[1].Platform)
	}
	if got.Package.Environments[1].Dependencies != nil {
		t.Errorf("readPackageFile() return got.Package.Environments[1].Dependencies = %s, nil", got.Package.Environments[1].Dependencies)
	}
	if got.Package.Environments[1].Verification != "false && false" {
		t.Errorf("readPackageFile() return got.Package.Environments[1].Verification = %s, wants 'false && false'", got.Package.Environments[1].Verification)
	}
	if got.Package.Environments[1].Privilege != true {
		t.Errorf("readPackageFile() return got.Package.Environments[1].Privilege = false, wants true")
	}
	if got.Package.Environments[1].Script != "echo \"This is sample install script. It is no excution anyware.\"\n" {
		t.Errorf("readPackageFile() return got.Package.Environments[1].Script = %s, wants echo \"This is sample install script. It is no excution anyware.\"\n", got.Package.Environments[1].Script)
	}
}
