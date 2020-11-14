// vim:set noexpandtab :
package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
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
					packageName := c.Args().First()
					pkg, err := readPackageFile("definitions/" + packageName + ".mopm.yaml")
					if err != nil {
						log.Fatal(err)
						return err
					}
					printPackage(pkg)
					return nil
				},
			},
			{
				Name:  "lint",
				Usage: "check package definition file format",
				Action: func(c *cli.Context) error {
					packagePath := c.Args().First()
					_, err := readPackageFile(packagePath)
					if err != nil {
						log.Fatal(err)
					}
					return err
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
					packageName := c.Args().First()
					pkg, err := readPackageFile("definitions/" + packageName + ".mopm.yaml")
					if err != nil {
						log.Fatal(err)
						return err
					}
					err = verifyPackage(pkg)
					if err != nil {
						log.Fatal(err)
						return err
					}
					return nil
				},
			},
			{
				Name:  "install",
				Usage: "install the package",
				Action: func(c *cli.Context) error {
					packageName := c.Args().First()
					pkg, err := readPackageFile("definitions/" + packageName + ".mopm.yaml")
					if err != nil {
						log.Fatal(err)
						return err
					}
					err = installPackage(pkg)
					if err != nil {
						log.Fatal(err)
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func readPackageFile(path string) (*Package, error) {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal("The package do not exists")
		return nil, err
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pkg := Package{}
	err = yaml.Unmarshal(buf, &pkg)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = lintPackage(&pkg)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &pkg, nil
}

func environmentOfTheMachine(pkg *Package) (*Environment, error) {
	machineEnvId := machineEnvId()
	for _, env := range pkg.Environments {
		if env.Architecture+"@"+env.Platform == machineEnvId {
			return &env, nil
		}
	}
	return nil, errors.New("Matched environment do not exist")
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
		envId := env.Architecture + "@" + env.Platform
		if machineEnvId == envId {
			fmt.Print("\x1b[32m" + envId + "\x1b[0m")
		} else {
			fmt.Print(envId)
		}
	}
	fmt.Println()
}

func lintPackage(pkg *Package) error {
	pkgNameRegex := regexp.MustCompile(`^[0-9a-z\-]+$`)
	if !pkgNameRegex.MatchString(pkg.Name) {
		return errors.New("package name must consist of a-z, 0-9 and -(hyphen) charactors")
	}
	urlRegex := regexp.MustCompile(`^https?://`)
	if !urlRegex.MatchString(pkg.Url) {
		return errors.New("package url must start with http(s):// ... ")
	}
	if pkg.Description == "" {
		return errors.New("package url must not be empty")
	}
	if len(pkg.Environments) == 0 {
		return errors.New("package environment must not be empty")
	}
	for _, env := range pkg.Environments {
		if env.Architecture != "amd64" {
			return errors.New("package environment architecture must be 'amd64'")
		}
		if env.Platform != "darwin" && env.Platform != "linux/ubuntu" {
			return errors.New("package environment architecture must be 'darwin' || 'linux/ubuntu'")
		}
		for _, dpkg := range env.Dependencies {
			if !pkgNameRegex.MatchString(dpkg) {
				return errors.New("package environment dependencies package name must consist of a-z, 0-9 and -(hyphen) charactors")
			}
		}
		if env.Verification == "" {
			return errors.New("package environment verification must not be empty")
		}
		if env.Privilege != true && env.Privilege != false {
			return errors.New("package environment architecture must be boolean")
		}
		if env.Script == "" {
			return errors.New("package environment script must not be empty")
		}
	}
	return nil
}

func verifyPackage(pkg *Package) error {
	env, err := environmentOfTheMachine(pkg)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		return err
	}

	if verifyPackage(pkg) == nil {
		fmt.Fprintln(os.Stderr, "The package is already installed")
		return nil
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
		log.Fatal(err)
		return err
	}

	if verifyPackage(pkg) != nil {
		err = errors.New("Finished installing script but failed to verify")
		log.Fatal(err)
		return err
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
