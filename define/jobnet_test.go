package define

import (
	"os"
	"path/filepath"
	"strings"

	"testing"
)

var testDir string = setTestDir()

func setTestDir() string {
	cacheDir := filepath.Join(os.TempDir(), "gocuto")
	os.MkdirAll(cacheDir, 0777)
	return cacheDir
}

const testJobNetName string = "testjobnet"

func TestReplaceJobnetTemplate(t *testing.T) {
	testFile := filepath.Join(testDir, testJobNetName+".json")
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("File create error: %s", err)
	}
	flow := `{
    "flow":"job1->job2->[job3,job4,job5->job6]->job7",
    "jobs":[
        {"name":"job2","node":"123.45.67.89","port":1234,"param":"$WAPARAM3"},
        {"name":"job7","env":"BUCKET=$WAPARAM1","param":"-f $WAPARAM2"}
    ]
}`
	f.WriteString(flow)
	f.Close()
	defer os.Remove(testFile)

	params := []string{"bucket1", "testfile.txt"}
	dynamicJobnetName, err := ReplaceJobnetTemplate(testDir, testJobNetName, params)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	jobnetFileName := filepath.Join(testDir, dynamicJobnetName+".json")
	if f, err = os.Open(jobnetFileName); err != nil {
		t.Fatalf("Not found to jobnetwork file[%v].", jobnetFileName)
	}
	buf := make([]byte, 512)
	size, _ := f.Read(buf)
	if size == 0 {
		t.Errorf("jobnet file[%v] is empty.", jobnetFileName)
	}
	contents := string(buf)
	if !strings.Contains(contents, "bucket1") {
		t.Errorf("file contents [%v], not contents [%v]", contents, "bucket1")
	}
	if !strings.Contains(contents, "testfile.txt") {
		t.Errorf("file contents [%v], not contents [%v]", contents, "testfile.txt")
	}
	f.Close()
}
