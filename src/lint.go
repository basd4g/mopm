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
	for _, env := range pkg.Environments {
		if env.Architecture != "amd64" {
			return errors.New("Package architecture must be 'amd64'")
		}
		if env.Platform != "darwin" && env.Platform != "linux/ubuntu" {
			return errors.New("Package architecture must be 'darwin' or 'linux/ubuntu'")
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
