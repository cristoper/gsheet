package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/cristoper/gsheet/gdrive"
	"github.com/urfave/cli/v2"
)

func CreateFolderAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("NAME is required")
	}
	name := c.Args().Get(0)
	svc, err := gdrive.NewServiceWithCtx(context.Background())
	if err != nil {
		return err
	}
	dir, err := svc.CreateFolder(name, c.String("parent"))
	if err == nil {
		fmt.Printf("Created directory named %s with id %s\n", dir.Name, dir.Id)
	}
	return err
}
