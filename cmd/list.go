package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func ListAction(c *cli.Context) error {
	var q string
	p := c.String("parent")
	if p != "" {
		q = fmt.Sprintf("'%s' in parents", p)
	}
	files, err := svc.Search(q)
	if err == nil {
		for _, f := range files {
			fmt.Fprintf(c.App.ErrWriter, "%-16s\t%1s\n", f.Name, f.Id)
		}
	}
	return err
}
