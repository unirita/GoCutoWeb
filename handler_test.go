package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShowJSONCache(t *testing.T) {
	server := httptest.NewServer(setupHandler())
	defer server.Close()

	cacheDir := filepath.Join(os.TempDir(), "gocuto")
	os.MkdirAll(cacheDir, 0777)
	cacheFile := filepath.Join(cacheDir, "20151007123456789.json")
	f, err := os.Create(cacheFile)
	if err != nil {
		t.Fatalf("File create error: %s", err)
	}
	f.WriteString("abcd")
	f.Close()
	defer os.Remove(cacheFile)

	output := testGetMessages(t, server.URL+"/caches/20151007123456789")
	if output != "abcd" {
		t.Errorf("output => %s, want %s", output, "abcd")
	}
}

func TestNoticeJobnet(t *testing.T) {
	server := httptest.NewServer(setupHandler())
	defer server.Close()

	cacheDir := filepath.Join(os.TempDir(), "gocuto")
	os.MkdirAll(cacheDir, 0777)
	cacheFile := filepath.Join(cacheDir, "testjobnet.json")
	f, err := os.Create(cacheFile)
	if err != nil {
		t.Fatalf("File create error: %s", err)
	}
	f.WriteString("abcd")
	f.Close()
	defer os.Remove(cacheFile)

	form := make(url.Values)
	form.Set("jobnetwork", "testjobnet")
	form.Set("bucket", "bucket1")
	form.Set("file", "testdata")

	output := testPostFormMessages(t, server.URL+"/notice", form)
	fmt.Println(output)
	if !strings.Contains(output, "testjobnet") {
		t.Errorf("output 0size")
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
