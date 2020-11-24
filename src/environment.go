// vim:set noexpandtab :
package main

import (
	"errors"
)

type Environment struct {
	Architecture string
	Platform     string
	Dependencies []string
	Verification string
	Privilege    bool
	Script       string
}

func (env Environment) String() string {
	priv := ""
	if env.Privilege {
		priv = "(need privilege)"
	}
	envId := env.Architecture + "@" + env.Platform
	if machineEnvId() == envId {
		return "\x1b[32m" + envId + priv + "\x1b[0m"
	}
	return envId + priv
}

func (env Environment) Verify() bool {
	return execBash(env.Verification, true) == nil
}

func (env Environment) DependenciesNotInstalled() []string {
	var ret []string
	for _, depName := range env.Dependencies {
		depEnv, err := findPackageEnvironment(depName, machineEnvId())
		Exit1IfError(err)
		if !depEnv.Verify() {
			ret = append(ret, depName)
		}
	}
	return ret
}

func findPackageEnvironment(packageName string, envId string) (*Environment, error) {
	pkgFiles, err := findAllPackageFile(packageName)
	if err != nil {
		return nil, err
	}
	for _, pkgFile := range pkgFiles {
		for _, env := range pkgFile.Package.Environments {
			if env.Architecture+"@"+env.Platform == envId {
				return &env, nil
			}
		}
	}
	return nil, errors.New("Matched environment does not exist")
}
