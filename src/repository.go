// vim:set noexpandtab :
package main

import (
	"io/ioutil"
	"os"
	"strings"
)

type Repository struct {
	url string
	dir string
}

func repositories() []Repository {
	defaultPackageRepoUrl := "https://github.com/basd4g/mopm-defs.git"

	path := mopmDir() + "/repos-url"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		message("Create the file because it does not exist: " + path)
		err = ioutil.WriteFile(path, []byte(defaultPackageRepoUrl), 0644)
		Exit1IfError(err)
	}
	buf, err := ioutil.ReadFile(path)
	Exit1IfError(err)

	var repos []Repository
	for _, line := range strings.Split(string(buf), "\n") {
		if line != "" && !strings.HasPrefix(line, "#") {
			url :=  strings.Trim(line, " ")
			repos = append(repos, Repository{
				url: url,
				dir: repoDir(url),
			})
		}
	}
	if len(repos) == 0 {
		Exit1("package repository url is not found in the file: " + path)
	}
	return repos
}

func repoDir(url string) string {
	repo := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://"), ".git")
	return mopmDir() + "/repos/" + repo
}
