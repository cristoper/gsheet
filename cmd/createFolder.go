package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

func CreateFolderAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("NAME is required")
	}
	name := c.Args().Get(0)
	dir, err := svc.CreateFolder(name, c.String("parent"))
	if err == nil {
		fmt.Fprintf(c.App.ErrWriter, "Created directory named %s with id %s\n", dir.Name, dir.Id)
	}
	return err
}
