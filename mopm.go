// vim:set noexpandtab :
package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type Environment struct {
	Architecture string
	Platform     string
	Dependencies []string
	Verification string
	Privilege    bool
	Script       string
}

type Package struct {
	Name         string
	Url          string
	Description  string
	Environments []Environment
}

func main() {
	app := &cli.App{
		Name:    "mopm",
		Usage:   "Mopm (Manager Of Package Maganger) is meta package manager for cross platform environment.",
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:  "search",
				Usage: "search package",
				Action: func(c *cli.Context) error {
					return search(c.Args().First())
				},
			},
			{
				Name:  "lint",
				Usage: "check package definition file format",
				Action: func(c *cli.Context) error {
					return lint(c.Args().First())
				},
			},
			{
				Name:    "environment",
				Aliases: []string{"env"},
				Usage:   "check the machine environment",
				Action: func(c *cli.Context) error {
					fmt.Println(machineEnvId())
					return nil
				},
			},
			{
				Name:    "verify",
				Aliases: []string{"vrf"},
				Usage:   "verify the package to be installed or not",
				Action: func(c *cli.Context) error {
					return verify(c.Args().First())
				},
			},
			{
				Name:  "install",
				Usage: "install the package",
				Action: func(c *cli.Context) error {
					return install(c.Args().First())
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func search(packageName string) error {
	pkg, err := readPackageFile("definitions/" + packageName + ".mopm.yaml")
	if err != nil {
		message(err.Error())
		return err
	}
	printPackage(pkg)
	return nil
}

func lint(packagePath string) error {
	_, err := readPackageFile(packagePath)
	if err != nil {
		message(err.Error())
	} else {
		message("lint passed")
	}
	return err
}

func verify(packageName string) error {
	pkg, err := readPackageFile("definitions/" + packageName + ".mopm.yaml")
	if err != nil {
		message(err.Error())
		return err
	}
	return verifyPackage(pkg)
}

func install(packageName string) error {
	pkg, err := readPackageFile("definitions/" + packageName + ".mopm.yaml")
	if err != nil {
		message(err.Error())
		return err
	}
	err = installPackage(pkg)
	if err != nil {
		message(err.Error())
		if err.Error() == "The package is already installed" {
			return nil
		}
		return err
	}
	message("Installed successfully.")
	return nil
}

func readPackageFile(path string) (*Package, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("The package do not exist: %s\nWrapped: %w", path, err)
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pkg := Package{}
	err = yaml.Unmarshal(buf, &pkg)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse yaml file: %s\nWrapped: %w", path, err)
	}
	err = lintPackage(&pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func printPackage(pkg *Package) {
	fmt.Println("name:         " + pkg.Name)
	fmt.Println("url:          " + pkg.Url)
	fmt.Println("description:  " + pkg.Description)
	fmt.Print("environments: ")
	machineEnvId := machineEnvId()
	for i, env := range pkg.Environments {
		if i != 0 {
			fmt.Print(", ")
		}
		needPrivilege := ""
		if env.Privilege {
			needPrivilege = "(need privilege)"
		}
		envId := env.Architecture + "@" + env.Platform
		if machineEnvId == envId {
			fmt.Print("\x1b[32m" + envId + needPrivilege + "\x1b[0m")
		} else {
			fmt.Print(envId + needPrivilege)
		}
	}
	fmt.Println()
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
		return errors.New("Package url must not be empty")
	}
	if len(pkg.Environments) == 0 {
		return errors.New("Package environment must not be empty")
	}
	for _, env := range pkg.Environments {
		if env.Architecture != "amd64" {
			return errors.New("Package environment architecture must be 'amd64'")
		}
		if env.Platform != "darwin" && env.Platform != "linux/ubuntu" {
			return errors.New("Package environment architecture must be 'darwin' || 'linux/ubuntu'")
		}
		for _, dpkg := range env.Dependencies {
			if !pkgNameRegex.MatchString(dpkg) {
				return errors.New("Package environment dependencies package name must consist of a-z, 0-9 and -(hyphen) charactors")
			}
		}
		if env.Verification == "" {
			return errors.New("Package environment verification must not be empty")
		}
		if env.Privilege != true && env.Privilege != false {
			return errors.New("Package environment architecture must be boolean")
		}
		if env.Script == "" {
			return errors.New("Package environment script must not be empty")
		}
	}
	return nil
}

func verifyPackage(pkg *Package) error {
	env, err := environmentOfTheMachine(pkg)
	if err != nil {
		return err
	}
	err = execBash(env.Verification)
	if err != nil {
		return errors.New("The package is not installed")
	}
	return nil
}

func installPackage(pkg *Package) error {
	env, err := environmentOfTheMachine(pkg)
	if err != nil {
		return err
	}

	if verifyPackage(pkg) == nil {
		return errors.New("The package is already installed")
	}

	// | package\user | root  | unroot |
	// | ----         | ----  | ----   |
	// | root         | OK    | FAIL   |
	// | unroot       | OK(*) | OK     |
	// (*)  If mopm is runnning on sudo (Need unroot username to get with $SUDO_USER)
	machinePrivilege := (os.Getuid() == 0)
	isSudo := (machinePrivilege && os.Getenv("SUDO_USER") != "")

	if env.Privilege == machinePrivilege {
		err = execBash(env.Script)
	} else if !env.Privilege && isSudo {
		err = execBashUnsudo(env.Script)
	} else {
		err = errors.New("Check privilege to install this package")
	}
	if err != nil {
		return err
	}

	if verifyPackage(pkg) != nil {
		return errors.New("Finished installing script but failed to verify")
	}
	return nil
}

func environmentOfTheMachine(pkg *Package) (*Environment, error) {
	machineEnvId := machineEnvId()
	for _, env := range pkg.Environments {
		if env.Architecture+"@"+env.Platform == machineEnvId {
			return &env, nil
		}
	}
	return nil, errors.New("Matched environment does not exist")
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

func execBash(script string) error {
	cmd := exec.Command("bash")
	cmd.Stdin = bytes.NewBufferString("#!/bin/bash -e\n" + script + "\n")
	return cmd.Run()
}

func execBashUnsudo(script string) error {
	cmd := exec.Command("sudo", "--user="+os.Getenv("SUDO_USER"), "bash")
	cmd.Stdin = bytes.NewBufferString("#!/bin/bash -e\n" + script + "\n")
	return cmd.Run()
}

func message(s string) {
	fmt.Fprintln(os.Stderr, s)
}
