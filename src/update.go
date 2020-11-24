// vim:set noexpandtab :
package main

import (
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

func update(_ *cli.Context) {
	for _, url := range packageRepositories() {
		message(url)
		path := repoUrl2repoPath(url)
		_, err := os.Stat(path)
		if err != nil {
			message("Directory does not exist: " + path + "\nClone")
			gitClone(path, url)
		} else {
			gitPull(path)
		}
	}
}

func gitClone(path string, url string) {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stderr,
	})
	Exit1IfError(err)
}

func gitPull(path string) {
	r, err := git.PlainOpen(path)
	Exit1IfError(err)
	w, err := r.Worktree()
	Exit1IfError(err)
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stderr,
	})
	if err != nil && err.Error() != "already up-to-date" {
		Exit1IfError(err)
	}
}
