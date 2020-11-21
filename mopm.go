// vim:set noexpandtab :
package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Environment struct {
	Architecture string
	Platform     string
	Dependencies []string
	Verification string
	Privilege    bool
	Script       string
}

func (env Environment) Verify() bool {
	return execBash(env.Verification) == nil
}

func (env Environment) DependenciesNotInstalled() []string {
	var ret []string
	for _, depName := range env.Dependencies {
		depEnv, err := findPackageEnvironment(depName, machineEnvId())
		checkIfError(err)
		if !depEnv.Verify() {
			ret = append(ret, depName)
		}
	}
	return ret
}

type Package struct {
	Name         string
	Url          string
	Description  string
	Environments []Environment
}

type PackageFile struct {
	Package *Package
	Path    string
}

func (pkg Package) String() string {
	out := new(bytes.Buffer)
	fmt.Fprintf(out, "name:         %s\n", pkg.Name)
	fmt.Fprintf(out, "url:          %s\n", pkg.Url)
	fmt.Fprintf(out, "description:  %s\n", pkg.Description)
	fmt.Fprintf(out, "environments: ")

	for i, env := range pkg.Environments {
		if i != 0 {
			fmt.Fprint(out, ", ")
		}
		fmt.Fprint(out, env)
	}
	return string(out.Bytes())
}

