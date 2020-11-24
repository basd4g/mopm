// vim:set noexpandtab :
package main

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
)

func update() {
	for _, repo := range repositories() {
		message(repo.url)
		_, err := os.Stat(repo.dir)
		if err != nil {
			message("Directory does not exist: " + repo.dir + "\nClone")
			gitClone(repo.dir, repo.url)
		} else {
			gitPull(repo.dir)
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
