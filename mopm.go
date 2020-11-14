package main

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
)

type Package struct {
	Name         string
	Url          string
	Description  string
	Environments []struct {
		Architecture string
		Platform     string
		Dependencies []string
		Verification string
		Privilege    string
		Script       string
	}
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
					packagePath := "definitions/" + packageName + ".mopm.yaml"
					_, err := os.Stat(packagePath)
					if err != nil {
						log.Fatal(err)
						return err
					}

					pkg, err := readPackageFile(packagePath)
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
					pkg, err := readPackageFile(packagePath)
					if err != nil {
						log.Fatal(err)
						return err
					}
					err = checkPackageFormat(pkg)
					return err
				},
			},
			{
				Name:  "environment",
				Usage: "check the machine environment",
				Action: func(c *cli.Context) error {
					platform, err := readPlatform()
					if err != nil {
						log.Fatal(err)
						return err
					}
					fmt.Println("platform:     " + platform)
					fmt.Println("architecture: " + runtime.GOARCH)
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
	// file is exist?
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal("Error: The package do not exists")
		return nil, err
	}
	// read file
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
	return &pkg, nil
}

func printPackage(pkg *Package) {
	fmt.Println("name:         " + pkg.Name)
	fmt.Println("url:          " + pkg.Url)
	fmt.Println("description:  " + pkg.Description)
	fmt.Print("environments: ")
	for i, env := range pkg.Environments {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Print("(" + env.Architecture + ", " + env.Platform + ")")
	}
	fmt.Println()
}

func checkPackageFormat(pkg *Package) error {
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
		if env.Architecture != "x86_64" {
			return errors.New("package environment architecture must be 'x86_64'")
		}
		if env.Platform != "darwin" && env.Platform != "ubuntu" {
			return errors.New("package environment architecture must be 'darwin' || 'ubuntu'")
		}
		for _, dpkg := range env.Dependencies {
			if !pkgNameRegex.MatchString(dpkg) {
				return errors.New("package environment dependencies package name must consist of a-z, 0-9 and -(hyphen) charactors")
			}
		}
		if env.Verification == "" {
			return errors.New("package environment verification must not be empty")
		}
		if env.Privilege != "root" && env.Privilege != "unnecessary" && env.Privilege != "never" {
			return errors.New("package environment architecture must be 'root', 'unnecessary' or 'never'")
		}
		if env.Script == "" {
			return errors.New("package environment script must not be empty")
		}
	}
	return nil
}

func readPlatform() (string, error) {
	if runtime.GOOS != "linux" {
		return runtime.GOOS, nil
	}
	// read file
	buf, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	for _, line := range regexp.MustCompile(`\r\n|\n\r|\n|\r`).Split(string(buf), -1) {
		if strings.HasPrefix(line, "NAME=\"") && strings.HasSuffix(line, "\"") {
			distributionName := strings.TrimSpace(strings.ToLower(line[6 : len(line)-1]))
			return "linux/" + distributionName, nil
		}
	}
	return "linux", nil
}
