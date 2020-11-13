package main

import (
  "os"
  "fmt"
  "github.com/urfave/cli"
)

func main() {
  app := cli.NewApp()

  app.Name = "mopm"
  app.Usage = "Mopm (Manager Of Package Maganger) is meta package manager for cross platform environment."
  app.Version = "0.0.1"
  app.Action = func  (context *cli.Context) error {
    fmt.Println("hello, world!")
    return nil
  }

  app.Run(os.Args)
}
