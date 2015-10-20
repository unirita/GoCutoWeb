package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/unirita/gocutoweb/config"
)

var testJobnetName string = "testjobnet"
var cacheDir string
var configFileName string

func init() {
	cacheDir = filepath.Join(os.TempDir(), "gocuto")
	os.MkdirAll(cacheDir, 0777)
	configFileName = filepath.Join(os.TempDir(), "gocuto", "web.ini")

	createConfigFile()
	createDummyUtil()
	fmt.Println("tempdir = ", cacheDir)
}

func createConfigFile() {
	configContents := `[server]
listen_port = 1234
master_dir = "%v"

[log]
output_dir     = "%v"
output_level   = "info"
max_size_kb    = 10240
max_generation = 2

[jobnet]
jobnet_dir    = "%v"
`
	f, _ := os.Create(configFileName)
	dir := strings.Replace(cacheDir, "\\", "/", -1)
	f.WriteString(fmt.Sprintf(configContents, dir, dir, dir))
	f.Close()
}

func createDummyUtil() {
	var content, dummyUtil string

	if runtime.GOOS == "windows" {
		content =
			`@echo off
echo {"status":0,"message":"","pid":1234,"network":{"instance":123,"name":"%s_20151007123456789"}}
exit 0
`
		dummyUtil = filepath.Join(cacheDir, "realtime.bat")
	} else {
		content =
			`#!/bin/sh
echo \{\"status\":0,\"message\":\"\",\"pid\":1234,\"network\":\{\"instance\":123,\"name\":\"%s_20151007123456789\"\}\}
exit 0
`
		dummyUtil = filepath.Join(cacheDir, "realtime")
	}
	f, _ := os.Create(dummyUtil)
	f.WriteString(fmt.Sprintf(content, testJobnetName))
	f.Close()
	os.Chmod(dummyUtil, 0777)
}

func TestShowJSONCache(t *testing.T) {
	server := httptest.NewServer(setupHandler())
	defer server.Close()

	cacheFile := filepath.Join(cacheDir, "20151007123456.789.json")
	f, err := os.Create(cacheFile)
	if err != nil {
		t.Fatalf("File create error: %s", err)
	}
	f.WriteString("abcd")
	f.Close()
	defer os.Remove(cacheFile)

	output := testGetMessages(t, server.URL+"/caches/20151007123456.789")
	if output != "abcd" {
		t.Errorf("output => %s, want %s", output, "abcd")
	}
}

func TestNoticeJobnet(t *testing.T) {

	err := config.Load(configFileName)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(setupHandler())
	defer server.Close()

	cacheDir := filepath.Join(os.TempDir(), "gocuto")
	os.MkdirAll(cacheDir, 0777)
	cacheFile := filepath.Join(cacheDir, "testjobnet.json")
	f, err := os.Create(cacheFile)
	if err != nil {
		t.Fatalf("File create error: %s", err)
	}
	flow := `{
    "flow":"job1->job2->[job3,job4,job5->job6]->job7",
    "jobs":[
		{"name":"job1","path":"s3dladapter","param":"-b $WAPARAM1 -f $WAPARAM2"},
        {"name":"job2","node":"123.45.67.89","port":1234},
        {"name":"job7","env":"RESULT=$WAPARAM3"}
    ]
}`
	f.WriteString(flow)
	f.Close()
	defer os.Remove(cacheFile)

	form := make(url.Values)
	form.Set("jobnetwork", testJobnetName)
	form.Set("param1", "bucket1")
	form.Set("param2", "testdata")

	output := testPostFormMessages(t, server.URL+"/notice", form)
	if !strings.Contains(output, testJobnetName) {
		t.Errorf("output => %v, want contains(%v)", output, testJobnetName)
	}
}

func testGetMessages(t *testing.T, requestURL string) string {
	res, err := http.Get(requestURL)
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Logf("URL[%s]", requestURL)
		t.Errorf("res.StatusCode => %d, want %d", res.StatusCode, http.StatusOK)
	}

	resMsg, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Responce read error occured: %s", err)
	}

	return string(resMsg)
}

func testPostFormMessages(t *testing.T, requestURL string, form url.Values) string {
	res, err := http.PostForm(requestURL, form)
	//	json := `{"jobnetwork":"testjobnet","bucket":"bucket1","file":"testdata"}`
	//	res, err := http.Post(requestURL, "application/json", strings.NewReader(json))
	if err != nil {
		t.Fatalf("Error occured: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Logf("URL[%s]", requestURL)
		t.Errorf("res.StatusCode => %d, want %d", res.StatusCode, http.StatusOK)
	}

	resMsg, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Responce read error occured: %s", err)
	}

	return string(resMsg)
}
