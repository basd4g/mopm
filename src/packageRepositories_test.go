// vim:set noexpandtab :
package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestPackageRepositories(t *testing.T) {
	_ = os.Remove("/home/mopmuser/.mopm/repos-url")
	got := packageRepositories()
	if len(got) != 1 || got[0] != "https://github.com/basd4g/mopm-defs.git" {
		t.Errorf("packageRepositories() = '%s', want 'https://github.com/basd4g/mopm-defs.git', if repos-url is not exist", got)
	}
	_ = ioutil.WriteFile("/home/mopmuser/.mopm/repos-url", []byte("#comment line\nhttps://example.com/hello.git\n#comment line\nhttp://www.example/hoge.git"), 0644)
	got = packageRepositories()
	if len(got) != 2 || got[0] != "https://example.com/hello.git" || got[1] != "http://www.example/hoge.git" {
		t.Errorf("packageRepositories() = {'%s', '%s'}, want {'https://example.com/hello.git', 'http://www.example/hoge.git'}", got[0], got[1])
	}
	_ = os.Remove("/home/mopmuser/.mopm/repos-url")

}

func TestRepoUrl2repoPath(t *testing.T) {
	got := repoUrl2repoPath("https://github.com/basd4g/mopm-defs.git")
	if got != "/home/mopmuser/.mopm/repos/github.com/basd4g/mopm-defs" {
		t.Errorf("repoUrl2repoPath(\"https://github.com/basd4g/mopm-defs.git\") = %s, want '/home/mopmuser/.mopm/repos/github.com/basd4g/mopm-defs'", got)
	}
}
