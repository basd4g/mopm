// vim:set noexpandtab :
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func execBash(script string, silently bool) error {
	cmd := exec.Command("bash")
	return cmdRun(cmd, "#!/bin/bash -e\n"+script+"\n", silently)
}

func execBashUnsudo(script string, silently bool) error {
	cmd := exec.Command("sudo", "--user="+os.Getenv("SUDO_USER"), "bash")
	return cmdRun(cmd, "#!/bin/bash -e\n"+script+"\n", silently)
}

func cmdRun(cmd *exec.Cmd, stdinString string, silently bool) error {
	cmd.Stdin = bytes.NewBufferString(stdinString)
	if silently {
		return cmd.Run()
	}

	mopmDir := mopmDir()

	logFile, err := os.OpenFile(mopmDir+"/stdout.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Exit1IfError(err)
	fmt.Fprintf(logFile, "#MOPM:LOG:TIME----- %s -----\n", time.Now())
	cmd.Stdout = io.MultiWriter(os.Stdout, logFile)

	logFileError, err := os.OpenFile(mopmDir+"/stderr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Exit1IfError(err)
	fmt.Fprintf(logFileError, "#MOPM:LOG:TIME----- %s -----\n", time.Now())
	cmd.Stderr = io.MultiWriter(os.Stderr, logFileError)

	return cmd.Run()
}
