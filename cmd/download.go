package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

func DownloadAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("FILE_ID is required")
	}
	id := c.Args().Get(0)
	contents, err := svc.FileContents(id)
	if err != nil {
		return err
	}
	fmt.Fprint(c.App.Writer, string(contents))
	return nil
}
