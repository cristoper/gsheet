package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func NewSheetAction(c *cli.Context) error {
	return sheetSvc.NewSheet(c.String("id"), c.String("name"))
}

func DeleteSheetAction(c *cli.Context) error {
	return sheetSvc.DeleteSheet(c.String("id"), c.String("name"))
}

func ClearSheetAction(c *cli.Context) error {
	return sheetSvc.Clear(c.String("id"), c.StringSlice("range")...)
}

func RangeSheetAction(c *cli.Context) error {
	info, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	// Set sep based on --sep flag, parsing escape sequences like \t into a rune
	sep, err := strconv.Unquote(`"` + c.String("sep") + `"`)
	if err != nil || len(sep) == 0 {
		fmt.Fprintf(os.Stderr, "Error parsing --sep; using default (',')\n")
		sep = ","
	}
	sheetSvc.Sep = rune(sep[0])

	if info.Mode()&os.ModeCharDevice > 0 {
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
		resp, err := sheetSvc.UpdateRangeCSV(c.String("id"), c.String("range"), os.Stdin)
		if err != nil {
			return err
		}
		fmt.Printf("Updated %d cells\n", resp.UpdatedCells)
	}
	return nil
}
