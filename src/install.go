// vim:set noexpandtab :
package main

import (
	"errors"
	"os"
)

func install(pkgName string) {
	installedAny := false
	pushInstallPkg([]string{pkgName})
	for len(installPkgStack) > 0 {
		pkgName = popInstallPkg()
		if findInstallPkg(pkgName) {
			Exit1("dependencies is looped")
		}

		env, err := findPackageEnvironment(pkgName, machineEnvId())
		Exit1IfError(err)

		if env.Verify() {
			continue
		}
		installedAny = true

		deps := env.DependenciesNotInstalled()
		if len(deps) != 0 {
			pushInstallPkg([]string{pkgName})
			pushInstallPkg(deps)
			continue
		}

		err = installExec(env.Privilege, env.Script)
		Exit1IfError(err)

		if !env.Verify() {
			Exit1("Finished installing script but failed to verify")
		}
		message("Installed " + pkgName)
	}
	if !installedAny {
		message("The package is already installed.")
		return
	}
	message("Installed successfully.")
}

var installPkgStack []string

func pushInstallPkg(ss []string) {
	installPkgStack = append(installPkgStack, ss...)
}

func popInstallPkg() string {
	ret := installPkgStack[len(installPkgStack)-1]
	installPkgStack = installPkgStack[:len(installPkgStack)-1]
	return ret
}

func findInstallPkg(str string) bool {
	for _, s := range installPkgStack {
		if s == str {
			return true
		}
	}
	return false
}

func installExec(privilege bool, script string) error {
	// | package\user | root  | unroot |
	// | ----         | ----  | ----   |
	// | root         | OK    | FAIL   |
	// | unroot       | OK(*) | OK     |
	// (*)  If mopm is runnning on sudo (Need unroot username to get with $SUDO_USER)
	if privilege == machinePrivilege() {
		return execBash(script, false)
	}
	isSudo := (machinePrivilege() && os.Getenv("SUDO_USER") != "")
	if !privilege && isSudo {
		return execBashUnsudo(script, false)
	}
	return errors.New("Check privilege to install this package")
}
