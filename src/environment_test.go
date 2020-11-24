// vim:set noexpandtab :
package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

//func (env Environment) String() string {
//func (env Environment) Verify() bool {
//func (env Environment) DependenciesNotInstalled() []string {

func TestFindPackageEnvironment(t *testing.T) {
	_ = os.MkdirAll("/home/mopmuser/.mopm", 0755)
	_ = ioutil.WriteFile("/home/mopmuser/.mopm/repos-url", []byte("https://github.com/basd4g/mopm.git"), 0644)
	_ = os.MkdirAll("/home/mopmuser/.mopm/repos/github.com/basd4g/mopm", 0755)
	exec.Command("cp", "-r", "/home/mopmuser/definitions", "/home/mopmuser/.mopm/repos/github.com/basd4g/mopm").Run()

	env, err := findPackageEnvironment("a", "amd64@linux/alpine_linux")
	if err != nil {
		t.Errorf("findPackageEnvironment() return err: %s", err)
	}
	if env.Architecture != "amd64" {
		t.Errorf("findPackageEnvironment() returned env.Architecture: %s, want 'amd64'", env.Architecture)
	}
	if env.Platform != "linux/alpine_linux" {
		t.Errorf("findPackageEnvironment() returned env.Platform: %s, want 'linux/alpine_linux'", env.Platform)
	}
	if env.Dependencies[0] != "b" || env.Dependencies[1] != "c" || env.Dependencies[2] != "d" {
		t.Errorf("findPackageEnvironment() returned env.Dependencies: %s, want {'b', 'c', 'd'}", env.Dependencies)
	}
	if env.Verification != "find /tmp/mopm-defs-test/a-is-installed" {
		t.Errorf("findPackageEnvironment() returned env.Verification: %s, want 'find /tmp/mopm-defs-test/a-is-installed'", env.Verification)
	}
	if env.Privilege {
		t.Errorf("findPackageEnvironment() returned env.Privilege: true, want false")
	}
	cmpTxt := "mkdir -p /tmp/mopm-defs-test\ntouch /tmp/mopm-defs-test/a-is-installed\necho \"install pkg a is finished (stdout message test)\"\necho \"install pkg a is finished (stderr message test)\" 1>&2\n"
	if env.Script != cmpTxt {
		t.Errorf("findPackageEnvironment() returned env.Script: %s, want '%s'", env.Script, cmpTxt)
	}
}
