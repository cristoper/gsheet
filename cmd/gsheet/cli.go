package main

import (
	"github.com/urfave/cli/v2"
)

// Set with -ldflags "-X main.version=VERSION"
var version string

// Configure all of the urfave/cli commands, flags and arguments
var app = &cli.App{
	Name:    "gsheet",
	Version: version,
	Usage:   "upload and download Google Sheet data from the cli",
	Flags:   []cli.Flag{},
	Commands: []*cli.Command{
		// Sheets
		{
			Name:     "csv",
			Usage:    "Pipe csv data to range or read it from range",
			Action:   rangeSheetAction,
			Category: "Sheets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "id",
					Usage:   "id of the spreadsheet document",
					EnvVars: []string{"GSHEET_ID"},
				},
				&cli.StringFlag{
					Name:  "range",
					Usage: "Sheet range to update or get (A1 notation)",
				},
				&cli.BoolFlag{
					Name:  "append",
					Usage: "If set, append to end of any data in range",
				},
				&cli.StringFlag{
					Name:  "sep",
					Value: ",",
					Usage: `Record separator (use '\t' for tab)`,
				},
				&cli.BoolFlag{
					Name:     "read",
					Usage:    "Force gsheet to read from range instead of write to range. This is useful if stdin is set to a non-character device such as when running a script from cron.",
					Required: false,
					Value:    false,
				},
			},
		},
		{
			Name:     "title",
			Usage:    "Get the title of a sheet by its id",
			Action:   titleByIdAction,
			Category: "Sheets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "id",
					Usage:   "id of the spreadsheet document",
					EnvVars: []string{"GSHEET_ID"},
				},
				&cli.Int64Flag{
					Name:  "sheetid",
					Usage: "id of the sheet to get the title of",
				},
			},
		},
		{
			Name:      "sheetInfo",
			Usage:     "Dump info about the spreadsheet as json",
			Action:    sheetInfoAction,
			Category:  "Sheets",
			ArgsUsage: "SHEET_ID",
		},
		{
			Name:     "clear",
			Usage:    "Clear all values from given range",
			Action:   clearSheetAction,
			Category: "Sheets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "id",
					Usage:   "id of the spreadsheet document",
					EnvVars: []string{"GSHEET_ID"},
				},
				&cli.StringSliceFlag{
					Name:  "range",
					Usage: "Sheet range to update or get (A1 notation)",
				},
			},
		},
		{
			Name:     "newSheet",
			Usage:    "Create a new sheet",
			Action:   newSheetAction,
			Category: "Sheets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "id",
					Usage:   "id of the spreadsheet document",
					EnvVars: []string{"GSHEET_ID"},
				},
				&cli.StringFlag{
					Name:  "name",
					Usage: "title to give the new sheet",
					Value: "NewSheet",
				},
			},
		},
		{
			Name:     "deleteSheet",
			Usage:    "Delete the named sheet",
			Action:   deleteSheetAction,
			Category: "Sheets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "id",
					Usage:   "id of the spreadsheet document",
					EnvVars: []string{"GSHEET_ID"},
				},
				&cli.StringFlag{
					Name:  "name",
					Usage: "name of the sheet to delete",
				},
			},
		},
		{
			Name:     "sort",
			Usage:    "Sort a sheet by column(s)",
			Action:   sortSheetAction,
			Category: "Sheets",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "id",
					Usage:   "id of the spreadsheet document",
					EnvVars: []string{"GSHEET_ID"},
				},
				&cli.StringFlag{
					Name:  "name",
					Usage: "name of the sheet to sort",
				},
				&cli.BoolFlag{
					Name:    "ascending",
					Aliases: []string{"asc"},
					Usage:   "If set sorts in ascending order; otherwise sorts in descending order",
					Value:   false,
				},
				&cli.Int64Flag{
					Name:    "column",
					Aliases: []string{"c", "col"},
					Usage:   "Column index to sort (0=A, 1=B, ...)",
					Value:   0,
				},
			},
		},

		// Files
		{
			Name:      "createFolder",
			Usage:     "Creates a new folder",
			ArgsUsage: "NAME",
			Action:    createFolderAction,
			Category:  "Files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "parent",
					Usage:   "The id of a parent folder to act on.",
					Value:   "root",
					EnvVars: []string{"GSHEET_PARENT"},
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Delete file(s) from drive (careful, does not trash them!)",
			ArgsUsage: "FILE_ID [FILE_ID...]",
			Action:    deleteAction,
			Category:  "Files",
		},
		{
			Name:     "list",
			Usage:    "List file names and ids",
			Action:   listAction,
			Category: "Files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "parent",
					Usage:   "id of the folder to list (use 'root' for drive root)",
					EnvVars: []string{"GSHEET_PARENT"},
				},
			},
		},
		{
			Name:      "upload",
			Usage:     "Upload a file to Google Drive.",
			Action:    uploadAction,
			ArgsUsage: "[FILE]",
			Category:  "Files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "parent",
					Usage:   "id of the folder to upload to",
					Value:   "root",
					EnvVars: []string{"GSHEET_PARENT"},
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
			Action:    downloadAction,
			ArgsUsage: "FILE_ID",
			Category:  "Files",
		},
		{
			Name:      "info",
			Usage:     "Dump all file's metadata as json to stdout",
			Action:    infoAction,
			ArgsUsage: "FILE_ID",
			Category:  "Files",
		},
	},
}
