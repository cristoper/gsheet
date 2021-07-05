package main

import (
	"github.com/urfave/cli/v2"
)

// Configure all of the urfave/cli commands, flags and arguments
var app = &cli.App{
	Name:  "gsheet",
	Usage: "upload and download Google Sheet data from the cli",
	Flags: []cli.Flag{},
	Commands: []*cli.Command{
		{
			Name:      "createFolder",
			Usage:     "Creates a new folder",
			ArgsUsage: "NAME",
			Action:    CreateFolderAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "parent",
					Usage:       "The id of a parent folder to act on.",
					DefaultText: "root",
				},
			},
		},
		{
			Name:   "delete",
			Usage:  "Delete a file from drive",
			Action: DeleteAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Usage:    "id of the file to delete",
					Required: true,
				},
			},
		},
		{
			Name:   "list",
			Usage:  "List file names and ids in the folder specified by --parent.",
			Action: ListAction,
		},
		{
			Name:      "upload",
			Usage:     "Upload a file to Google Drive. If a file with the same name exists in parent, it is replaced with the new file. If FILE is not given, reads from stdin (in which case --name is required).",
			Action:    UploadAction,
			ArgsUsage: "[FILE]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "name",
					Usage:       "Name to give the uploaded file",
					DefaultText: "Name of input file",
				},
			},
		},
	},
}
