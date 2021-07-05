package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

func InfoAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("FILE_ID is required")
	}
	id := c.Args().Get(0)
	file, err := svc.GetInfo(id)
	if err != nil {
		return err
	}
	jsonBytes, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprint(c.App.Writer, string(jsonBytes))
	return nil
}
