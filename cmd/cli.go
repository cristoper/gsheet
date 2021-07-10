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
		// Sheets
		{
			Name:     "csv",
			Usage:    "Pipe csv data to range or read it from range",
			Action:   RangeSheetAction,
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
				&cli.StringFlag{
					Name:  "sep",
					Value: ",",
					Usage: `Record separator (use '\t' for tab)`,
				},
			},
		},
		{
			Name:     "clear",
			Usage:    "Clear all values from given range",
			Action:   ClearSheetAction,
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
			Name:     "new",
			Usage:    "Create a new sheet",
			Action:   NewSheetAction,
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
			Name:     "delete",
			Usage:    "Delete the named sheet",
			Action:   DeleteSheetAction,
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
			Action:   SortSheetAction,
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
			Action:    CreateFolderAction,
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
					Name:    "parent",
					Usage:   "id of the folder to list (use 'root' for drive root)",
					EnvVars: []string{"GSHEET_PARENT"},
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
