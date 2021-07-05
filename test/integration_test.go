// Integration tests for gdrive package
package gsheet_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	. "github.com/cristoper/gsheet/gdrive"
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

var svc = func() *Service {
	svc, err := NewServiceWithCtx(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return svc
}()

// A big ol' integration script that exercises the major features of both the
// gdrive and gsheet packages
func TestGDriveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// The plan:
	// Create a test directory
	// Upload new file to test directory
	// Upload existing file to test directory
	// - Update range of data from file
	// - Get range of data from file
	// Upload a non-spreadsheet file to test directory
	// Delete test directory and all files

	testName := fmt.Sprintf("gsheet_test_%d", time.Now().UnixNano())
	testdir, err := svc.CreateFolder(testName, "root")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created test directory with id %s", testdir.Id)

	// new file
	testfile, err := svc.CreateOrUpdateFile(testName+".csv", testdir.Id, strings.NewReader(testData))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created test spreadsheet named %s with id %s and mimetype %s", testfile.Name, testfile.Id, testfile.MimeType)
	if testfile.MimeType != "application/vnd.google-apps.spreadsheet" {
		t.Fatal(".csv was not converted to gsheet")
	}
	contents, err := svc.FileContents(testfile.Id)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(contents), "\r", "") != testData {
		t.Fatalf("File contents do not match (got '%v' but expected '%v')", contents, []byte(testData))
	}

	// update file
	updatefile, err := svc.CreateOrUpdateFile(testName+".csv", testdir.Id, strings.NewReader(modData))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Updated test spreadsheet named %s with id %s and mimetype %s", updatefile.Name, updatefile.Id, updatefile.MimeType)
	if updatefile.Id != testfile.Id {
		t.Fatal("Updated Id does not match created Id")
	}
	if updatefile.MimeType != "application/vnd.google-apps.spreadsheet" {
		t.Fatal("updated file was not converted to gsheet")
	}
	contents, err = svc.FileContents(updatefile.Id)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(contents), "\r", "") != modData {
		t.Fatalf("File contents do not match (got '%v' but expected '%v')", contents, []byte(modData))
	}

	// 'binary' file (ie, non workspace file type)
	binaryfile, err := svc.CreateOrUpdateFile(testName+".txt", testdir.Id, strings.NewReader(testData))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created test file named %s with id %s and mimetype %s", binaryfile.Name, binaryfile.Id, binaryfile.MimeType)
	if !strings.Contains(binaryfile.MimeType, "text/plain") {
		t.Fatalf("Test binary file unexpected mimetype: %s", binaryfile.MimeType)
	}
	contents, err = svc.FileContents(binaryfile.Id)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(contents), "\r", "") != testData {
		t.Fatalf("File contents do not match (got '%v' but expected '%v')", contents, []byte(testData))
	}

	// Delete test dir
	err = svc.DeleteFile(testdir.Id)
	if err != nil {
		t.Fatal(err)
	}
	testdirs, err := svc.FilesNamed(testName, "root")
	if err != nil {
		t.Fatal(err)
	}
	if len(testdirs) != 0 {
		t.Fatal("Failed to delete testdir")
	}
	t.Log("Deleted test directory.")

}
