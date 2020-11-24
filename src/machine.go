// vim:set noexpandtab :
package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func machinePlatform() string {
	if runtime.GOOS != "linux" {
		return runtime.GOOS
	}
	buf, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		panic("failed to read /etc/os-release inspite that your machine is linux")
	}
	for _, line := range regexp.MustCompile(`\r\n|\n\r|\n|\r`).Split(string(buf), -1) {
		if strings.HasPrefix(line, "NAME=\"") && strings.HasSuffix(line, "\"") {
			distributionName := strings.Replace(strings.TrimSpace(strings.ToLower(line[6:len(line)-1])), " ", "_", -1)
			return "linux/" + distributionName
		}
	}
	return "linux"
}

func machineEnvId() string {
	platform := machinePlatform()
	return runtime.GOARCH + "@" + platform
}

func machinePrivilege() bool {
	return os.Getuid() == 0
}
