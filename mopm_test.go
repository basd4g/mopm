// vim:set noexpandtab :
package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMachinePlatform(t *testing.T) {
	got := machinePlatform()
	if got != "linux/alpine_linux" {
		t.Errorf("machinePlatform() = %s, want 'alpine_linux'", got)
	}
}

func TestMachineEnvId(t *testing.T) {
	got := machineEnvId()
	if got != "amd64@linux/alpine_linux" {
		t.Errorf("machineEnvId() = %s, want 'amd64@linux/alpine_linux'", got)
	}
}

func TestLintPackage(t *testing.T) {
	pkg := Package{
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

func TestPackageRepositories(t *testing.T) {
	_ = os.Remove("/home/mopmuser/.mopm/repos-url")
	got := packageRepositories()
	if len(got) != 1 || got[0] != "https://github.com/basd4g/mopm-defs.git" {
		t.Errorf("packageRepositories() = '%s', want 'https://github.com/basd4g/mopm-defs.git', if repos-url is not exist", got)
	}
	_ = ioutil.WriteFile("/home/mopmuser/.mopm/repos-url", []byte("#comment line\nhttps://example.com/hello.git\n#comment line\nhttp://www.example/hoge.git"), 0644)
	got = packageRepositories()
	if len(got) != 2 || got[0] != "https://example.com/hello.git" || got[1] != "http://www.example/hoge.git" {
		t.Errorf("packageRepositories() = {'%s', '%s'}, want {'https://example.com/hello.git', 'http://www.example/hoge.git'}", got[0], got[1])
	}
	_ = os.Remove("/home/mopmuser/.mopm/repos-url")

}

func TestRepoUrl2repoPath(t *testing.T) {
	got := repoUrl2repoPath("https://github.com/basd4g/mopm-defs.git")
	if got != "/home/mopmuser/.mopm/repos/github.com/basd4g/mopm-defs" {
		t.Errorf("repoUrl2repoPath(\"https://github.com/basd4g/mopm-defs.git\") = %s, want '/home/mopmuser/.mopm/repos/github.com/basd4g/mopm-defs'", got)
	}
}

func TestMachinePrivilege(t *testing.T) {
	got := machinePrivilege()
	if got {
		t.Errorf("machinePrivilege() = true, want false")
	}
}

/*

func (env Environment) Verify() bool {
func (env Environment) DependenciesNotInstalled() []string {
func (pkg Package) String() string {
func (pkgFile PackageFile) String() string {
func (env Environment) String() string {
func main() {
func update(_ *cli.Context) {
func gitClone(path string, url string) {
func gitPull(path string) {
func search(c *cli.Context) {
func checkPrivilege(c *cli.Context) {
func lint(c *cli.Context) {
func verify(c *cli.Context) {
func PushInstallPkg(ss []string) {
func PopInstallPkg() string {
func FindInstallPkg(str string) bool {
func install(c *cli.Context) {
func installExec(privilege bool, script string) error {
func findAllPackageFile(packageName string) ([]PackageFile, error) {
func findPackageEnvironment(packageName string, envId string) (*Environment, error) {
func execBash(script string, silently bool) error {
func execBashUnsudo(script string, silently bool) error {
func cmdRun(cmd *exec.Cmd, stdinString string, silently bool) error {
func message(s string) {
func Exit1IfError(err error) {
func Exit1(s string) {
*/
