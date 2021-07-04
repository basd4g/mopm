// vim:set noexpandtab :
package main

import (
	"errors"
	"regexp"
)

func lintPackage(pkg *Package) error {
	pkgNameRegex := regexp.MustCompile(`^[0-9a-z\-]+$`)
	if !pkgNameRegex.MatchString(pkg.Name) {
		return errors.New("Package name must consist of a-z, 0-9 and -(hyphen) charactors")
	}
	urlRegex := regexp.MustCompile(`^https?://`)
	if !urlRegex.MatchString(pkg.Url) {
		return errors.New("Package url must start with http(s):// ... ")
	}
	if pkg.Description == "" {
		return errors.New("Package description must not be empty")
	}
	if len(pkg.Environments) == 0 {
		return errors.New("Package must not be empty")
	}
	platformRegex := regexp.MustCompile(`^(darwin|(linux(#[a-z\-_]+)?))$`)
	for _, env := range pkg.Environments {
		if env.Architecture != "amd64" && env.Architecture != "arm64" {
			return errors.New("Package architecture must be 'amd64' or 'arm64'")
		}
		if !platformRegex.MatchString(env.Platform) {
			return errors.New("Package architecture must be 'darwin', 'linux' or 'linux#DISTNAME'")
		}
		for _, dpkg := range env.Dependencies {
			if !pkgNameRegex.MatchString(dpkg) {
				return errors.New("Package dependencies must consist of a-z, 0-9 and -(hyphen) charactors")
			}
		}
		if env.Verification == "" {
			return errors.New("Package verification must not be empty")
		}
		if env.Script == "" {
			return errors.New("Package script must not be empty")
		}
	}
	return nil
}
