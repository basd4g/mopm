// vim:set noexpandtab :
package main

import (
	"os"
	"os/user"
)

func homeDir() string {
	if !machinePrivilege() {
		usr, err := user.Current()
		Exit1IfError(err)
		return usr.HomeDir
	}
	sudoUserName := os.Getenv("SUDO_USER")
	if sudoUserName == "" {
		Exit1("Please excute with sudo if you excute mopm by root")
	}
	usr, err := user.Lookup(sudoUserName)
	Exit1IfError(err)
	return usr.HomeDir
}

func mopmDir() string {
	mopmDir := homeDir() + "/.mopm"
	if f, err := os.Stat(mopmDir); os.IsNotExist(err) || !f.IsDir() {
		// directory '~/.mopm' is not exist
		err = os.Mkdir(mopmDir, 0777)
		Exit1IfError(err)
	}
	return mopmDir
}
