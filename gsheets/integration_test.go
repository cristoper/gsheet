// Integration tests for drive package
package gsheets

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/cristoper/gsheet/gdrive"
)

var (
	testData = `Col1,Col2,Col3
1,2,3
1,2,3
1,2,3`

	modData = `Col1,Col2,Col3
3,2,1
3,2,1
3,2,1`
)

var (
	svcDrive = func() *gdrive.Service {
		svc, err := gdrive.NewServiceWithCtx(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		return svc
	}()
	svcSheet = func() *Service {
		svc, err := NewServiceWithCtx(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		return svc
	}()
)

// Big ol' integration script to drive the sheet package
func TestSheetIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// Plan:
	// - (Create new spreadsheet with drive package)
	// - Create new sheet
	// - Update data in sheet via csv
    // - Append data to sheet
	// - Clear sheet
	// - Delete sheet
	// - (Use drive package to delete document)

	testName := fmt.Sprintf("gsheet_test_%d", time.Now().UnixNano())
	testfile, err := svcDrive.CreateOrUpdateFile(testName+".csv",
		"root", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created test spreadsheet named %s with id %s and mimetype %s", testfile.Name, testfile.Id, testfile.MimeType)

	// Defer delete test so we try to clean up if any of the other tests fail
	defer func() {
		err = svcDrive.DeleteFile(testfile.Id)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Test file deleted.")
	}()

	err = svcSheet.NewSheet(testfile.Id, "TEST")
	if err != nil {
		t.Fatal(err)
	}
	sheetId, err := svcSheet.SheetFromTitle(testfile.Id, "TEST")
	if err != nil {
		t.Fatal(err)
	}
	if sheetId == nil {
		t.Fatal("New sheet not found")
	}
	t.Logf("Created new sheet with id %d", sheetId)

	resp, err := svcSheet.UpdateRangeCSV(testfile.Id, "TEST", strings.NewReader(testData))
	if err != nil {
		t.Fatal(err)
	}
	if resp.UpdatedCells != 12 {
		t.Fatal("Unexpected number of cells updated")
	}
	vals, err := svcSheet.GetRangeCSV(testfile.Id, "TEST")
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(vals), "\n", "") !=
		strings.ReplaceAll(testData, "\n", "") {
		t.Log("Updated data does not match test data")
		t.Log(vals, []byte(testData))
		t.Fail()
	}

    // test append
    appendResp, err := svcSheet.AppendRangeCSV(testfile.Id, "TEST", strings.NewReader(testData))
    if err != nil {
        t.Fatal(err)
    }
    if appendResp.Updates.UpdatedCells != 12 {
        t.Fatal("Unexpected number of cells updated")
    }
	vals, err = svcSheet.GetRangeCSV(testfile.Id, "TEST")
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(vals), "\n", "") !=
		strings.ReplaceAll(testData + testData, "\n", "") {
		t.Log("Updated data does not match test data")
		t.Log(vals, []byte(testData))
		t.Fail()
	}

	err = svcSheet.Clear(testfile.Id, "TEST")
	if err != nil {
		t.Fatal(err)
	}
	vals, err = svcSheet.GetRangeCSV(testfile.Id, "TEST")
	if err != nil {
		t.Fatal(err)
	}
	if string(vals) != "" {
		t.Fatal("Did not clear sheet.")
	}

	err = svcSheet.DeleteSheet(testfile.Id, "TEST")
	if err != nil {
		t.Fatal(err)
	}
	gService := svcSheet.SpreadsheetsService()
	ss, err := gService.Get(testfile.Id).Do()
	if err != nil {
		t.Fatal(err)
	}
	if len(ss.Sheets) != 1 {
		t.Fatal("Unexpected number of sheets after deleting TEST")
	}
}
