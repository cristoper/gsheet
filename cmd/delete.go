package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func DeleteAction(c *cli.Context) error {
	id := c.String("id")
	err := svc.DeleteFile(id)
	if err != nil {
		fmt.Fprintf(c.App.ErrWriter, "Deleted file %s\n", id)
	}
	return err
}
