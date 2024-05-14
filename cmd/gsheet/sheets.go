package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func newSheetAction(c *cli.Context) error {
	if c.String("id") == "" {
		return fmt.Errorf("The --id flag is required")
	}
	return sheetSvc.NewSheet(c.String("id"), c.String("name"))
}

func deleteSheetAction(c *cli.Context) error {
	if c.String("id") == "" {
		return fmt.Errorf("The --id flag is required")
	}
	return sheetSvc.DeleteSheet(c.String("id"), c.String("name"))
}

func titleByIdAction(c *cli.Context) error {
	if c.String("id") == "" {
		return fmt.Errorf("The --id flag is required")
	}
	if c.String("sheetid") == "" {
		return fmt.Errorf("The --sheetid flag is required")
	}
	id := c.String("id")
	sheetId := c.Int64("sheetid")
	title, err := sheetSvc.TitleFromSheetId(id, sheetId)
	if err != nil {
		return err
	}
	if title == nil {
		return fmt.Errorf("No title found for sheetid %d", sheetId)
	}
	fmt.Fprint(c.App.Writer, *title)
	return nil
}

func sheetInfoAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("SHEET_ID is required")
	}
	id := c.Args().Get(0)
	info, err := sheetSvc.SpreadsheetsService().Get(id).Do()
	if err != nil {
		return err
	}
	jsonBytes, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprint(c.App.Writer, string(jsonBytes))
	return nil
}

func clearSheetAction(c *cli.Context) error {
	return sheetSvc.Clear(c.String("id"), c.StringSlice("range")...)
}

func sortSheetAction(c *cli.Context) error {
	return sheetSvc.Sort(c.String("id"), c.String("name"), c.Bool("ascending"),
		c.Int64("column"))
}

func rangeSheetAction(c *cli.Context) error {
	info, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	outInfo, err := os.Stdout.Stat()
	if err != nil {
		return err
	}

	inTty := info.Mode()&os.ModeCharDevice > 0
	outTty := outInfo.Mode()&os.ModeCharDevice > 0

	// Set sep based on --sep flag, parsing escape sequences like \t into a rune
	sep, err := strconv.Unquote(`"` + c.String("sep") + `"`)
	if err != nil || len(sep) == 0 {
		fmt.Fprintf(os.Stderr, "Error parsing --sep; using default (',')\n")
		sep = ","
	}
	sheetSvc.Sep = rune(sep[0])

	// if stdin is connected to a tty
	// or if neither stdin nor stdout are connected to a tty and there is no
	// stdin data to read (like when run from cron)
	if inTty || (!inTty && !outTty && info.Size() == 0) {
		// stdin is not connected to a pipe or file
		// get data
		vals, err := sheetSvc.GetRangeCSV(c.String("id"), c.String("range"))
		if err != nil {
			return err
		}
		fmt.Println(string(vals))
	} else {
		// otherwise stdin is connected to a pipe or file
		// send data
		if c.Bool("append") {
			// append
			resp, err := sheetSvc.AppendRangeCSV(c.String("id"), c.String("range"), os.Stdin)
			if err != nil {
				return err
			}
			fmt.Printf("Updated %d cells\n", resp.Updates.UpdatedCells)
		} else {
			// overwrite
			resp, err := sheetSvc.UpdateRangeCSV(c.String("id"), c.String("range"), os.Stdin)
			if err != nil {
				return err
			}
			fmt.Printf("Updated %d cells\n", resp.UpdatedCells)
		}
	}
	return nil
}
