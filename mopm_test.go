package main

import (
	"os"
	"testing"
)

var pkg = Package{
	Name:        "package-name",
	Url:         "https://example.com",
	Description: "This sentence is package description!",
	Environments: []Environment{
		{
			Architecture: "amd64",
			Platform:     "linux/ubuntu",
			Dependencies: []string{},
			Verification: "verificationCommand",
			Privilege:    false,
			Script:       "installComannds",
		},
	},
}

func TestMachinePlatform(t *testing.T) {
	got := machinePlatform()
	if got != "linux/ubuntu" {
		t.Errorf("machinePlatform() = %s, want 'linux/ubuntu'", got)
	}
}

func TestMachineEnvId(t *testing.T) {
	got := machineEnvId()
	if got != "amd64@linux/ubuntu" {
		t.Errorf("machineEnvId() = %s, want 'amd64@linux/ubuntu'", got)
	}
}

func TestLintPackage(t *testing.T) {
	// success
	got := lintPackage(&pkg)
	if got != nil {
		t.Errorf("lintPackage(&pkg) = %s, want nil", got)
	}

	// fail: name
	pkg.Name = "with space"
	got = lintPackage(&pkg)
	expected := "Package name must consist of a-z, 0-9 and -(hyphen) charactors"
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Name = "package-name"

	// fail: url
	pkg.Url = "ftp://example.com/hoge/fuga"
	got = lintPackage(&pkg)
	expected = "Package url must start with http(s):// ... "
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Url = "https://example.com"

	// fail: description
	pkg.Description = ""
	got = lintPackage(&pkg)
	expected = "Package description must not be empty"
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Description = "This sentence is package description!"

	// fail: architecture
	pkg.Environments[0].Architecture = "unknown"
	got = lintPackage(&pkg)
	expected = "Package architecture must be 'amd64'"
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Environments[0].Architecture = "amd64"

	// fail: platform
	pkg.Environments[0].Platform = "unknown"
	got = lintPackage(&pkg)
	expected = "Package architecture must be 'darwin' or 'linux/ubuntu'"
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Environments[0].Platform = "linux/ubuntu"

	// fail: dependencies
	pkg.Environments[0].Dependencies = []string{"valid-package-name", "unvalid package name"}
	got = lintPackage(&pkg)
	expected = "Package dependencies must consist of a-z, 0-9 and -(hyphen) charactors"
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Environments[0].Dependencies = []string{}

	// fail: verification
	pkg.Environments[0].Verification = ""
	got = lintPackage(&pkg)
	expected = "Package verification must not be empty"
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Environments[0].Verification = "verificationCommand"

	// fail: verification
	pkg.Environments[0].Script = ""
	got = lintPackage(&pkg)
	expected = "Package script must not be empty"
	if got != nil && got.Error() != expected {
		t.Errorf("lintPackage(&pkg) = '%s', want '%s'", got, expected)
	}
	pkg.Environments[0].Script = "installCommands"
}

func TestReadPackageFile(t *testing.T) {
	dir, _ := os.Getwd()
	got, err := readPackageFile(dir + "/sample.mopm.yaml")
	if err != nil {
		t.Errorf("readPackageFile() return error: %s, want nil", err)
	}
	if got.Name != "sample" {
		t.Errorf("readPackageFile() return got.Name = %s, want 'sample'", got.Name)
	}
	if got.Url != "https://github.com/basd4g/mopm" {
		t.Errorf("readPackageFile() return got.Url = %s, want 'https://github.com/basd4g/mopm'", got.Url)
	}
	if got.Description != "This is sample package definition script. It cannot be installed." {
		t.Errorf("readPackageFile() return got.Description = %s, want 'This is sample package definition script. It cannot be installed.'", got.Description)
	}
	if got.Environments[0].Architecture != "amd64" {
		t.Errorf("readPackageFile() return got.Environments[0].Architecture = %s, wants 'amd64'", got.Environments[0].Architecture)
	}
	if got.Environments[0].Platform != "darwin" {
		t.Errorf("readPackageFile() return got.Environments[0].Platform = %s, wants 'darwin'", got.Environments[0].Platform)
	}
	if got.Environments[0].Dependencies != nil {
		t.Errorf("readPackageFile() return got.Environments[0].Dependencies = %s, nil", got.Environments[0].Dependencies)
	}
	if got.Environments[0].Verification != "false && false" {
		t.Errorf("readPackageFile() return got.Environments[0].Verification = %s, wants 'false && false'", got.Environments[0].Verification)
	}
	if got.Environments[0].Privilege != false {
		t.Errorf("readPackageFile() return got.Environments[0].Privilege = true, wants false")
	}
	if got.Environments[0].Script != "echo \"This is sample install script. It is no excution anyware.\"\n" {
		t.Errorf("readPackageFile() return got.Environments[0].Script = %s, wants echo \"This is sample install script. It is no excution anyware.\"\n", got.Environments[0].Script)
	}
	if got.Environments[1].Architecture != "amd64" {
		t.Errorf("readPackageFile() return got.Environments[1].Architecture = %s, wants 'amd64'", got.Environments[1].Architecture)
	}
	if got.Environments[1].Platform != "linux/ubuntu" {
		t.Errorf("readPackageFile() return got.Environments[1].Platform = %s, wants 'linux/ubuntu'", got.Environments[1].Platform)
	}
	if got.Environments[1].Dependencies != nil {
		t.Errorf("readPackageFile() return got.Environments[1].Dependencies = %s, nil", got.Environments[1].Dependencies)
	}
	if got.Environments[1].Verification != "false && false" {
		t.Errorf("readPackageFile() return got.Environments[1].Verification = %s, wants 'false && false'", got.Environments[1].Verification)
	}
	if got.Environments[1].Privilege != true {
		t.Errorf("readPackageFile() return got.Environments[1].Privilege = false, wants true")
	}
	if got.Environments[1].Script != "echo \"This is sample install script. It is no excution anyware.\"\n" {
		t.Errorf("readPackageFile() return got.Environments[1].Script = %s, wants echo \"This is sample install script. It is no excution anyware.\"\n", got.Environments[1].Script)
	}
}

/*
TODO: Write tests for the following functions...

func main() {
func search(packageName string) error {
func lint(packagePath string) error {
func verify(packageName string) error {
func install(packageName string) error {
func printPackage(pkg *Package) {
func verifyPackage(pkg *Package) error {
func installPackage(pkg *Package) error {
func environmentOfTheMachine(pkg *Package) (*Environment, error) {
func execBash(script string) error {
func execBashUnsudo(script string) error {
func message(s string) {
*/