func (pkgFile PackageFile) String() string {
	return fmt.Sprintf("path:         %s\n%s", pkgFile.Path, pkgFile.Package)
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

func main() {
	app := &cli.App{
		Name:    "mopm",
		Usage:   "Mopm (Manager Of Package Maganger) is meta package manager for cross platform environment.",
		Version: "0.0.1",
		Commands: []cli.Command{
			{
				Name:   "search",
				Usage:  "search package",
				Action: search,
			},
			{
				Name:   "check-privilege",
				Usage:  "check package to need privilege on this envirionment",
				Action: checkPrivilege,
			},

			{
				Name:   "update",
				Usage:  "download latest package definition files",
				Action: update,
			},
			{
				Name:   "lint",
				Usage:  "check package definition file format",
				Action: lint,
			},
			{
				Name:    "environment",
				Aliases: []string{"env"},
				Usage:   "check the machine environment",
				Action: func(_ *cli.Context) {
					fmt.Println(machineEnvId())
				},
			},
			{
				Name:    "verify",
				Aliases: []string{"vrf"},
				Usage:   "verify the package to be installed or not",
				Action:  verify,
			},
			{
				Name:   "install",
				Usage:  "install the package",
				Action: install,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

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
	checkIfError(err)
}

func gitPull(path string) {
	r, err := git.PlainOpen(path)
	checkIfError(err)
	w, err := r.Worktree()
	checkIfError(err)
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stderr,
	})
	if err != nil && err.Error() != "already up-to-date" {
		checkIfError(err)
	}
	message("hello")
}

func search(c *cli.Context) {
	packageName := c.Args().First()
	pkgFiles, err := findAllPackageFile(packageName)
	checkIfError(err)
	for _, pkgFile := range pkgFiles {
		fmt.Println(pkgFile, "\n")
	}
}

func checkPrivilege(c *cli.Context) {
	packageName := c.Args().First()
	env, err := findPackageEnvironment(packageName, machineEnvId())
	checkIfError(err)
	fmt.Println(env.Privilege)
}

func lint(c *cli.Context) {
	packagePath := c.Args().First()
	_, err := readPackageFile(packagePath)
	checkIfError(err)
	message("lint passed")
}

func verify(c *cli.Context) {
	packageName := c.Args().First()
	env, err := findPackageEnvironment(packageName, machineEnvId())
	checkIfError(err)
	fmt.Println(env.Verify())
}

var installPkgStack []string

func PushInstallPkg(ss []string) {
	installPkgStack = append(installPkgStack, ss...)
}

func PopInstallPkg() string {
	ret := installPkgStack[len(installPkgStack)-1]
	installPkgStack = installPkgStack[:len(installPkgStack)-1]
	return ret
}

func FindInstallPkg(str string) bool {
	for _, s := range installPkgStack {
		if s == str {
			return true
		}
	}
	return false
}

func install(c *cli.Context) {
	installedAny := false
	pkgName := c.Args().First()
	PushInstallPkg([]string{pkgName})
	for len(installPkgStack) > 0 {
		pkgName = PopInstallPkg()
		if FindInstallPkg(pkgName) {
			checkIfError(errors.New("dependencies is looped"))
		}

		env, err := findPackageEnvironment(pkgName, machineEnvId())
		checkIfError(err)

		if env.Verify() {
			continue
		}
		installedAny = true

		deps := env.DependenciesNotInstalled()
		if len(deps) != 0 {
			PushInstallPkg([]string{pkgName})
			PushInstallPkg(deps)
			continue
		}

		err = installExec(env.Privilege, env.Script)
		checkIfError(err)

		if !env.Verify() {
			checkIfError(errors.New("Finished installing script but failed to verify"))
		}
		message("Installed " + pkgName)
	}
	if !installedAny {
		message("The package is already installed.")
		return
	}
	message("Installed successfully.")
}

func installExec(privilege bool, script string) error {
	// | package\user | root  | unroot |
	// | ----         | ----  | ----   |
	// | root         | OK    | FAIL   |
	// | unroot       | OK(*) | OK     |
	// (*)  If mopm is runnning on sudo (Need unroot username to get with $SUDO_USER)
	if privilege == machinePrivilege() {
		return execBash(script)
	}
	isSudo := (machinePrivilege() && os.Getenv("SUDO_USER") != "")
	if !privilege && isSudo {
		return execBashUnsudo(script)
	}
	return errors.New("Check privilege to install this package")
}

func homeDir() string {
	if !machinePrivilege() {
		usr, err := user.Current()
		checkIfError(err)
		return usr.HomeDir
	}
	sudoUserName := os.Getenv("SUDO_USER")
	if sudoUserName == "" {
		checkIfError(errors.New("Please excute with sudo if you excute mopm by root"))
	}
	usr, err := user.Lookup(sudoUserName)
	checkIfError(err)
	return usr.HomeDir
}

func mopmDir() string {
	mopmDir := homeDir() + "/.mopm"
	if f, err := os.Stat(mopmDir); os.IsNotExist(err) || !f.IsDir() {
		// directory '~/.mopm' is not exist
		err = os.Mkdir(mopmDir, 0777)
		checkIfError(err)
	}
	return mopmDir
}

func packageRepositories() []string {
	defaultPackageRepoUrl := "https://github.com/basd4g/mopm-defs.git"

	path := mopmDir() + "/repos-url"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		message("Create the file because it does not exist: " + path)
		err = ioutil.WriteFile(path, []byte(defaultPackageRepoUrl), 0644)
		checkIfError(err)
	}
	buf, err := ioutil.ReadFile(path)
	checkIfError(err)

	var repos []string
	for _, repo := range strings.Split(string(buf), "\n") {
		if repo != "" && !strings.HasPrefix(repo, "#") {
			repos = append(repos, strings.Trim(repo, " "))
		}
	}
	if len(repos) == 0 {
		checkIfError(errors.New("package repository url is not found in the file: " + path))
	}
	return repos
}

func repoUrl2repoPath(url string) string {
	repo := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://"), ".git")
	return mopmDir() + "/repos/" + repo
}

func findAllPackageFile(packageName string) ([]PackageFile, error) {
	var pkgFiles []PackageFile
	for _, url := range packageRepositories() {
		path := repoUrl2repoPath(url) + "/definitions/" + packageName + ".yaml"
		pkgFile, err := readPackageFile(path)
		if err == nil {
			pkgFiles = append(pkgFiles, pkgFile)
		}
	}
	return pkgFiles, nil
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

func readPackageFile(path string) (PackageFile, error) {
	_, err := os.Stat(path)
	if err != nil {
		return PackageFile{}, fmt.Errorf("The package does not exist: %s\nWrapped: %w", path, err)
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return PackageFile{}, err
	}

	pkg := Package{}
	err = yaml.Unmarshal(buf, &pkg)
	if err != nil {
		return PackageFile{}, fmt.Errorf("Failed to parse yaml file: %s\nWrapped: %w", path, err)
	}
	err = lintPackage(&pkg)
	if err != nil {
		return PackageFile{}, err
	}
	return PackageFile{
		Package: &pkg,
		Path:    path,
	}, nil
}

func lintPackage(pkg *Package) error {
	pkgNameRegex := regexp.MustCompile(`^[0-9a-z\-]+$`)
	if !pkgNameRegex.MatchString(pkg.Name) {
		return errors.New("Package name must consist of a-z, 0-9 and -(hyphen) charactors")
	}
	urlRegex := regexp.MustCompile(`^https?://`)
	if !urlRegex.MatchString(pkg.Url) {
		return errors.New("Package url must start with http(s):// ... ")
	}
	if pkg.Description == "" {
		return errors.New("Package description must not be empty")
	}
	if len(pkg.Environments) == 0 {
		return errors.New("Package must not be empty")
	}
	for _, env := range pkg.Environments {
		if env.Architecture != "amd64" {
			return errors.New("Package architecture must be 'amd64'")
		}
		if env.Platform != "darwin" && env.Platform != "linux/ubuntu" {
			return errors.New("Package architecture must be 'darwin' or 'linux/ubuntu'")
		}
		for _, dpkg := range env.Dependencies {
			if !pkgNameRegex.MatchString(dpkg) {
				return errors.New("Package dependencies must consist of a-z, 0-9 and -(hyphen) charactors")
			}
		}
		if env.Verification == "" {
			return errors.New("Package verification must not be empty")
		}
		if env.Script == "" {
			return errors.New("Package script must not be empty")
		}
	}
	return nil
}

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
			distributionName := strings.TrimSpace(strings.ToLower(line[6 : len(line)-1]))
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

func execBash(script string) error {
	cmd := exec.Command("bash")
	return cmdRun(cmd, "#!/bin/bash -e\n"+script+"\n")
}

func execBashUnsudo(script string) error {
	cmd := exec.Command("sudo", "--user="+os.Getenv("SUDO_USER"), "bash")
	return cmdRun(cmd, "#!/bin/bash -e\n"+script+"\n")
}

func cmdRun(cmd *exec.Cmd, stdinString string) error {
	mopmDir := mopmDir()

	cmd.Stdin = bytes.NewBufferString(stdinString)
	logFile, err := os.OpenFile(mopmDir+"/stdout.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkIfError(err)
	fmt.Fprintf(logFile, "#MOPM:LOG:TIME----- %s -----\n", time.Now())
	cmd.Stdout = io.MultiWriter(os.Stdout, logFile)

	logFileError, err := os.OpenFile(mopmDir+"/stderr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkIfError(err)
	fmt.Fprintf(logFileError, "#MOPM:LOG:TIME----- %s -----\n", time.Now())
	cmd.Stderr = io.MultiWriter(os.Stderr, logFileError)
	return cmd.Run()
}

func message(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func checkIfError(err error) {
	if err != nil {
		message(err.Error())
		os.Exit(1)
	}
}
