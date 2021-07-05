package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

func DeleteAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("Missing FILE_ID")
	}
	for _, id := range c.Args().Slice() {
		err := svc.DeleteFile(id)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.App.ErrWriter, "Deleted file %s\n", id)
	}
	return nil
}
