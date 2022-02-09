// Package gsheets is a package which provides utilities for manipulating (get
// and update data, clear, sort) Google Sheets documents.
// This can be more simple than using Google's API for common tasks (especially
// for sending and receiving csv data to and from Sheets); for anything more
// complicated use Google's golang sdk directly:
// https://pkg.go.dev/google.golang.org/api@v0.50.0/sheets/v4
package gsheets

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"google.golang.org/api/sheets/v4"
)

// Define an interface so we can mock the SpreadsheetsService if we need to
type ssService interface {
	Create(*sheets.Spreadsheet) *sheets.SpreadsheetsCreateCall
	Get(string) *sheets.SpreadsheetsGetCall
	BatchUpdate(string, *sheets.BatchUpdateSpreadsheetRequest) *sheets.SpreadsheetsBatchUpdateCall
}

// Define an interface so we can mock the SpreadsheetsValuesService if we need to
type valueService interface {
	BatchGet(string) *sheets.SpreadsheetsValuesBatchGetCall
	BatchUpdate(string, *sheets.BatchUpdateValuesRequest) *sheets.SpreadsheetsValuesBatchUpdateCall
	BatchClear(string, *sheets.BatchClearValuesRequest) *sheets.SpreadsheetsValuesBatchClearCall
}

// Service is a wrapper around both SpreadsheetsService and SpreadsheetsValuesService
type Service struct {
	Sep    rune // record separator when [un]serializing csv
	ctx    context.Context
	sheet  ssService
	values valueService
}

// NewServiceWithCtx creates and wraps a new Service with the provided context
func NewServiceWithCtx(ctx context.Context) (*Service, error) {
	ssvc, err := sheets.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &Service{
		Sep:    ',',
		ctx:    ctx,
		sheet:  ssvc.Spreadsheets,
		values: ssvc.Spreadsheets.Values,
	}, nil
}

// SpreadsheetsService returns a pointer to the wrapped SpreadsheetsService
func (svc *Service) SpreadsheetsService() *sheets.SpreadsheetsService {
	return svc.sheet.(*sheets.SpreadsheetsService)
}

// NewSheet creates a new sheet on spreadsheet identified by 'id'
func (svc *Service) NewSheet(id, title string) error {
	_, err := svc.sheet.BatchUpdate(id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			&sheets.Request{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: title,
					},
				},
			},
		},
	}).Context(svc.ctx).Do()
	return err
}

// SheetFromTitle returns the sheetID for the sheet with 'title' in the
// spreadsheet doc identified by 'id'.
// If no error is encountered and no matching sheet is found, both return
// values will be nil.
func (svc *Service) SheetFromTitle(id, title string) (*int64, error) {
	ss, err := svc.sheet.Get(id).Context(svc.ctx).Do()
	if err != nil {
		return nil, err
	}
	var sheetId *int64
	for _, sheet := range ss.Sheets {
		if sheet.Properties.Title == title {
			sheetId = &sheet.Properties.SheetId
			break
		}
	}
	return sheetId, nil
}

// DeleteSheet deletes the sheet with 'title' from spreadsheet doc identified
// by 'id'
func (svc *Service) DeleteSheet(id, title string) error {
	// find sheet matching title
	sheetId, err := svc.SheetFromTitle(id, title)
	if err != nil {
		return err
	}
	if sheetId == nil {
		return errors.New("No sheet found with title: " + title)
	}

	_, err = svc.sheet.BatchUpdate(id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			&sheets.Request{
				DeleteSheet: &sheets.DeleteSheetRequest{
					SheetId: *sheetId,
				},
			},
		},
	}).Context(svc.ctx).Do()
	return err
}

// GetRangeRaw gets unformatted values in 'a1Range' from the spreadsheet doc
// identified by 'id'.
// You must type switch the resulting [][]interface{} (outer slice is rows, inner
// slice is value per column)
// A1 syntax: https://developers.google.com/sheets/api/guides/concepts
func (svc *Service) GetRangeRaw(id string, a1Range string) ([][]interface{}, error) {
	resp, err := svc.values.BatchGet(id).
		Context(svc.ctx).
		MajorDimension("ROWS").
		Ranges(a1Range).
		ValueRenderOption("UNFORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, err
	}
	return resp.ValueRanges[0].Values, nil
}

// GetRangeFormatted gets formatted values in 'a1Range' from the spreadsheet
// doc identified by 'id'.
// All values are returned as strings, formatted as they display in the spreadsheet document
func (svc *Service) GetRangeFormatted(id string, a1Range string) ([][]string, error) {
	resp, err := svc.values.BatchGet(id).
		Context(svc.ctx).
		MajorDimension("ROWS").
		Ranges(a1Range).
		ValueRenderOption("FORMATTED_VALUE").
		Do()
	if err != nil {
		return nil, err
	}

	values := resp.ValueRanges[0].Values
	// make a [][]string to hold typecast values
	var stringVals = make([][]string, len(values))
	for i := range stringVals {
		stringVals[i] = make([]string, len(values[i]))
	}
	// cast each interface{} to a string
	for r, row := range values {
		for c, v := range row {
			stringVals[r][c] = v.(string)
		}
	}
	return stringVals, nil
}

