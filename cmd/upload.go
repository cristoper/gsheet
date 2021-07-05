package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func UploadAction(c *cli.Context) error {
	var inFile io.ReadCloser = nil
	var name string
	if c.IsSet("name") {
		name = c.String("name")
	}
	if c.NArg() < 1 {
		fmt.Fprintln(c.App.ErrWriter, "No FILE given; creating empty file on drive.")
	} else {
		switch f := c.Args().First(); f {
		case "-":
			inFile = os.Stdin
		default:
			var err error
			inFile, err = os.Open(f)
			if err != nil {
				return err
			}
			defer inFile.Close()
			if name == "" {
				name = filepath.Base(f)
			}
		}
	}
	if name == "" {
		return errors.New("Must specify --name if FILE is not given a path")
	}
	file, err := svc.CreateOrUpdateFile(name, c.String("parent"), inFile)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.ErrWriter, "Uploaded file as %s\n", file.Id)
	return nil
}
