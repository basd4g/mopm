// vim:set noexpandtab :
package main

import (
	"testing"
)

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
	expected = "Package architecture must be 'darwin', 'linux' or 'linux/DISTNAME'"
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
