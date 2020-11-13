package main

import (
  "os"
  "fmt"
  "log"
  "github.com/urfave/cli"
)

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
          fmt.Printf("%s\n", c.Args().First())
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