// GetRangeCSV returns values in 'a1Range' from the spreadsheet doc identified
// by 'id' in csv format.
func (svc *Service) GetRangeCSV(id, a1Range string) ([]byte, error) {
	rows, err := svc.GetRangeFormatted(id, a1Range)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte{})
	csvW := csv.NewWriter(buf)
	csvW.Comma = svc.Sep
	err = csvW.WriteAll(rows)
	if err != nil {
		return nil, err
	}
	csvW.Flush()
	err = csvW.Error()
	return buf.Bytes(), err
}

// UpdateRangeRaw updates the values in 'a1Range' in the spreadsheet doc
// identified by 'id' to 'values'.
// Each value in 'values' must be a string, float, int, or bool and must fit
// within the dimensions of 'a1Range'.
func (svc *Service) UpdateRangeRaw(id, a1Range string, values [][]interface{}) (*sheets.UpdateValuesResponse, error) {

	resp, err := svc.values.BatchUpdate(id, &sheets.BatchUpdateValuesRequest{
		Data: []*sheets.ValueRange{
			&sheets.ValueRange{
				MajorDimension: "ROWS",
				Range:          a1Range,
				Values:         values,
			},
		},
		ValueInputOption: "RAW",
	}).
		Context(svc.ctx).
		Do()
	if err != nil {
		return nil, err
	}
	return resp.Responses[0], nil
}

// UpdateRangeStrings update values in 'a1Range' in the spreadsheet doc
// identified by 'id' to 'values'.
// Values will be parsed by Google Sheets as if they were typed in by the user
// (so strings containing numerals may be converted to numbers, etc.)
func (svc *Service) UpdateRangeStrings(id, a1Range string, values [][]string) (*sheets.UpdateValuesResponse, error) {

	// make a [][]interface{} to hold typecast values
	var vals = make([][]interface{}, len(values))
	for i := range vals {
		vals[i] = make([]interface{}, len(values[i]))
	}
	// cast each string to interface{}
	for r, row := range values {
		for c, v := range row {
			vals[r][c] = v
		}
	}

	resp, err := svc.values.BatchUpdate(id, &sheets.BatchUpdateValuesRequest{
		Data: []*sheets.ValueRange{
			&sheets.ValueRange{
				MajorDimension: "ROWS",
				Range:          a1Range,
				Values:         vals,
			},
		},
		ValueInputOption: "USER_ENTERED",
	}).
		Context(svc.ctx).
		Do()
	if err != nil {
		return nil, err
	}
	return resp.Responses[0], nil
}

// UpdateRangeCSV update values in 'a1Range' in the spreadsheet doc identified
// by 'id' to 'values'.
// 'values' is an io.Reader which supplies text in csv format.
// Values will be parsed by Google Sheets as if they were typed in by the user
// (so strings containing numerals may be converted to numbers, etc.)
func (svc *Service) UpdateRangeCSV(id, a1Range string, values io.Reader) (*sheets.UpdateValuesResponse, error) {
	csvR := csv.NewReader(values)
	csvR.FieldsPerRecord = -1 // disable field checks
	csvR.Comma = svc.Sep
	rows, err := csvR.ReadAll()
	if err != nil {
		return nil, err
	}
	return svc.UpdateRangeStrings(id, a1Range, rows)
}

// Clear clears the value of all 'a1Ranges' in the spreadsheet doc identified
// by 'id'.
func (svc *Service) Clear(id string, a1Ranges ...string) error {
	_, err := svc.values.BatchClear(id, &sheets.BatchClearValuesRequest{
		Ranges: a1Ranges,
	}).Context(svc.ctx).Do()
	return err
}

// Sort sorts the sheet titled 'name' on the spreadsheet doc identified by 'id'
// by 'column'.
// If asc is true, sort ascending; otherwise sort descending
// Column is the column index rather than A1 notation (0=A, 1=B, ...)
func (svc *Service) Sort(id, name string, asc bool, column int64) error {
	sheetId, err := svc.SheetFromTitle(id, name)
	if err != nil {
		return err
	}

	if sheetId == nil {
		return fmt.Errorf("No sheet titled %s found", name)
	}

	order := "DESCENDING"
	if asc {
		order = "ASCENDING"
	}

	sortSpec := sheets.Request{
		SortRange: &sheets.SortRangeRequest{
			Range: &sheets.GridRange{
				SheetId: *sheetId,
			},
			SortSpecs: []*sheets.SortSpec{
				&sheets.SortSpec{
					DimensionIndex: column,
					SortOrder:      order,
				}},
		},
	}
	_, err = svc.sheet.BatchUpdate(id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&sortSpec},
	}).Context(svc.ctx).Do()
	return err
}
