// Integration tests for gdrive package
package gdrive

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
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

var svcDrive = func() *Service {
	svc, err := NewServiceWithCtx(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return svc
}()

// Big ol' integration script to drive the drive package
func TestDriveIntegration(t *testing.T) {
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
	testdir, err := svcDrive.CreateFolder(testName, "root")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created test directory with id %s", testdir.Id)

	// Defer the delete test so that we at least attempt to clean up if any other tests fail
	defer func() {
		err = svcDrive.DeleteFile(testdir.Id)
		if err != nil {
			t.Fatal(err)
		}
		testdirs, err := svcDrive.FilesNamed(testName, "root")
		if err != nil {
			t.Fatal(err)
		}
		if len(testdirs) != 0 {
			t.Fatal("Failed to delete testdir")
		}
		t.Log("Deleted test directory.")
	}()

	// new file
	testfile, err := svcDrive.CreateOrUpdateFile(testName+".csv", testdir.Id, strings.NewReader(testData))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created test spreadsheet named %s with id %s and mimetype %s", testfile.Name, testfile.Id, testfile.MimeType)
	if testfile.MimeType != "application/vnd.google-apps.spreadsheet" {
		t.Fatal(".csv was not converted to gsheet")
	}
	contents, err := svcDrive.FileContents(testfile.Id)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(contents), "\r", "") != testData {
		t.Fatalf("File contents do not match (got '%v' but expected '%v')", contents, []byte(testData))
	}

	// update file
	updatefile, err := svcDrive.CreateOrUpdateFile(testName+".csv", testdir.Id, strings.NewReader(modData))
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
	contents, err = svcDrive.FileContents(updatefile.Id)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(contents), "\r", "") != modData {
		t.Fatalf("File contents do not match (got '%v' but expected '%v')", contents, []byte(modData))
	}

	// 'binary' file (ie, non workspace file type)
	binaryfile, err := svcDrive.CreateOrUpdateFile(testName+".txt", testdir.Id, strings.NewReader(testData))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created test file named %s with id %s and mimetype %s", binaryfile.Name, binaryfile.Id, binaryfile.MimeType)
	if !strings.Contains(binaryfile.MimeType, "text/plain") {
		t.Fatalf("Test binary file unexpected mimetype: %s", binaryfile.MimeType)
	}
	contents, err = svcDrive.FileContents(binaryfile.Id)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ReplaceAll(string(contents), "\r", "") != testData {
		t.Fatalf("File contents do not match (got '%v' but expected '%v')", contents, []byte(testData))
	}
}
