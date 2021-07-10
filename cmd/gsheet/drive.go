package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func DeleteAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("Missing FILE_ID")
	}
	for _, id := range c.Args().Slice() {
		err := driveSvc.DeleteFile(id)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.App.ErrWriter, "Deleted file %s\n", id)
	}
	return nil
}

func CreateFolderAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("NAME is required")
	}
	name := c.Args().Get(0)
	dir, err := driveSvc.CreateFolder(name, c.String("parent"))
	if err == nil {
		fmt.Fprintf(c.App.ErrWriter, "Created directory named %s with id %s\n", dir.Name, dir.Id)
	}
	return err
}

func DownloadAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("FILE_ID is required")
	}
	id := c.Args().Get(0)
	contents, err := driveSvc.FileContents(id)
	if err != nil {
		return err
	}
	fmt.Fprint(c.App.Writer, string(contents))
	return nil
}

func InfoAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("FILE_ID is required")
	}
	id := c.Args().Get(0)
	file, err := driveSvc.GetInfo(id)
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
	file, err := driveSvc.CreateOrUpdateFile(name, c.String("parent"), inFile)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.ErrWriter, "Uploaded file as %s\n", file.Id)
	return nil
}

func ListAction(c *cli.Context) error {
	var q string
	p := c.String("parent")
	if p != "" {
		q = fmt.Sprintf("'%s' in parents", p)
	}
	files, err := driveSvc.Search(q)
	if err == nil {
		for _, f := range files {
			fmt.Fprintf(c.App.ErrWriter, "%-16s\t%1s\n", f.Name, f.Id)
		}
	}
	return err
}
