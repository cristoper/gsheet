package main

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func DownloadAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("FILE_ID is required")
	}
	//name := c.Args().Get(0)
	return nil
}
