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
          _, error := os.Stat(packagePath)
          if error != nil {
            log.Fatal("Error: The package did not exists")
            return error
          }
          readFile(packagePath)
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

func readFile (path string) {
  buf, err := ioutil.ReadFile(path)
  if err != nil {
    log.Fatal(err)
    return
  }

  data, err := ReadOnStruct(buf)
  if err != nil {
    log.Fatal(err)
    return
  }
  fmt.Println(*data)
}

func ReadOnStruct(fileBuffer []byte) (*Package, error) {
  p := Package{}
  err := yaml.Unmarshal(fileBuffer, &p)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  return &p, nil
}
