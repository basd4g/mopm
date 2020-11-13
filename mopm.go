package main

import (
  "os"
  "fmt"
  "log"
  "io/ioutil"
  "github.com/urfave/cli"
  "github.com/go-yaml/yaml"
)

type Package struct {
  Name string
  Url string
  Description string
  Environments []struct {
    Architecture string
    Platform string
    Dependencies []string
    Verification string
    Privilege string
    Script string
  }
}

func main() {
  app := &cli.App {
    Name: "mopm",
    Usage: "Mopm (Manager Of Package Maganger) is meta package manager for cross platform environment.",
    Version: "0.0.1",
    Commands: []*cli.Command{
      {
        Name: "search",
        Usage: "search package",
        Action: func (c *cli.Context) error {
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
    },
  }

  err := app.Run(os.Args)
  if err !=nil {
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
    fmt.Print( "(" + env.Architecture + ", " + env.Platform + ")" )
  }
  fmt.Println()
}
