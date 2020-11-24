// vim:set noexpandtab :
package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func packageRepositories() []string {
	defaultPackageRepoUrl := "https://github.com/basd4g/mopm-defs.git"

	path := mopmDir() + "/repos-url"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		message("Create the file because it does not exist: " + path)
		err = ioutil.WriteFile(path, []byte(defaultPackageRepoUrl), 0644)
		Exit1IfError(err)
	}
	buf, err := ioutil.ReadFile(path)
	Exit1IfError(err)

	var repos []string
	for _, repo := range strings.Split(string(buf), "\n") {
		if repo != "" && !strings.HasPrefix(repo, "#") {
			repos = append(repos, strings.Trim(repo, " "))
		}
	}
	if len(repos) == 0 {
		Exit1("package repository url is not found in the file: " + path)
	}
	return repos
}

func repoUrl2repoPath(url string) string {
	repo := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://"), ".git")
	return mopmDir() + "/repos/" + repo
}
