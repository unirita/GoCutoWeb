package define

import (
	"log"
	"os"
	"path/filepath"

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
        {"name":"job2","node":"123.45.67.89","port":1234},
        {"name":"job7","env":"RESULT=$MJjob3:RC$"}
    ]
}`
	f.WriteString(flow)
	f.Close()
	defer os.Remove(testFile)

	dynamicJobnetName, err := ReplaceJobnetTemplate(testDir, testJobNetName, "bucket", "testfile.txt")
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	jobnetFileName := filepath.Join(testDir, dynamicJobnetName+".json")
	if f, err = os.Open(jobnetFileName); err != nil {
		t.Errorf("Not found to jobnetwork file[%v].", jobnetFileName)
	}
	f.Close()
	log.Println(dynamicJobnetName)
}
