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

func TestHomeDir(t *testing.T) {
	got := homeDir()
	if got != "/home/basd4g" {
		t.Errorf("homeDir() = %s, want '/home/basd4g'", got)
	}
}

func TestPackageRepositories(t *testing.T) {
	got := packageRepositories()
	if len(got) != 1 || got[0] != "https://github.com/basd4g/mopm-defs.git" {
		t.Errorf("packageRepositories() = %s, want 'https://github.com/basd4g/mopm-defs.git'", got)
	}
}

func TestRepoUrl2repoPath(t *testing.T) {
	got := repoUrl2repoPath("https://github.com/basd4g/mopm-defs.git")
	if got != "/home/basd4g/.mopm/github.com/basd4g/mopm-defs" {
		t.Errorf("repoUrl2repoPath(\"https://github.com/basd4g/mopm-defs.git\") = %s, want '/home/basd4g/.mopm/github.com/basd4g/mopm-defs'", got)
	}
}

/*
func (pkg Package) String() string {
func (pkgFile PackageFile) String() string {
func (env Environment) String() string {
func main() {
func update() error {
func gitClone(path string, url string) {
func gitPull(path string) {
func search(packageName string) error {
func lint(packagePath string) error {
func verify(packageName string) error {
func verifyExec(env *Environment) error {
func install(packageName string) error {
func installExec(privilege bool, script string) error {
func findAllPackageFile(packageName string) ([]PackageFile, error) {
func findPackageEnvironment(packageName string, envId string) (*Environment, error) {
func readPackageFile(path string) (PackageFile, error) {
func lintPackage(pkg *Package) error {
func machinePlatform() string {
func machineEnvId() string {
func machinePrivilege() bool {
:q
func execBashFunc(script string) error {
func execBashUnsudoFunc(script string) error {
func message(s string) {
func checkIfError(err error) {
*/
