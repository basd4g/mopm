// vim:set noexpandtab :
package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "mopm",
		Usage:   "Mopm (Manager Of Package Maganger) is meta package manager for cross platform environment.",
		Version: "0.0.2",
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
				Name:  "update",
				Usage: "download latest package definition files",
				Action: func(_ *cli.Context) {
					update()
				},
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
				Name:  "install",
				Usage: "install the package",
				Action: func(c *cli.Context) {
					install(c.Args().First())
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func search(c *cli.Context) {
	packageName := c.Args().First()
	pkgFiles, err := findAllPackageFile(packageName)
	Exit1IfError(err)
	for _, pkgFile := range pkgFiles {
		fmt.Print(pkgFile, "\n\n")
	}
}

func checkPrivilege(c *cli.Context) {
	packageName := c.Args().First()
	env, err := findPackageEnvironment(packageName, machineEnvId())
	Exit1IfError(err)
	fmt.Println(env.Privilege)
}

func verify(c *cli.Context) {
	packageName := c.Args().First()
	env, err := findPackageEnvironment(packageName, machineEnvId())
	Exit1IfError(err)
	fmt.Println(env.Verify())
}

func lint(c *cli.Context) {
	packagePath := c.Args().First()
	_, err := readPackageFile(packagePath)
	Exit1IfError(err)
	message("lint passed")
}

func message(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func Exit1IfError(err error) {
	if err != nil {
		Exit1(err.Error())
	}
}

func Exit1(s string) {
	message(s)
	os.Exit(1)
}
