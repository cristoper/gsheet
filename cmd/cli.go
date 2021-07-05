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
			Category:  "Files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "parent",
					Usage: "The id of a parent folder to act on.",
					Value: "root",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete file(s) from drive (careful, does not trash them!)",
			ArgsUsage: "FILE_ID [FILE_ID...]",
			Action:    DeleteAction,
			Category:  "Files",
		},
		{
			Name:     "list",
			Usage:    "List file names and ids",
			Action:   ListAction,
			Category: "Files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "parent",
					Usage: "id of the folder to list (use 'root' for drive root)",
				},
			},
		},
		{
			Name:      "upload",
			Usage:     "Upload a file to Google Drive.",
			Action:    UploadAction,
			ArgsUsage: "[FILE]",
			Category:  "Files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "parent",
					Usage: "id of the folder to upload to",
					Value: "root",
				},
				&cli.StringFlag{
					Name:        "name",
					Usage:       "Name to give the uploaded file",
					DefaultText: "Name of input file",
				},
			},
		},
		{
			Name:      "download",
			Usage:     "Download a file from google drive and send it to stdout",
			Action:    DownloadAction,
			ArgsUsage: "FILE_ID",
			Category:  "Files",
		},
		{
			Name:      "info",
			Usage:     "Dump all file's metadata as json to stdout",
			Action:    InfoAction,
			ArgsUsage: "FILE_ID",
			Category:  "Files",
		},
	},
}
