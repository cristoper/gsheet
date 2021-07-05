package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func ListAction(c *cli.Context) error {
	fmt.Printf("args: %s\n", c.Args())
	return nil
}
